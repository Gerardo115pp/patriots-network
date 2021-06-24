package patriotblockchain

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ProofWork(block *Block) bool {
	return strings.HasPrefix(block.Hash, "000")
}

func Clamp(x int, min int, max int) int {
	if x > max {
		return x
	} else if x < min {
		return min
	}
	return x
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func LogFatal(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func StringToUint(s string) uint {
	return uint(StringToInt(s))
}

func StringToInt(s string) int {
	new_int, err := strconv.Atoi(s)
	if err == nil {
		return new_int
	} else {
		panic(err)
	}
}

func StringToUint64(s string) uint64 {
	new_uint64, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		return new_uint64
	} else {
		panic(err)
	}
}

func ShaAsInt64(s string) uint64 {
	var hash_bytes [20]byte = sha1.Sum([]byte(s))
	return binary.BigEndian.Uint64(hash_bytes[:])
}
