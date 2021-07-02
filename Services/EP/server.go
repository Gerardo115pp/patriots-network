package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/gorilla/websocket"
)

const DEBUG = false
const FILES_DIRECTORY = "./users_files"
const SERVER_DATA = "./operational_data"

var JD_ADDRESS string = ""
var JD_PORT string = "4000"
var JD_HOST string = "0.0.0.0"

type GWdata struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Server struct {
	port                 int
	hasher               hash.Hash
	users                []User
	bad_response         []byte
	ok                   []byte
	socket               websocket.Upgrader
	messages             []Message
	chat_message_channel chan *Message
	user_join            chan *User
	user_leaves          chan *User
	GWserver             GWdata
}

func (self *Server) debugUsersLogged(identified string) {
	fmt.Println("\n\nLoggeados on: ", identified)
	for _, users := range self.users {
		fmt.Printf(" %s : %t\n", users.User_name, users.Is_logged)
	}
	fmt.Print("\n\n")
}

func (self Server) debugLog(debug_message string) {
	if DEBUG {
		fmt.Printf("\n\n(DEBUG MESSAGE) %s\n\n", debug_message)
	}
}

func (self *Server) init(port int) {
	self.port = port
	self.hasher = sha1.New()
	self.ok = []byte("{\"response\":\"ok\"}")
	self.bad_response = []byte("{\"response\":\"bad\"}")
	self.socket = *(self.setSocket())
	self.chat_message_channel = make(chan *Message)
	self.user_join = make(chan *User)
	self.user_leaves = make(chan *User)

	// getting gw data
	response, err := http.Get(fmt.Sprintf("http://%s/GW", JD_ADDRESS))
	if err != nil {
		fmt.Println("Warning:", err.Error())
		return
	}

	if response.StatusCode == 200 {
		response_data, err := self.parseFormAsMap(response)
		if err != nil {
			log.Fatal(err)
		}
		self.GWserver = GWdata{
			Host: response_data["host"],
			Port: response_data["port"],
		}
		fmt.Println("GW connection recived")
	} else {
		fmt.Println("Warning:", response.Status, "Failed to connect to JD, couldnt stablish a connection to a GW")
	}
}

func (self *Server) addUser(user_name string, connection *websocket.Conn) (*User, error) {
	var new_user *User = new(User)

	for h, user := range self.users {
		if user.User_name == user_name {
			if !user.Is_logged {
				new_user = &(self.users[h])
				new_user.connection = connection
				new_user.Is_logged = true
				return new_user, nil
			} else {
				return new_user, fmt.Errorf("User is already logged")
			}
		} else {
			fmt.Printf("\ndebug: %s !== %s\n", user.User_name, user_name)
		}
	}
	return new_user, nil
}

func (self *Server) broadcast(msg *Message, resend bool) {
	fmt.Printf("\nBroadcasting message '''%s''' from '%s'\n", msg.Body, msg.Sender)
	for h, user := range self.users {
		if user.Is_logged && (user.User_name != msg.Sender || resend) {
			self.users[h].write(msg)
		} else if user.User_name == msg.Sender || !resend {
			fmt.Println("Skiping user '", user.User_name)
		}

		if DEBUG {
			fmt.Printf("\nDEBUG---->\nUSER:%v\nMESSSAGE: %v\n resend: %t\n", user, msg, resend)
		}
	}
}

func (self *Server) corsEnabler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	}
}

func (self *Server) castMessageToBlockchain(message *Message) {
	message_json, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	gw_url := fmt.Sprintf("http://%s:%s/register-transaction", self.GWserver.Host, self.GWserver.Port)
	http.Post(gw_url, "application/json", bytes.NewBuffer(message_json))
	fmt.Println("Transaction:", string(message_json), "casted to ", gw_url)
}

func (self *Server) getHash(s string) (hsh string) {
	_, err := self.hasher.Write([]byte(s))
	defer self.hasher.Reset()
	panicIfErr(err)
	hsh = hex.EncodeToString(self.hasher.Sum(nil))
	return
}

func (self *Server) generateResponse(msg string) (response []byte) {
	response = []byte(fmt.Sprintf("{\"response\":\"%s\"}", msg))
	return
}

func (self *Server) getUsers(response http.ResponseWriter, request *http.Request) {
	fmt.Println("\nServing 'GET' users request: ", request.Host)
	if DEBUG {
		self.debugUsersLogged(fmt.Sprintf("getUsers from %s", request.RemoteAddr))
	}
	users_json, err := json.Marshal(self.users)
	panicIfErr(err)

	self.setHeaders(&response)
	_, _ = response.Write(users_json)
}

func (self *Server) getPortStr() string {
	return fmt.Sprintf(":%d", self.port)
}

func (self *Server) getUserName(response http.ResponseWriter, request *http.Request) {
	var user_hash string = request.FormValue("hash")
	if user_hash == "" {
		_, _ = response.Write(self.bad_response)
	}
	self.setHeaders(&response)
	for _, user := range self.users {
		if user.hash == user_hash {
			_, _ = response.Write(self.generateResponse(user.User_name))
			return
		}
	}
	_, _ = response.Write(self.bad_response)
}

func (self *Server) loadState() {
	var users_state_file string = path.Join(SERVER_DATA, "users.json")
	var users_state []byte
	var err error

	err = self.requestMessagesToGW()
	if err != nil {
		log.Fatal(err)
	}

	// loading users
	if self.pathExists(users_state_file) {
		users_state, err = ioutil.ReadFile(users_state_file)
		panicIfErr(err)
		panicIfErr(json.Unmarshal(users_state, &(self.users)))

		for h := range self.users {
			self.users[h].chat = self
			self.users[h].Is_logged = false
			self.users[h].hash = self.getHash(self.users[h].User_name)
		}
	}
}

func (self *Server) lisintOnWebsocket() {
	var new_message *Message
	var do_broadcast, resend bool

	for {
		do_broadcast = false
		resend = false
		select {
		case new_message = <-self.chat_message_channel:
			self.castMessageToBlockchain(new_message)
			self.messages = append(self.messages, *new_message)
			do_broadcast = true
			resend = true
			if err := self.saveServerStatus(); err != nil && DEBUG {
				fmt.Println("Error while saveing server status:", err)
			}
		case user := <-self.user_join:
			new_message = newMessage(fmt.Sprintf("El usuario %s se a unido a la sala", user.User_name), user.User_name, user.Color, true, false)
			do_broadcast = true
			if err := self.saveServerStatus(); err != nil && DEBUG {
				fmt.Println("Error while saveing server status:", err)
			}
		case user := <-self.user_leaves:
			if DEBUG {
				self.debugUsersLogged(fmt.Sprintf("User: %s logged out", user.User_name))
			}
			new_message = newMessage(fmt.Sprintf("El usuario %s a dejado la sala", user.User_name), user.User_name, user.Color, true, false)
			do_broadcast = true
		}

		if do_broadcast {
			self.broadcast(new_message, resend)
		}
	}
}

func (self *Server) isUserRegisterd(user_name string) (user_requested *User, is_registerd bool) {
	user_requested = new(User)
	for h, user := range self.users {
		if user.User_name == user_name {
			self.debugLog(fmt.Sprintf("User '%s' was already logged", user_name))
			user_requested = &(self.users[h])

			is_registerd = true
			return
		}
	}
	is_registerd = false
	return
}

func (self *Server) setSocket() *websocket.Upgrader {
	var socket websocket.Upgrader = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(r *http.Request) bool {
			log.Printf("\n%s %s%s %v\n", r.Method, r.Host, r.RequestURI, r.Proto)
			return true
		},
	}
	return &socket
}

func (self *Server) socketRequestHandler(response http.ResponseWriter, request *http.Request) {
	connection, err := self.socket.Upgrade(response, request, nil)
	var sender string
	panicIfErr(err)
	var url_parameters url.Values = request.URL.Query()
	if url_parameters.Get("user_name") == "" {
		sender = fmt.Sprintf("anon-%d", getRandomInt64())
	} else {
		sender = url_parameters.Get("user_name")
	}

	new_user, err := self.addUser(sender, connection)
	if err != nil {
		fmt.Println("\nError on user login: ", err)
		return
	}

	fmt.Printf("\nUSER: %s logged in\n", new_user.User_name)

	self.user_join <- new_user
	new_user.read()
}

func (self Server) setHeaders(response *http.ResponseWriter) {
	(*response).Header().Set("Content-Type", "application/json")
	(*response).Header().Set("Access-Control-Allow-Origin", "*")
}

func (self Server) pathExists(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func (self *Server) requestMessagesToGW() error {
	response, err := http.Get(fmt.Sprintf("http://%s:%s/blockchain?format=data", self.GWserver.Host, self.GWserver.Port))
	if err != nil {
		return err
	}

	var new_message *Message = new(Message)
	var messages_bytes []byte
	var encoded_messages []string
	messages_bytes, err = ioutil.ReadAll(response.Body)

	err = json.Unmarshal(messages_bytes, &encoded_messages)
	if err != nil {
		return err
	}

	fmt.Printf("Recived %d messages from GW\n", len(encoded_messages))
	for _, encoded_message := range encoded_messages {
		if err = json.Unmarshal([]byte(encoded_message), new_message); err != nil {
			fmt.Println("Error while loading message:", encoded_message, ".")
			return err
		}
		self.messages = append(self.messages, *new_message)
	}

	return err
}

func (self *Server) retriveFile(response http.ResponseWriter, request *http.Request) {
	fmt.Println("\nServing a 'POST' file request from:", request.RemoteAddr)

	user_hash := request.FormValue("user_uuid")
	file, header, err := request.FormFile("file")
	panicIfErr(err)
	panicIfErr(self.saveFile(file, header))
	self.setHeaders(&response)
	_, _ = response.Write(self.ok)
	self.userSendedFile(user_hash, header)
}

func (self *Server) retriveMessages(response http.ResponseWriter, request *http.Request) {
	self.setHeaders(&response)
	if request.Method == http.MethodGet {
		var messages_data []byte
		messages_data, err := json.Marshal(self.messages)
		panicIfErr(err)
		_, _ = response.Write(messages_data)
	} else {
		_, _ = response.Write(self.bad_response)
	}
}

func (self *Server) saveFile(file multipart.File, file_header *multipart.FileHeader) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if !self.pathExists(FILES_DIRECTORY) {
		panicIfErr(os.Mkdir(FILES_DIRECTORY, 0755))
	}
	return ioutil.WriteFile(fmt.Sprintf("%s/%s", FILES_DIRECTORY, file_header.Filename), data, 0666)
}

func (self *Server) saveServerStatus() error {
	var users_data []byte

	users_data, err := json.Marshal(self.users)
	if err != nil {
		return err
	}

	if !self.pathExists(SERVER_DATA) {
		err = os.Mkdir(SERVER_DATA, 0755)
		if err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/users.json", SERVER_DATA), users_data, 0666); err != nil {
		return err
	}
	return nil
}

func (self *Server) parseFormAsMap(request *http.Response) (map[string]string, error) {
	var parsed_form map[string]string = make(map[string]string)
	var err error
	body_data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return make(map[string]string), err
	}
	err = json.Unmarshal(body_data, &parsed_form)
	return parsed_form, err
}

func (self *Server) registerUser(response http.ResponseWriter, request *http.Request) {
	var user_name string = request.FormValue("user_name")
	var color string = request.FormValue("color")

	var new_user *User
	new_user, user_is_registred := self.isUserRegisterd(user_name)
	if !user_is_registred {
		fmt.Printf("\nRegistrated new user '%s'\n", user_name)
		new_user = new(User)
		new_user.init(user_name, color, self.getHash(user_name), self)
		self.users = append(self.users, *new_user)
	}
	self.setHeaders(&response)
	_, _ = response.Write(self.generateResponse(new_user.hash))
}

func (self *Server) run() {

	self.loadState()

	var file_server http.Handler = http.FileServer(http.Dir("./users_files"))

	http.HandleFunc("/register-user", self.registerUser)
	http.HandleFunc("/retrive-messages-history", self.retriveMessages)
	http.HandleFunc("/login", self.socketRequestHandler)
	http.HandleFunc("/get-username", self.getUserName)
	http.HandleFunc("/get-users", self.getUsers)
	http.HandleFunc("/send-file", self.retriveFile)
	http.HandleFunc("/static/", self.corsEnabler(http.StripPrefix("/static/", file_server)))

	go self.lisintOnWebsocket()

	fmt.Println("Listening on:", self.getPortStr())
	panicIfErr(http.ListenAndServe(self.getPortStr(), nil))
}

func (self *Server) userSendedFile(user_uuid string, file_header *multipart.FileHeader) {
	var sender *User
	for h, user := range self.users {
		if user.hash == user_uuid {
			sender = &(self.users[h])
		}
	}
	fmt.Printf("\n\n'%s' sended a file called '%s'\n\n", sender.User_name, file_header.Filename)

	var new_message *Message = newMessage(
		file_header.Filename,
		sender.User_name,
		sender.Color,
		false,
		true)
	self.chat_message_channel <- new_message
}

func main() {
	if param := os.Getenv("JD_PORT"); param != "" {
		JD_PORT = param
	}

	if param := os.Getenv("JD_HOST"); param != "" {
		JD_HOST = param
	}
	JD_ADDRESS = fmt.Sprintf("%s:%s", JD_HOST, JD_PORT)

	var server *Server = new(Server)
	server.init(5000)
	server.run()
}
