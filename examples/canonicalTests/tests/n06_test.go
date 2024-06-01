package main

import (
	"testing"
	"time"
)

func Test06(t *testing.T) {
}

// TP send on closed
// TP recv on closed
func n06() {
	c := make(chan int, 1)

	go func() {
		c <- 1
		<-c
	}()

	time.Sleep(300 * time.Millisecond) // prevent actual send on closed channel
	close(c)
}
