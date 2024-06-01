package main

import (
	"testing"
	"time"
)

func Test22(t *testing.T) {
	n22()
}

// TN
func n22() {
	x := make(chan int)

	go func() {
		x <- 1
	}()

	go func() {
		x <- 1
	}()

	<-x
	<-x

	time.Sleep(300 * time.Millisecond)
}
