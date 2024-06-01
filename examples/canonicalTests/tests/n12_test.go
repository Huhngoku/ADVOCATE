package main

import (
	"testing"
	"time"
	"sync"
)

// TN: no send to closed channel because of once
func Test12(t *testing.T) {
	n12()
}
func n12() {
	c := make(chan int, 1)

	once := sync.Once{}

	go func() {
		once.Do(func() {
			c <- 1
		})
	}()

	go func() {
		once.Do(func() {
			close(c)
		})
	}()

	time.Sleep(100 * time.Millisecond)
}
