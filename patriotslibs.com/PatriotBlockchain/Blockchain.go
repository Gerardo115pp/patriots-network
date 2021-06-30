package patriotblockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const GENESIS_HASH = "0000000000000000000000000000000000000000000"

type Blockchain struct {
	blocks *List
}

func (self *Blockchain) AddBlock(block *Block) error {
	if !self.blocks.IsEmpty() && block.Previous != self.blocks.root.NodeContent.(*Block).Hash {
		return fmt.Errorf("block previous hash doesnt match blockchain head hash")
	}
	ok, err := block.verifyChecksum()
	if ok {
		self.blocks.Push(block)
	}
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

	for h := range blocks_slice {
		self.blocks.Append(&(blocks_slice[(len(blocks_slice)-1)-h])) // loads the slice in reverse
	}
	_, err = self.ValidateBlockchain()
	return err
}

func (self *Blockchain) DataToJson() []byte {
	var data_fields []string
	data_fields = self.blocks.Map(func(ln *ListNode) string { return string(ln.NodeContent.(*Block).Data) })
	json_bytes, err := json.Marshal(data_fields)
	if err != nil {
		LogFatal(err)
	}
	return json_bytes
}

func (self *Blockchain) Genesis(data []byte) *Block {
	var genesis_block *Block
	if self.blocks == nil || self.blocks.length == 0 {
		genesis_block = CreateBlock(data, GENESIS_HASH, 0)
	}
	return genesis_block
}

func (self *Blockchain) Head() Block {
	return *(self.blocks.root.NodeContent.(*Block))
}

func (self *Blockchain) HeadHash() string {
	return self.blocks.root.NodeContent.(*Block).Hash
}

func (self *Blockchain) Load(filename string) error {
	var filedata []byte
	filedata, err := ioutil.ReadFile(fmt.Sprintf("%s.json", filename))
	if err != nil {
		return err
	}

	return self.FromBytes(filedata)
}

func (self *Blockchain) NewBlock(data []byte) *Block {
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
