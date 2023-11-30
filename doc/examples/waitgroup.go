package main

import "sync"
import "time"
import "fmt"

/*

Add can act like Done

Counter not allowed to become negative


*/

func negativeCounterPanic() {
	var g sync.WaitGroup

	g.Add(1)

	func() {
		g.Add(-1)
	}()

	func() {
		g.Done()
	}()

	time.Sleep(1 * time.Second)

	g.Wait()
	fmt.Printf("done")

}

func negativeCounterPanic2() {
	var g sync.WaitGroup

	g.Add(1)

	func() {
		g.Done()
	}()

	func() {
		g.Done()
	}()

	time.Sleep(1 * time.Second)

	g.Wait()
	fmt.Printf("done")

}

func addInBetween() {
	var g sync.WaitGroup

	g.Add(1)

	func() {
		g.Done()
	}()

	func() {
		g.Add(1)
	}()

	func() {
		g.Done()
	}()

	time.Sleep(1 * time.Second)

	g.Wait()
	fmt.Printf("done")

}

// Happens-before relation

// Add can act like Done
// go-race does not issue a race here.
// Hence,
//
//	A < W  and D < W
func addActsLikeDone() {
	var g sync.WaitGroup
	shared := 1

	g.Add(3)

	read := func(x int) {}

	func() {
		read(shared)
		g.Add(-2) // A
	}()

	func() {
		read(shared)
		g.Done() // D
	}()

	time.Sleep(1 * time.Second)

	g.Wait() // W
	shared = 3
	fmt.Printf("done")

}

// out of sync reset via add leads to race with wait
func resetViaAdd() {
	var g sync.WaitGroup

	g.Add(3)

	func() {
		g.Add(-2)
	}()

	func() {
		g.Done()
	}()

	time.Sleep(1 * time.Second)

	go func() {
		time.Sleep(2 * time.Second)
		g.Add(1)
	}()
	g.Wait()

	time.Sleep(1 * time.Second)

	fmt.Printf("done")

	time.Sleep(3 * time.Second)

	go func() {
		g.Done()
	}()

	g.Wait()

	fmt.Printf("done2")

}

func main() {

	// negativeCounterPanic3()

	// addActsLikeDone()

	resetViaAdd()

}
