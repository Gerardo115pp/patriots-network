package main

import "fmt"

type Peer struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (self *Peer) Json() string {
	return fmt.Sprintf("{\"host\": \"%s\", \"port\": \"%s\"}", self.Host, self.Port)
}

func (self *Peer) String() string {
	return fmt.Sprintf("%s:%s", self.Host, self.Port)
}
