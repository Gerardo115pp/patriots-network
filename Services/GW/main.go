package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	patriot_blockchain "github.com/Gerardo115pp/PatriotLib/PatriotBlockchain"
	patriot_router "github.com/Gerardo115pp/PatriotLib/PatriotRouter"
)

const JD_address = "192.168.1.79:4000"
const PEERS_FILE = "peers.json"

type GWnode struct {
	blockchain  *patriot_blockchain.Blockchain
	router      *patriot_router.Router
	used_port   int
	known_peers []Peer
}

func (self *GWnode) boot() {
	// construction
	self.router = patriot_router.CreateRouter()

	// loading peers
	var loaded_peers []Peer = make([]Peer, 0)

	if patriot_blockchain.PathExists("peers.json") {
		data, err := ioutil.ReadFile("peers.json")
		if err != nil {
			patriot_blockchain.LogFatal(err)
		}

		err = json.Unmarshal(data, &loaded_peers)
		if err != nil {
			patriot_blockchain.LogFatal(err)
		}
	} else {
		loaded_peers = self.requestPeers()
		peers_data, err := json.Marshal(loaded_peers)
		if err != nil {
			patriot_blockchain.LogFatal(err)
		}
		ioutil.WriteFile(PEERS_FILE, peers_data, 0700)
	}

	// we are the boot peer
	if len(loaded_peers) == 0 {
		fmt.Println("We are the boot peer")
	}

	self.registerAsPeer()

	// loading blockchain
	if patriot_blockchain.PathExists("blockchain") {
		self.blockchain.Load("blockchain")
	} else {

	}

	self.run()
}

func (self *GWnode) parseFormAsMap(request *http.Request) (map[string]string, error) {
	var parsed_form map[string]string = make(map[string]string)
	var err error
	if content_type := request.Header.Get("Content-Type"); content_type == "application/json" {
		body_data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return make(map[string]string), err
		}
		err = json.Unmarshal(body_data, &parsed_form)
	}
	return parsed_form, err
}

func (self *GWnode) registerAsPeer() {
	var use_port int = 4005
	for {
		if conn, _ := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprint(use_port)), time.Millisecond*800); conn == nil {
			break
		} else {
			conn.Close()
			fmt.Println("Port", use_port, "is busy")
			use_port++
		}
	}
	form_data := []byte(fmt.Sprintf("{\"port\": \"%d\"}", use_port))
	response, _ := http.Post(fmt.Sprintf("http://%s/GW", JD_address), "application/json", bytes.NewBuffer(form_data))
	if response.StatusCode != 200 {
		patriot_blockchain.LogFatal(fmt.Errorf("couldnt register to JD"))
	}
	self.used_port = use_port
}

func (self *GWnode) requestPeers() []Peer {
	var network_peers []Peer
	response, err := http.Get(fmt.Sprintf("http://%s/GW", JD_address))
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	if response.StatusCode == 404 {
		// we are the boot peer
		return make([]Peer, 0)
	}

	body_data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		patriot_blockchain.LogFatal(err)
	}

	err = json.Unmarshal(body_data, &network_peers)
	if err != nil {
		fmt.Println("Warning on requestPeers:", err.Error())
	}
	return network_peers
}

func (self *GWnode) registerTransaction(response http.ResponseWriter, request *http.Request) {
	transaction_data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(400)
		fmt.Fprint(response, "invalid format")
	} else {
		fmt.Println("Transaction:", string(transaction_data))
	}
}

func (self *GWnode) run() {
	self.router.RegisterRoute(patriot_router.NewRoute("/register-transaction", true), self.registerTransaction)

	fmt.Println("Awaiting connections on port:", self.used_port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", self.used_port), self.router); err != nil {
		patriot_blockchain.LogFatal(err)
	}
}

func main() {
	var GW_node *GWnode = new(GWnode)
	GW_node.boot()
}
