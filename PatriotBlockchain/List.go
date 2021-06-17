package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const DEBUG = false
const SEPARATOR = "<,>"

type Content interface {
	getId() uint
	toString() string
	String() string
	toJson() []byte
}

type ConditionalFunc func(Content) bool

type ListNode struct {
	NodeContent Content
	Prev        *ListNode
	Next        *ListNode
}

type List struct {
	length            int
	root              *ListNode
	iterator_position uint
}

func (self *List) append(content Content) int {
	var new_node *ListNode = new(ListNode)
	new_node.NodeContent = content

	if self.length == 0 {
		self.root = new_node
	} else {
		current_node := self.root
		for current_node.Next != nil {
			current_node = current_node.Next
		}
		current_node.Next = new_node
		new_node.Prev = current_node
	}
	self.length++
	return self.length
}

func (self *List) clear() {
	self.root = nil
	self.length = 0
}

func (self *List) exists(content Content) Content {
	if self.root != nil {
		var current_node *ListNode = self.root
		for current_node != nil {

			if current_node.NodeContent.String() == content.String() {
				return current_node.NodeContent
			}
			current_node = current_node.Next
		}
	}
	return nil
}

func (self *List) filter(conditionf ConditionalFunc) *List {
	var new_list *List = new(List)
	if self.root != nil {
		current_node := self.root
		for current_node != nil {
			if conditionf(current_node.NodeContent) {
				new_list.append(current_node.NodeContent)
			}
			current_node = current_node.Next
		}
	}
	return new_list
}

func (self *List) get(position int) Content {
	if position >= self.length {
		panic(fmt.Errorf("index %d in list of length %d is out of range", position, self.length))
	}

	var current_node *ListNode = self.root
	for h := 0; h != position && current_node != nil; h++ {
		current_node = current_node.Next
	}
	return current_node.NodeContent
}

func (self *List) mapFunc(callback func(*ListNode) string) []string {
	var current_node *ListNode = self.root
	var map_results []string
	for current_node != nil {
		map_results = append(map_results, callback(current_node))
		current_node = current_node.Next
	}
	return map_results
}

func (self *List) push(c Content) int {
	var new_node *ListNode = new(ListNode)
	new_node.NodeContent = c
	self.root.Prev = new_node
	new_node.Next = self.root
	self.root = new_node
	self.length++
	return self.length
}

func (self *List) pop() Content {
	if self.root != nil {
		var current_node *ListNode = self.root
		if self.length == 1 {
			self.root = nil
			self.length = 0
			return current_node.NodeContent
		}

		for current_node.Next != nil {
			current_node = current_node.Next
		}
		current_node.Prev.Next = nil

		self.length--
		return current_node.NodeContent
	}
	return nil
}

func (self *List) remove(content Content) {
	var indirect **ListNode = &self.root
	for (*indirect) != nil && !((*indirect).NodeContent.getId() == content.getId()) {
		indirect = &((*indirect).Next)
	}
	if ((*indirect).Next) != nil {
		((*indirect).Next).Prev = (*indirect).Prev
	}
	*indirect = ((**indirect).Next)
	self.length--
}

func (self *List) searchBy(condition func(Content) bool) Content {
	for current_node := self.root; current_node != nil; current_node = current_node.Next {
		if condition(current_node.NodeContent) {
			return current_node.NodeContent
		}
	}
	return nil
}

func (self *List) save(file_name string) bool {
	var serialized_string string
	var current_node *ListNode = self.root
	for true {
		serialized_string += current_node.NodeContent.toString()
		current_node = current_node.Next
		if current_node != nil {
			serialized_string += SEPARATOR
		} else {
			break
		}
	}

	return ioutil.WriteFile(file_name, []byte(serialized_string), 0666) == nil
}

func (self *List) toJson() string {
	var json_content []string
	for current_node := self.root; current_node != nil; current_node = current_node.Next {
		json_content = append(json_content, string(current_node.NodeContent.toJson()))
	}
	return fmt.Sprintf("[%s]", strings.Join(json_content, ","))
}

func (self *List) toString() (rstring string) {
	if self.root == nil {
		return "{Empty-List}"
	}

	node_iter := self.root
	node_position := 0
	for {
		p_node := "nil"
		n_node := "nil"
		if node_position != 0 {
			p_node = node_iter.Prev.NodeContent.toString()
		}
		if node_position != (self.length - 1) {
			n_node = node_iter.Next.NodeContent.toString()
		}

		rstring += fmt.Sprintf("(%s<<(%s, p%d)>>%s)", p_node, node_iter.NodeContent.toString(), node_position, n_node)
		node_position++
		node_iter = node_iter.Next
		if node_iter == nil {
			break
		}
		rstring += " - "
	}
	return rstring
}

func (self *List) toSlice() []Content {
	var list_slice []Content
	for current_node := self.root; current_node != nil; current_node = current_node.Next {
		list_slice = append(list_slice, current_node.NodeContent)
	}
	return list_slice
}
