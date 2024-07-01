package main

import (
	"sync"
	"testing"
	"time"
)

func Test18(t *testing.T) {
	n18()
}

// TP
func n18() {
	ch := make(chan int, 1)
	var g sync.WaitGroup

	g.Add(1)

	func() {
		g.Done()
		ch <- 1
	}()

	time.Sleep(100 * time.Millisecond)
	g.Wait()
	close(ch)
}
