package main

import "fmt"

type Peer struct {
	host string `json: "host"`
	port string `json: "port"`
}

func (self *Peer) Json() string {
	return fmt.Sprintf("{\"host\": %s, \"port\": %s}", self.host, self.port)
}

func (self *Peer) String() string {
	return fmt.Sprintf("%s:%s", self.host, self.port)
}
