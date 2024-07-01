package main

import (
	"testing"
	"time"
)

func Test51(t *testing.T) {
	n51()
}

// =============== Leaking Channels ===============

// leaking because of chan with possible partner
func n51() {
	c := make(chan int, 0)

	go func() {
		c <- 1
		println(1)
	}()

	go func() {
		c <- 1
		println(2)
	}()

	<-c
	time.Sleep(200 * time.Millisecond)
}
