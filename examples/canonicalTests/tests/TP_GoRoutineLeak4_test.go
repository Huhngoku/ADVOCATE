package main

import (
	"testing"
	"time"
)

func Test42(t *testing.T) {
	n42()
}

func n42() {
	c := make(chan int, 0)

	go func() {
		<-c
	}()

	go func() {
		c <- 1
	}()

	go func() {
		<-c
	}()

	time.Sleep(100 * time.Millisecond)
}
