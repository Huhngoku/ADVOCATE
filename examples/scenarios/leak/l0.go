package main

import (
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	go func() {
		ch <- 1
	}()
	go func() {
		close(ch)
	}()

	time.Sleep(1 * time.Second)
}
