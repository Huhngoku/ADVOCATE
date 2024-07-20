package main

import (
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Done()
	wg.Done()
}

// Example Negative WaitGroup Count
