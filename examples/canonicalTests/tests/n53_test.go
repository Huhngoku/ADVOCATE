package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n53()
}

// leak because of select with possible partner
func n53() {
	c := make(chan int, 0)

	go func() {
		<-c
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)

		select {
		case c <- 1:
		}
	}()

	c <- 1

	time.Sleep(200 * time.Millisecond)
}
