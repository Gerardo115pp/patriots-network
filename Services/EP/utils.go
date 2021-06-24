package main

import (
	"math/rand"
	"time"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getRandom() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func getRandomInt64() int64 {
	return getRandom().Int63()
}
