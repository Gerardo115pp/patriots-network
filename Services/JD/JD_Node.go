package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	patriot_router "github.com/Gerardo115pp/PatriotLib/PatriotRouter"
)

type GWagent struct {
	Host         string `json:"host"`
	Lisenting_on string `json:"port"`
	connections  uint
	relaible     bool
}

type XZagent struct {
	Host         string `json:"host"`
	Lisenting_on string `json:"port"`
}

type JD struct {
	router   *patriot_router.Router
	GWagents map[string]*GWagent
	XZagents map[string]*XZagent
	port     string
	host     string
}

func (self *JD) composeGWagentFromRequest(request *http.Request) *GWagent {
	var gw_agent *GWagent = new(GWagent)
	gw_agent.Host, _, _ = net.SplitHostPort(request.RemoteAddr)
	gw_agent.Lisenting_on = request.URL.Port()
	gw_agent.relaible = true
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

func (self *JD) handleGW(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		if len(self.GWagents) == 0 {
			response.WriteHeader(404)
			fmt.Fprint(response, "no GW registred")
		} else {
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
		var gw_code string = request.FormValue("code")
		var lisinting_port string = request.FormValue("port")
		if gw_code != "" && lisinting_port != "" {
			var gw_agent *GWagent = self.composeGWagentFromRequest(request)
			gw_agent.Lisenting_on = lisinting_port
			fmt.Printf("New GW agent with host '%s:%s'\n", gw_agent.Host, gw_agent.Lisenting_on)
			self.GWagents[gw_code] = gw_agent
		} else {
			response.WriteHeader(400)
			response.Write(self.composeError("missing code"))
		}
	}
}

func (self *JD) handleXZ(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		// register a gw agent
		var lisinting_port string = request.FormValue("port")
		if lisinting_port != "" {
			var xz_agent *XZagent = new(XZagent)
			var xz_code string = self.getXZcode()
			xz_agent.Lisenting_on = lisinting_port
			xz_agent.Host, _, _ = net.SplitHostPort(request.RemoteAddr)

			fmt.Printf("New XZ agent with host '%s:%s' and code %s \n", xz_agent.Host, xz_agent.Lisenting_on, xz_code)

			self.XZagents[xz_code] = xz_agent
		} else {
			response.WriteHeader(400)
			response.Write(self.composeError("missing code"))
		}
	}
}

func (self *JD) run() {
	self.router.RegisterRoute(patriot_router.NewRoute(`/GW`, false), self.handleGW)
	self.router.RegisterRoute(patriot_router.NewRoute(`/XZ`, false), self.handleXZ)

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
	return new_jd_node
}
