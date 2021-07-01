package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	patriot_blockchain "github.com/Gerardo115pp/PatriotLib/PatriotBlockchain"
	patriot_router "github.com/Gerardo115pp/PatriotLib/PatriotRouter"
)

var DEBUG bool = false

type GWagent struct {
	Host         string `json:"host"`
	Lisenting_on string `json:"port"`
	connections  uint
	identifier   uint64
}

type XZagent struct {
	Host         string `json:"host"`
	Lisenting_on string `json:"port"`
}

type JD struct {
	router       *patriot_router.Router
	GWagents     map[string]*GWagent
	XZagents     map[string]*XZagent
	port         string
	host         string
	mining_state bool
}

func (self *JD) awaitXZ(request_data *bytes.Buffer, xz *XZagent, resolution_channel chan []byte) {
	var client *http.Client = &http.Client{}
	request, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/mine", xz.Host, xz.Lisenting_on), request_data) // we should try to implement a context request to cancel requests before attempting to reimplement Transport
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Failed to connect to xz node %s:%s, reason: %s\n", xz.Host, xz.Lisenting_on, err.Error())
		return
	}
	response_data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		patriot_blockchain.LogFatal(err)
	}

	resolution_channel <- response_data
}

func (self *JD) broadcastGWConnection(new_gw *GWagent) {
	gw_data, err := json.Marshal(new_gw)
	if err != nil {
		patriot_blockchain.LogFatal(err)
	}

	var request_body *bytes.Buffer

	var gw_count uint = 1
	for _, gw := range self.GWagents {
		fmt.Printf("\rBoadcasting (%d/%d)", gw_count, len(self.GWagents)-1)
		if gw.Host == new_gw.Host && gw.Lisenting_on == new_gw.Lisenting_on {
			continue
		}

		request_body = bytes.NewBuffer(gw_data)
		request_address := fmt.Sprintf("http://%s:%s/peers", gw.Host, gw.Lisenting_on)

		castDebug(fmt.Sprintf("gw data: %s\n", string(gw_data)))

		http.Post(request_address, "application/json", request_body)
	}
	fmt.Println("done")
}

func (self *JD) broadcastMiningRequest(block_data []byte) {
	request_body := bytes.NewBuffer(block_data)
	var block_resolution_channel chan []byte = make(chan []byte)
	if len(self.XZagents) > 0 {
		self.mining_state = true

		for _, xz_miner := range self.XZagents {
			go self.awaitXZ(request_body, xz_miner, block_resolution_channel)
		}

		var new_block_data []byte
		elapsed_xz := 0
		for elapsed_xz < len(self.XZagents) {
			new_block_data = <-block_resolution_channel
			if len(new_block_data) == 0 {
				elapsed_xz++
			} else {
				break
			}
		}
		fmt.Println("New block created.")
		self.castBlockToGws(new_block_data)
	} else {
		fmt.Println("No miners avaliable, stacking request.")
	}
}

func (self *JD) castBlockToGws(block_data []byte) {
	var http_client *http.Client = &http.Client{}
	var request_body *bytes.Buffer

	for _, gw := range self.GWagents {
		request_body = bytes.NewBuffer(block_data)
		request, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/new-block", gw.Host, gw.Lisenting_on), request_body)
		if err != nil {
			fmt.Println("While casting new block:", err.Error())
			continue
		}
		http_client.Do(request)
	}
}

func (self *JD) composeGWagentFromRequest(request *http.Request) *GWagent {
	var gw_agent *GWagent = new(GWagent)
	gw_agent.Host, _, _ = net.SplitHostPort(request.RemoteAddr)
	gw_agent.Lisenting_on = request.URL.Port()
	gw_agent.connections = 0
	if self.port == "" {
		self.port = "80"
	}
	return gw_agent
}

func (self *JD) composeJson(key string, value string) []byte {
	return []byte(fmt.Sprintf("{\"%s\": %s}", key, value))
}

func (self *JD) composeResponse(response_value string) []byte {
	return self.composeJson("response", response_value)
}

func (self *JD) composeError(error_value string) []byte {
	return self.composeJson("error", error_value)
}

func (self *JD) getXZcode() string {
	if len(self.XZagents) < 1000 {
		var xz_code string = fmt.Sprint(len(self.XZagents) + 1)
		xz_code = "0000"[:4-len(xz_code)%4] + xz_code
		return xz_code
	} else {
		return "a001"
	}
}

func (self *JD) getLazyestGW() *GWagent {
	var lazyiest *GWagent
	var lowest_connection_count uint = 9999
	for code, gw := range self.GWagents {
		if gw.connections < lowest_connection_count {
			lazyiest = self.GWagents[code]
			lowest_connection_count = gw.connections
		}
	}
	return lazyiest
}

func (self *JD) handleAllGWs(response http.ResponseWriter, request *http.Request) {
	var gw_identifier uint64 = 0
	if param_id := request.URL.Query().Get("code"); param_id != "" {
		gw_identifier = patriot_blockchain.StringToUint64(param_id)
	}
	fmt.Println(gw_identifier, "requested peers")

	if len(self.GWagents) == 0 {
		response.WriteHeader(200)
		fmt.Fprint(response, "[]")
		return
	}

	var gw_agents []*GWagent = make([]*GWagent, 0)
	var gw_agents_connected []byte
	for _, gw := range self.GWagents {
		if gw.identifier != gw_identifier {
			gw_agents = append(gw_agents, gw)
		} else {
			fmt.Println("Excluding client GW")
		}
	}

	gw_agents_connected, err := json.Marshal(gw_agents)
	castDebug(fmt.Sprintf("%s, len=%d", string(gw_agents_connected), len(gw_agents_connected)))
	if err != nil {
		patriot_blockchain.LogFatal(err)
	}
	response.WriteHeader(200)
	response.Header().Set("Content-Type", "application/json")
	response.Write(gw_agents_connected)
}

func (self *JD) handleGW(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		if len(self.GWagents) == 0 {
			response.WriteHeader(404)
			fmt.Fprint(response, "no GW registred")
		} else {
			fmt.Println("Getting lazyest GW")
			var gw_node *GWagent = self.getLazyestGW()
			var gw_data []byte
			gw_data, err := json.Marshal(gw_node)
			if err != nil {
				response.WriteHeader(500)
				fmt.Fprint(response, "Server error sorry for the inconvinience")
				fmt.Println("An error occurred while tying to serve GW node:", err.Error())
				return
			}
			response.WriteHeader(200)
			response.Write(gw_data)
		}
	case http.MethodPost:
		// register a gw agent
		post_form, err := self.parseFormAsMap(request)
		if err != nil {
			response.WriteHeader(400)
			response.Write(self.composeError("wrong form data"))
			return
		}
		var lisenting_port string = post_form["port"]
		if lisenting_port != "" {
			var gw_agent *GWagent = self.composeGWagentFromRequest(request)
			gw_agent.Lisenting_on = lisenting_port

			gw_remote_address := fmt.Sprintf("%s:%s", gw_agent.Host, gw_agent.Lisenting_on)
			gw_agent.identifier = patriot_blockchain.ShaAsInt64(gw_remote_address)
			if _, exists := self.GWagents[gw_remote_address]; exists {
				fmt.Printf("Recived a duplicated GW request from %s\n", gw_remote_address)
			} else {

				fmt.Printf("New GW agent with host '%s:%s'\n", gw_agent.Host, gw_agent.Lisenting_on)
				self.GWagents[gw_remote_address] = gw_agent
				if len(self.GWagents) > 1 {
					fmt.Println("Broadcasting gw connection")
					defer self.broadcastGWConnection(gw_agent)
				}

			}

			fmt.Println("Network peers:", len(self.GWagents))

			response.WriteHeader(200)
			fmt.Fprintf(response, fmt.Sprint(gw_agent.identifier))
		} else {
			fmt.Println("Error lisenting port was:", lisenting_port)
			response.WriteHeader(400)
			response.Write(self.composeError("missing code"))
		}
	}
}

func (self *JD) handleXZ(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPut:
		fmt.Println("New block mining request")
		block_data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println("Request error:", err.Error())
			response.WriteHeader(400)
			fmt.Fprintf(response, "wrong request format")
			return
		}

		go self.broadcastMiningRequest(block_data)

		response.WriteHeader(200)
		response.Write([]byte("ok"))
	case http.MethodPost:
		// register a gw agent
		form_values, err := self.parseFormAsMap(request)
		if err != nil {
			fmt.Println("XZ node registration error due to a bad request")
			response.WriteHeader(400)
			return
		}

		listening_port := form_values["port"]
		fmt.Println("form_values:", form_values)

		if listening_port != "" {
			var xz_agent *XZagent = new(XZagent)
			var xz_code string = self.getXZcode()
			xz_agent.Lisenting_on = listening_port
			xz_agent.Host, _, _ = net.SplitHostPort(request.RemoteAddr)

			fmt.Printf("New XZ agent with host '%s:%s' and code %s \n", xz_agent.Host, xz_agent.Lisenting_on, xz_code)

			self.XZagents[xz_code] = xz_agent
			response.WriteHeader(200)
		} else {
			response.WriteHeader(400)
			response.Write(self.composeError("missing code"))
		}
	}
}

func (self *JD) parseFormAsMap(request *http.Request) (map[string]string, error) {
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

func (self *JD) run() {
	self.router.RegisterRoute(patriot_router.NewRoute(`/GW-all`, true), self.handleAllGWs)
	self.router.RegisterRoute(patriot_router.NewRoute(`/GW`, true), self.handleGW)
	self.router.RegisterRoute(patriot_router.NewRoute(`/XZ`, true), self.handleXZ)

	fmt.Println("Lisinting on port:", self.port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", self.host, self.port), self.router); err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func createJD() *JD {
	var new_jd_node *JD = new(JD)
	new_jd_node.host = "0.0.0.0"
	new_jd_node.port = "4000"
	if custom_port := os.Getenv("JD_PORT"); custom_port != "" {
		new_jd_node.port = custom_port
	}
	new_jd_node.GWagents = make(map[string]*GWagent)
	new_jd_node.XZagents = make(map[string]*XZagent)
	new_jd_node.router = patriot_router.CreateRouter()
	new_jd_node.mining_state = false
	return new_jd_node
}
