package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type User struct {
	User_name  string `json:"user_name"`
	hash       string
	Color      string `json:"color"`
	Is_logged  bool   `json:"status"`
	connection *websocket.Conn
	chat       *Server
}

func (self *User) init(user_name string, color string, hash string, chat *Server) {
	self.User_name = user_name
	self.hash = hash
	self.Color = color
	self.Is_logged = false
	self.chat = chat
}

func (self *User) read() {
	for {
		if _, bmsg, err := self.connection.ReadMessage(); err != nil {
			fmt.Printf("Error on message: %v\n", err)
			break
		} else {
			self.chat.chat_message_channel <- newMessage(string(bmsg), self.User_name, self.Color, false, false)
		}
	}
	if DEBUG {
		fmt.Printf("\n\nLogging %s out\n\n", self.User_name)
	}
	self.Is_logged = false
	self.chat.user_leaves <- self
}

func (self *User) write(message *Message) {
	msg, _ := json.Marshal(message)
	if err := self.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
		fmt.Printf("Error while writing a message to user '%s': %v", self.User_name, err)
	}
}
