package main

const genesis_hash = "0000000000000000000000000000000000000000000"

type Blockchain struct {
	blocks *List
}

func (self *Blockchain) genesis(data []byte) {
	if self.blocks == nil {
		var genesis_block *Block = CreateBlock(data, genesis_hash, 0)

		self.blocks = new(List)
		self.blocks.append(genesis_block)
	}
}

func (self *Blockchain) validateBlockchain() bool {
	var previous_hash string = genesis_hash
	for _, block := range self.blocks.toSlice() {
		ok, _ := block.(*Block).verifyChecksum()
		if block.(*Block).Previous == previous_hash || ok {
			return false
		}
	}
	return true
}
