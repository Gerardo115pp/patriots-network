package patriotblockchain

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Block struct {
	Number   uint   `json:"block_num"`
	Previous string `json:"previous"`
	Data     []byte `json:"data"`
	Nonce    uint   `json:"nonce"`
	Hash     string `json:"hash"`
}

func (self *Block) getId() uint {
	return self.Number
}

func (self *Block) String() string {
	return string(self.toJson())
}

func (self *Block) toString() string {
	return string(self.Data)
}

func (self *Block) toJson() []byte {
	data, err := json.Marshal(self)
	if err != nil {
		logFatal(err)
	}
	return data
}

func (self *Block) verifyChecksum() (bool, error) {
	if self.Hash != "" {
		if !ProofWork(self) {
			return false, fmt.Errorf("Didnt passed proof of work")
		}

		if len(self.Previous) != len(self.Hash) {
			return false, fmt.Errorf("Length of previous and length of hash didnt match")
		}

		var block_hash string
		if block_hash = HashBlock(self); block_hash != self.Hash {
			return false, fmt.Errorf("Content hash didnt match")
		}
		return true, nil
	}
	return false, nil
}

func CreateBlock(data []byte, previous_hash string, block_num uint) *Block {
	var new_block *Block = new(Block)

	new_block.Previous = previous_hash
	new_block.Number = block_num
	new_block.Data = data
	new_block.Nonce = 3

	new_block.Hash = HashBlock(new_block)
	return new_block
}

func HashBlock(block *Block) string {
	var block_footprint string = fmt.Sprintf("%d%s%x", block.Nonce, block.Previous, block.Data)
	hasher := sha256.New()
	_, err := hasher.Write([]byte(block_footprint))
	if err != nil {
		logFatal(err)
	}
	return base64.RawStdEncoding.EncodeToString(hasher.Sum(nil))
}
