package main

import (
	"testing"
)

func Test24(t *testing.T) {
	n24()
}

// TN: No concurrent send on same channel
func n24() {
	x := make(chan int)

	go func() {

		<-x
		<-x
	}()

	x <- 1
	x <- 1
}
