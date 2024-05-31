
package main

import (
	"testing"
	"sync"
	"time"
)

// TP: possible send to closed channel recorded with once
func Test14(t *testing.T) {
	n14()
}
func n14() {
	c := make(chan int, 1)

	once := sync.Once{}

	go func() {
		once.Do(func() {
			c <- 1
		})
	}()

	go func() {
		time.Sleep(100 * time.Millisecond) // prevent actual send on closed channel
		once.Do(func() {
			// do nothing
		})
	}()

	time.Sleep(200 * time.Millisecond)
	close(c)
}
