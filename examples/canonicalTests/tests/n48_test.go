package main

import (
	"testing"
)

func Test48(t *testing.T) {
	n48()
}

func n48() {
	c := make(chan int, 0)
	d := make(chan int, 0)

	go func() {
		c <- 1
	}()

	select {
	case <-c:
	case <-d:
	}
}
