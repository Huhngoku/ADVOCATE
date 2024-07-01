package main

import (
	"testing"
)

func Test01(t *testing.T) {
	TN_SynchronousChannel()
}

// Synchronous channel.
// TN.
func TN_SynchronousChannel() {
	x := make(chan int)
	ch := make(chan int, 1)

	go func() {
		ch <- 1
		x <- 1
	}()

	<-x
	close(ch)
}
