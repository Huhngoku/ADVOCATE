package main

import (
	"sync"
	"testing"
	"time"
)

func Test02(t *testing.T) {
	n02()
}

// Wait group
// TN
func n02() {
	ch := make(chan int, 1)
	var g sync.WaitGroup

	g.Add(1)

	func() {
		ch <- 1
		g.Done()
	}()

	g.Wait()
	close(ch)

	time.Sleep(100 * time.Millisecond)

}
