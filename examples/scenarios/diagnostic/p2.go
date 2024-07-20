package main

import "time"

func main() {
	ch := make(chan int, 1)
	ch <- 1
	go func() {
		<-ch
	}()
	go func() {
		close(ch)
	}()
	time.Sleep(1 * time.Second)
}

// Example Rec on closed
