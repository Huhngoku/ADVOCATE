package main

import (
	"sync"
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n55()
}

// leak because of wait group
func n55() {
	w := sync.WaitGroup{}

	go func() {
		w.Add(1)
		w.Wait()
	}()

	time.Sleep(200 * time.Millisecond)
}
