package patriotblockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

const GENESIS_HASH = "0000000000000000000000000000000000000000000"

type Blockchain struct {
	blocks *List
	head   *Block
}

func (self *Blockchain) AddBlock(block *Block) error {
	if !self.blocks.IsEmpty() && block.Previous != self.head.Hash {
		return fmt.Errorf("block previous hash doesnt match blockchain head hash")
	}
	ok, err := block.verifyChecksum()
	if ok {
		self.blocks.Append(block)
	}
	self.head = block
	return err
}

func (self *Blockchain) Equals(other *Blockchain) bool {
	var self_footprint uint64 = ShaAsInt64(self.blocks.Json())
	var other_footprint uint64 = ShaAsInt64(other.blocks.Json())
	return self_footprint == other_footprint
}

func (self *Blockchain) FromBytes(blockchain_bytes []byte) (err error) {

	var blocks_slice []Block
	err = json.Unmarshal(blockchain_bytes, &blocks_slice)
	if err != nil {
		return err
	}

	sort.Slice(blocks_slice, func(i, j int) bool {
		return blocks_slice[i].Number < blocks_slice[j].Number
	})

	for h := range blocks_slice {
		self.AddBlock(&(blocks_slice[h]))
	}

	_, err = self.ValidateBlockchain()
	return err
}

func (self *Blockchain) DataToJson() []byte {
	var data_fields []string = make([]string, 0)
	self.blocks.Map(func(ln *ListNode) string {
		data_fields = append(data_fields, ln.NodeContent.(*Block).Data...)
		return ""
	}) // appends all the transactions of each block to the data_fields
	json_bytes, err := json.Marshal(data_fields)
	if err != nil {
		LogFatal(err)
	}
	return json_bytes
}

func (self *Blockchain) Genesis(data []string) *Block {
	var genesis_block *Block
	if self.blocks == nil || self.blocks.length == 0 {
		genesis_block = CreateBlock(data, GENESIS_HASH, 0)
	}
	return genesis_block
}

func (self *Blockchain) Head() Block {
	return *(self.head)
}

func (self *Blockchain) HeadHash() string {
	return self.Head().Hash
}

func (self *Blockchain) Load(filename string) error {
	var filedata []byte
	filedata, err := ioutil.ReadFile(fmt.Sprintf("%s.json", filename))
	if err != nil {
		return err
	}

	return self.FromBytes(filedata)
}

func (self *Blockchain) NewBlock(data []string) *Block {
	var previous_hash string = self.blocks.root.NodeContent.(*Block).Hash
	var block_num uint = uint(self.blocks.length)
	return CreateBlock(data, previous_hash, block_num)
}

func (self *Blockchain) Save(filename string) {
	var json_data string = self.blocks.Json()
	ioutil.WriteFile(fmt.Sprintf("%s.json", filename), []byte(json_data), 0604)
}

func (self *Blockchain) ToBytes() []byte {
	return []byte(self.blocks.Json())
}

func (self *Blockchain) ValidateBlockchain() (bool, error) {
	var previous_hash string = GENESIS_HASH
	for _, block := range self.blocks.Slice() {
		_, err := block.(*Block).verifyChecksum()
		if err != nil {
			return false, err
		}
		if block.(*Block).Previous != previous_hash {
			return false, fmt.Errorf("block hash '%s' doesnt match next block prev '%s'", previous_hash, block.(*Block).Previous)
		}
		previous_hash = block.(*Block).Hash
	}
	return true, nil
}

func (self *Blockchain) List() *List {
	return self.blocks
}

func (self *Blockchain) Length() int {
	return self.blocks.length
}

func CreateBlockchain() *Blockchain {
	var new_blockchain *Blockchain = new(Blockchain)
	new_blockchain.blocks = new(List)
	return new_blockchain
}
