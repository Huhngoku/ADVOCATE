
package main

import (
	"testing"
	"time"
	"sync"
)

func Test19(t *testing.T) {
	n19()
}
// FN
func n19() {
	ch := make(chan int, 1)
	m := sync.Mutex{}

	go func() {
		m.Lock()
		ch <- 1
		time.Sleep(100 * time.Millisecond)
		m.Unlock()
	}()

	time.Sleep(100 * time.Millisecond)
	if m.TryLock() {
		close(ch)
		m.Unlock()
	}
	time.Sleep(200 * time.Millisecond)
}
