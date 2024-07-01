package main

import (
	"testing"
	"time"
)

func Test40(t *testing.T) {
	n40()
}

func n40() {
	c := make(chan int, 0)

	go func() {
		<-c
	}()

	time.Sleep(100 * time.Millisecond)
}
