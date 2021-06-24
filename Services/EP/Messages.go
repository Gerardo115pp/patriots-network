package main

type Message struct {
	Uuid    int64  `json:"uuid"`
	Body    string `json:"body"`
	Sender  string `json:"sender"`
	Control bool   `json:"is_control_message"`
	Color   string `json:"color"`
	IsFile  bool   `json:"is_file"`
}

func newMessage(msg string, sender string, color string, is_control_message bool, is_file bool) *Message {
	var new_message *Message = new(Message)
	new_message.Body = msg
	new_message.Sender = sender
	new_message.Uuid = getRandomInt64()
	new_message.Color = color
	new_message.Control = is_control_message
	new_message.IsFile = is_file
	return new_message
}
