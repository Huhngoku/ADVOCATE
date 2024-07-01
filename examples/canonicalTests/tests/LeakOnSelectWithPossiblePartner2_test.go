package main

import (
	"testing"
	"time"
)

func Test53(t *testing.T) {
	n53()
}

// leak because of select with possible partner
func n53() {
	c := make(chan int, 0)
	d := make(chan int, 0)

	go func() {
		<-d
	}()

	go func() {
		<-c
	}()

	go func() {
		time.Sleep(300 * time.Millisecond)

		select {
		case c <- 1:
		case d <- 1:
		}
	}()

	c <- 1
	d <- 1

	time.Sleep(800 * time.Millisecond)
}
