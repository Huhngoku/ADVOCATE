package main

import (
	"testing"
	"time"
)

func Test54(t *testing.T) {
	n54()
}

// leak because of select without possible partner
func n54() {
	c := make(chan int, 0)
	d := make(chan int, 0)

	go func() {
		select {
		case c <- 1:
		case d <- 1:
		}
	}()

	time.Sleep(200 * time.Millisecond)
}
