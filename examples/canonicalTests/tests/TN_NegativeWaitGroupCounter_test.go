package main

import (
	"sync"
	"testing"
)

func Test25(t *testing.T) {
	n25()
}

// ============== Negative wait counter (Add before done) ==============
// no possible negative wait counter
func n25() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Done()
}
