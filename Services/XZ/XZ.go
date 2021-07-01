package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	patriot_blockchain "github.com/Gerardo115pp/PatriotLib/PatriotBlockchain"
	patriot_router "github.com/Gerardo115pp/PatriotLib/PatriotRouter"
)

const JS_ADDRESS = "192.168.1.79:4000"

type BlockData struct {
	Transactions []string `json:"transactions"`
	BlockNum     int      `json:"block_num"`
	Previous     string   `json:"previous"`
}

type XZ struct {
	port   int
	router *patriot_router.Router
}

func (self *XZ) boot() {
	// registers to JD as a miner
	var client *http.Client = new(http.Client)
	var request_body *bytes.Buffer = bytes.NewBufferString(fmt.Sprintf("{\"port\": \"%d\"}", self.port))
	var request *http.Request

	request, _ = http.NewRequest("POST", fmt.Sprintf("http://%s/XZ", JS_ADDRESS), request_body)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		patriot_blockchain.LogFatal(err)
	}

	if response.StatusCode == 200 {
		self.run()
	} else {
		fmt.Println("Rejected by JD.")
	}

}

func (self *XZ) run() {
	self.router.RegisterRoute(patriot_router.NewRoute("/mine", true), HandleMineBlockRequest)

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", self.port), self.router); err != nil {
		patriot_blockchain.LogFatal(err)
	}
}

func (self XZ) String() string {
	return fmt.Sprintf("XZ miner working on %d", self.port)
}

func HandleMineBlockRequest(response http.ResponseWriter, request *http.Request) {
	var block_data *BlockData = new(BlockData)
	bytes_data, err := ioutil.ReadAll(request.Body)
	fmt.Println("Request:", string(bytes_data))
	if err == nil {
		err = json.Unmarshal(bytes_data, block_data)
	}

	if err != nil {
		fmt.Println("Mine request error:", err.Error())
		response.WriteHeader(400)
		fmt.Fprint(response, "request format error")
		return
	}

	block_mined := MineBlock(block_data)

	response_data, err := json.Marshal(block_mined)

	response.WriteHeader(200)
	response.Write(response_data)
}

func MineBlock(block_data *BlockData) *patriot_blockchain.Block {
	var new_block *patriot_blockchain.Block = patriot_blockchain.CreateBlock(block_data.Transactions, block_data.Previous, uint(block_data.BlockNum))
	for !patriot_blockchain.ProofWork(new_block) {
		new_block.Nonce++
		new_block.Hash = patriot_blockchain.HashBlock(new_block)
		fmt.Printf("\rTrying with nonce '%d' -> %s", new_block.Nonce, new_block.Hash)
	}
	fmt.Println("\nDone")
	return new_block
}

func CreateXZ(port string) *XZ {
	var new_xz *XZ = new(XZ)
	new_xz.port = patriot_blockchain.StringToInt(port)
	new_xz.router = patriot_router.CreateRouter()

	return new_xz
}

func main() {
	var param_port string = os.Getenv("XZPORT")
	if param_port != "" {
		var xz *XZ = CreateXZ(param_port)
		fmt.Println(xz)
		xz.boot()
	} else {
		fmt.Println("No XZPORT variable was empty.")
		os.Exit(1)
	}
}
