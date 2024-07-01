package main

import (
	"testing"
	"time"
)

func Test41(t *testing.T) {
	n41()
}
func n41() {
	c := make(chan int, 0)

	go func() {
		close(c)
	}()

	time.Sleep(100 * time.Millisecond)
}
