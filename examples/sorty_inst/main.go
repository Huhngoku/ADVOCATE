package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ErikKassubek/deadlockDetectorGo/src/dedego"
)

func main() {
	dedego.Init(20)
	defer dedego.RunAnalyzer()
	defer time.Sleep(time.Millisecond)
	l := 10000000
	input := make([]int, l)
	for i := 0; i < l; i++ {
		input[i] = rand.Intn(l) + 1
	}
	SortSlice(input)
	fmt.Println(input)
}

var dedegoFetchOrder = make(map[int]int)
