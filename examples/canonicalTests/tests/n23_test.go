package main

import (
	"testing"
)

func Test23(t *testing.T) {
	n23()
}

// TN: No concurrent recv on same channel
func n23() {
	x := make(chan int, 2)

	go func() {
		x <- 1
		x <- 1
	}()

	<-x
	<-x
}
