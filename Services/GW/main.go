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

const TRANSACTION_LIMIT = 2
const JD_address = "192.168.1.79:4000"
const PEERS_FILE = "peers.json"

type GWnode struct {
	blockchain   *patriot_blockchain.Blockchain
	router       *patriot_router.Router
	used_port    int
	known_peers  []Peer
	transactions []string
}

func (self *GWnode) boot() {
	// construction
	self.router = patriot_router.CreateRouter()
	self.transactions = make([]string, 0)

	// loading peers
	var loaded_peers []Peer = make([]Peer, 0)

	// prompt to errors, when runing multiple GW from the same directory
	// if patriot_blockchain.PathExists("peers.json") {
	// 	data, err := ioutil.ReadFile("peers.json")
	// 	if err != nil {
	// 		patriot_blockchain.LogFatal(err)
	// 	}

	// 	err = json.Unmarshal(data, &loaded_peers)
	// 	if err != nil {
	// 		patriot_blockchain.LogFatal(err)
	// 	}
	// }

	if len(loaded_peers) == 0 {
		loaded_peers = self.requestPeers()
		peers_data, err := json.Marshal(loaded_peers)
		if err != nil {
			patriot_blockchain.LogFatal(err)
		}
		ioutil.WriteFile(PEERS_FILE, peers_data, 0700)
	}

	self.known_peers = loaded_peers

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

func (self *GWnode) broadcastTransactionToPeers(transaction []byte) {
	var http_client *http.Client = &http.Client{}
	var request *http.Request
	var request_body *bytes.Buffer = bytes.NewBuffer(transaction)

	for h, peer := range self.known_peers {
		fmt.Printf("Broadcasting to peers(%d/%d): %s", h+1, len(self.known_peers), peer.Port)

		request_url := fmt.Sprintf("http://%s:%s/register-transaction", peer.Host, peer.Port)
		request, _ = http.NewRequest("POST", request_url, request_body)
		request.Header.Set("X-from-peer", fmt.Sprint(self.used_port))

		http_client.Do(request)
	}
	fmt.Println("\nSuccess")
}

func (self *GWnode) handlePeers(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		fmt.Println("New peer connection recived")

		var new_peer Peer
		var peer_data map[string]string

		peer_data, err := self.parseFormAsMap(request)
		if err != nil {
			fmt.Println("Couldnt register peer due to a request format error from JD:", err.Error())
			return
		}

		new_peer.Host = peer_data["host"]
		new_peer.Port = peer_data["port"]

		self.known_peers = append(self.known_peers, new_peer)
		fmt.Println("Registred the new peer:", new_peer.Json())
		response.WriteHeader(200)
		fmt.Fprintf(response, "ok")
	}
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
	response, err := http.Get(fmt.Sprintf("http://%s/GW-all", JD_address))
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
		fmt.Println("Body data was:", string(body_data))
	}
	return network_peers
}

func (self *GWnode) registerTransaction(response http.ResponseWriter, request *http.Request) {
	transaction_data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(400)
		fmt.Fprint(response, "invalid format")
	} else {
		self.transactions = append(self.transactions, string(transaction_data))
		fmt.Printf("new transaction, now we have %d\n", len(self.transactions))
		if len(self.transactions) == TRANSACTION_LIMIT {
			// tell JD to mine the block
			fmt.Println("Transaction limit reached, requesting a blockmine")
		} else if sender := request.Header.Get("X-from-peer"); sender == "" {
			self.broadcastTransactionToPeers(transaction_data)
		} else {
			fmt.Println(" Got new transaction from", sender)
		}

	}
}

func (self *GWnode) run() {
	self.router.RegisterRoute(patriot_router.NewRoute("/register-transaction", true), self.registerTransaction)
	self.router.RegisterRoute(patriot_router.NewRoute("/peers", true), self.handlePeers)

	fmt.Println("Awaiting connections on port:", self.used_port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", self.used_port), self.router); err != nil {
		patriot_blockchain.LogFatal(err)
	}
}

func main() {
	var GW_node *GWnode = new(GWnode)
	GW_node.boot()
}
