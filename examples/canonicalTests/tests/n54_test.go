package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n54()
}

// leak because of select without possible partner
func n54() {
	c := make(chan int, 0)

	go func() {
		select {
		case c <- 1:
		}
	}()

	time.Sleep(200 * time.Millisecond)
}
