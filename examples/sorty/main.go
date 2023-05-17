package main

import (
	"fmt"
	"math/rand"
)

func main() {
	l := 10000000
	input := make([]int, l)
	for i := 0; i < l; i++ {
		input[i] = rand.Intn(l) + 1
	}
	SortSlice(input)
	fmt.Println(input)
}
