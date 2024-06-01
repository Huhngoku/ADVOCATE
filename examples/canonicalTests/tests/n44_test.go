package main

import (
	"sync"
	"testing"
	"time"
)

func Test44(t *testing.T) {
	n44()
}
func n44() {
	w := sync.WaitGroup{}

	go func() {
		time.Sleep(100 * time.Millisecond)
		w.Wait()
	}()

	w.Add(1)

	time.Sleep(100 * time.Millisecond)
}
