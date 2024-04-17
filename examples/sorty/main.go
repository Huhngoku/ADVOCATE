package main

import (
	"advocate"
	"math/rand"
)

func main() {
	if true {
		// init tracing
		advocate.InitTracing(-1)
		defer advocate.Finish()
	} else {
		// init replay
		advocate.EnableReplay()
		defer advocate.WaitForReplayFinish()
	}

	l := 10000000
	input := make([]int, l)
	rand.Seed(1) // added to create same sequence in replay
	for i := 0; i < l; i++ {
		input[i] = rand.Intn(l) + 1
	}
	SortSlice(input)
	_ = input
}
