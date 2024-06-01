
package main

import (
	"testing"
	"time"
	"sync"
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

	g.Wait()
	time.Sleep(100 * time.Millisecond)
	close(ch)
}
