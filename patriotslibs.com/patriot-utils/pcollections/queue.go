package pcollections

import "fmt"

type Qnode struct {
	next    *Qnode
	Content interface{}
}

func createQnode(content interface{}) *Qnode {
	var new_qnode *Qnode = new(Qnode)
	new_qnode.Content = content
	return new_qnode
}

type Queue struct {
	head *Qnode
	tail *Qnode
	len  uint
}

func (self *Queue) Clear() {
	self.head = nil
	self.tail = nil
	self.len = 0
}

func (self *Queue) Dequque() interface{} {
	if self.len == 0 {
		return nil
	} else if self.len == 1 {
		defer self.Clear()
		return self.head.Content
	}

	var current_node *Qnode = self.tail

	for current_node.next.next != nil {
		current_node = current_node.next
	}
	self.head = current_node
	current_node = self.head.next
	self.head.next = nil

	self.len--

	return current_node.Content
}

func (self *Queue) Enqueue(data interface{}) uint {
	var new_qnode *Qnode = createQnode(data)
	new_qnode.next = self.tail
	self.tail = new_qnode
	if self.head == nil {
		self.head = self.tail
	}
	self.len++
	return self.len
}

func (self *Queue) Len() uint {
	return self.len
}

func (self *Queue) String() string {
	var self_string string = "[tail]>"
	current_node := self.tail
	for current_node != nil {
		self_string += fmt.Sprintf(" %v", current_node.Content)
		if current_node.next != nil {
			self_string += ","
		}
		current_node = current_node.next
	}
	self_string += " <[head]"
	return self_string
}
