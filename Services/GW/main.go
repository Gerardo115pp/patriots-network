package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	patriot_blockchain "github.com/Gerardo115pp/PatriotLib/PatriotBlockchain"
	patriot_router "github.com/Gerardo115pp/PatriotLib/PatriotRouter"
)

const TRANSACTION_LIMIT = 2
const JD_address = "192.168.1.79:4000"
const PEERS_FILE = "peers.json"

var GW_NODE_NAME string

func getBlockchainFilename() string {
	filename := "blockchain"
	if GW_NODE_NAME != "" {
		filename += fmt.Sprintf(".%s", GW_NODE_NAME)
	}
	return filename
}

type GWnode struct {
	blockchain   *patriot_blockchain.Blockchain
	router       *patriot_router.Router
	used_port    int
	known_peers  []Peer
	transactions []string
	JD_Code      uint64
}

func (self *GWnode) addBlock(response http.ResponseWriter, request *http.Request) {
	var new_block *patriot_blockchain.Block = new(patriot_blockchain.Block)

	fmt.Println("New block arrived.")
	block_data, err := ioutil.ReadAll(request.Body)
	if err == nil {
		err = json.Unmarshal(block_data, new_block)
	}

	if err != nil {
		fmt.Println("Error while parsing block data:", err.Error(), "\nRequest body was:", string(block_data))
		response.WriteHeader(400)
		return
	}

	err = self.blockchain.AddBlock(new_block)
	if err != nil {
		fmt.Println("Error while adding block:", err.Error())
	}
	self.blockchain.Save(getBlockchainFilename())
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
	self.registerAsPeer()

	if len(loaded_peers) == 0 {
		loaded_peers = self.requestPeers()
		peers_data, err := json.Marshal(loaded_peers)
		if err != nil {
			patriot_blockchain.LogFatal(err)
		}
		ioutil.WriteFile(PEERS_FILE, peers_data, 0700)
	}

	fmt.Println("Peers loaded:", len(loaded_peers))
	self.known_peers = loaded_peers

	// we are the boot peer
	if len(loaded_peers) == 0 {
		fmt.Println("We are the boot peer")
	}

	self.loadBlockchain()

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

func (self *GWnode) clearTransactions() {
	fmt.Print("Clearing", len(self.transactions), "transactions:")
	self.transactions = make([]string, 0)
	fmt.Println(" Done")
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

func (self *GWnode) handleBlockchainDataReques(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		var format string = request.URL.Query().Get("format")
		switch format {
		case "complete":
			// return all the blockchain as is
			response.WriteHeader(200)
			response.Write(self.blockchain.ToBytes())
		case "data":
			// return just a list with the data field of each block
			response.WriteHeader(200)
			response.Write(self.blockchain.DataToJson())
		default:
			response.WriteHeader(400)
		}
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (self *GWnode) loadBlockchain() {
	var err error
	self.blockchain = patriot_blockchain.CreateBlockchain()

	if patriot_blockchain.PathExists(fmt.Sprintf("%s.json", getBlockchainFilename())) {
		// load form file
		fmt.Println("Loading blockchain from local storage")
		err = self.blockchain.Load(getBlockchainFilename())
	} else {
		fmt.Println("warning blockchain file doesnt exists:", getBlockchainFilename())
	}

	if len(self.known_peers) > 0 {
		other_blockchain := self.requestBlockchainToPeer()
		if err != nil && other_blockchain.Length() > 0 || other_blockchain.Length() > self.blockchain.Length() {
			fmt.Println("Blockchain loaded from peer")
			self.blockchain = other_blockchain
			self.blockchain.Save(getBlockchainFilename())
			err = nil
		}

	}

	if err != nil {
		fmt.Println("Unable to get valid blockchain copy")
		os.Exit(0)
	}
	fmt.Println("Blocks loaded:", self.blockchain.Length())
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
	response_data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		patriot_blockchain.LogFatal(err)
	}

	self.JD_Code = patriot_blockchain.StringToUint64(string(response_data))
	fmt.Println("JD_code:", self.JD_Code)
	self.used_port = use_port
}

func (self *GWnode) requestBlockMining() bool {
	var request_body *bytes.Buffer
	var request *http.Request
	var http_client *http.Client = new(http.Client)

	var previous_hash string = patriot_blockchain.GENESIS_HASH
	transactions_data, err := json.Marshal(self.transactions)
	if err != nil {
		patriot_blockchain.LogFatal(err)
	}

	if self.blockchain.Length() > 0 {
		previous_hash = self.blockchain.HeadHash()
	}

	encoded_transactions := base64.RawStdEncoding.EncodeToString(transactions_data)
	form_data := fmt.Sprintf("{\"transactions\": \"%s\", \"block_num\": %d, \"previous\": \"%s\" }", encoded_transactions, self.blockchain.Length(), previous_hash)
	request_body = bytes.NewBufferString(form_data)

	request, _ = http.NewRequest("PUT", fmt.Sprintf("http://%s/XZ", JD_address), request_body)

	response, err := http_client.Do(request)
	if err != nil {
		fmt.Println("Warning, error on send mining request:", err.Error())
	}

	if response.StatusCode == 200 {
		self.syncTransactions()
		return true
	} else {
		fmt.Print("JD gave a negative awnser to mine block request..\n")
		return false
	}
}

func (self *GWnode) requestBlockchainToPeer() *patriot_blockchain.Blockchain {
	var peer_blockchain *patriot_blockchain.Blockchain = patriot_blockchain.CreateBlockchain()
	var peer_index int = 0
	if len(self.known_peers) > 1 {
		rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(self.known_peers) - 1) // random number int [0,peers_count)
	}

	response, err := http.Get(fmt.Sprintf("http://%s/blockchain?format=complete", self.known_peers[peer_index].String()))
	if err != nil {
		fmt.Println("Warning on request blockchain to peer:", err.Error())
		return peer_blockchain
	}
	blockchain_bytes, err := ioutil.ReadAll(response.Body)
	if err == nil {
		peer_blockchain.FromBytes(blockchain_bytes)
	}

	if err != nil {
		fmt.Println("Peer blockchain request yielded an invalid blockchain:", err.Error())
	}
	return peer_blockchain
}

func (self *GWnode) requestPeers() []Peer {
	var network_peers []Peer
	response, err := http.Get(fmt.Sprintf("http://%s/GW-all?code=%d", JD_address, self.JD_Code))
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

	for _, peer := range network_peers {
		fmt.Println(peer.String())
	}

	return network_peers
}

func (self *GWnode) registerTransaction(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		// transaction recived
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
				if !self.requestBlockMining() {
					// handle Probable JD absens, turn GW booting node into JD?
				}

			} else if sender := request.Header.Get("X-from-peer"); sender == "" {
				self.broadcastTransactionToPeers(transaction_data)
			} else {
				fmt.Println(" Got new transaction from", sender)
			}
		}
	} else {
		// syncing request recived:
		sender, err := ioutil.ReadAll(request.Body)
		if err == nil {
			fmt.Println("Syncing request sent by", string(sender))
			self.clearTransactions()
		}
	}
}

func (self *GWnode) run() {
	self.router.RegisterRoute(patriot_router.NewRoute("/register-transaction", true), self.registerTransaction)
	self.router.RegisterRoute(patriot_router.NewRoute("/peers", true), self.handlePeers)
	self.router.RegisterRoute(patriot_router.NewRoute("/blockchain", true), self.handleBlockchainDataReques)
	self.router.RegisterRoute(patriot_router.NewRoute("/new-block", true), self.addBlock)

	fmt.Println("Awaiting connections on port:", self.used_port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", self.used_port), self.router); err != nil {
		patriot_blockchain.LogFatal(err)
	}
}

func (self *GWnode) syncTransactions() {
	// clears transactions on the entire gw network
	var err error
	var client *http.Client = &http.Client{}
	var request *http.Request
	var request_body *bytes.Buffer
	var body_content []byte = []byte(fmt.Sprint(self.used_port))

	fmt.Print("sync done with:")
	for _, peer := range self.known_peers {
		request_body = bytes.NewBuffer(body_content)
		request, err = http.NewRequest("DELETE", fmt.Sprintf("http://%s:%s/register-transaction", peer.Host, peer.Port), request_body)
		if err == nil {
			client.Do(request)
		}
		if err != nil {
			fmt.Printf(", Skiping %s due to: %s", peer.String(), err.Error())
		} else {
			fmt.Printf(", %s", peer.String())
		}
	}
	fmt.Println("\n Done syncing")

	self.clearTransactions()
}

func main() {
	GW_NODE_NAME = os.Getenv("GW_NAME")

	fmt.Println("GWV:", os.Getenv("GWV"))
	var GW_node *GWnode = new(GWnode)
	GW_node.boot()

}
