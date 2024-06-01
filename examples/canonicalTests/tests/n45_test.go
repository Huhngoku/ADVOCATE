package main

import (
	"sync"
	"testing"
	"time"
)

func Test45(t *testing.T) {
	n45()
}

func n45() {
	m := sync.Mutex{}

	go func() {
		m.Lock()
		m.Lock()
	}()

	time.Sleep(100 * time.Millisecond)
}
