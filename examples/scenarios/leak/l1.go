package main

import (
	"time"
)

func main() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	go func() {
		<-ch //A
	}()
	go func() {
		<-ch //B
	}()

	time.Sleep(1 * time.Second)
}

// Example Leak on unbuffered channel with possible partner
// A or B might block forever, causing a leak. However they have a potential partner that can unblock them.
