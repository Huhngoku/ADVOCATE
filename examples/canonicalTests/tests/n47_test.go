package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n47()
}

func n47() {
	c := make(chan int, 0)
	d := make(chan int, 0)

	go func() {
		c <- 1
	}()

	go func() {
		d <- 1
	}()

	select {
	case <-c:
	case <-d:
	}

	close(c)

	time.Sleep(100 * time.Millisecond)
}
