package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n52()
}

// leaking because of chan without possible partner
func n52() {
	c := make(chan int, 0)

	go func() {
		c <- 1
	}()

	time.Sleep(200 * time.Millisecond)
}
