package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	go func() {
		<-ch //A
		//code will not be executed because of leak
		fmt.Println("Never executed")
	}()
	time.Sleep(1 * time.Second)
}

// Example Leak on unbuffered channel without possible partner
// In this case the goroutine blocks at A and there is no other goroutine that can unblock it.
