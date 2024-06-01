
package main

import (
	"testing"
	"time"
	"sync"
)

// FN: possible send to closed channel not recorded because of once
func Test13(t *testing.T) {
	n13()
}
func n13() {
	c := make(chan int, 1)

	once := sync.Once{}

	close(c)

	go func() {
		time.Sleep(100 * time.Millisecond) // prevent actual send on closed channel
		once.Do(func() {
			c <- 1
		})
	}()

	go func() {
		once.Do(func() {
			// do nothing
		})
	}()

	time.Sleep(100 * time.Millisecond)
}
