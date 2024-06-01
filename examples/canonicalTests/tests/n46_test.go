package main

import (
	"testing"
	"time"
)

func Test46(t *testing.T) {
	n46()
}

// =============== Select Partner ===============
func n46() {
	c := make(chan int, 0)

	go func() {
		select {
		case c <- 1:
		}
	}()

	c <- 1

	time.Sleep(100 * time.Millisecond)
}
