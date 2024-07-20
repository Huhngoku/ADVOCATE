package main

import "time"

func main() {
	ch := make(chan int, 1)
	go func() {
		ch <- 1
	}()
	go func() {
		close(ch)
	}()
	time.Sleep(1 * time.Second)
}

// Example Send on closed
