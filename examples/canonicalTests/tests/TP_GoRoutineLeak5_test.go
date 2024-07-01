package main

import (
	"testing"
	"time"
)

func Test43(t *testing.T) {
	n43()
}
func n43() {
	c := make(chan int, 0)

	go func() {
		<-c
	}()

	go func() {
		c <- 1
	}()

	go func() {
		c <- 1
	}()

	time.Sleep(100 * time.Millisecond)
}
