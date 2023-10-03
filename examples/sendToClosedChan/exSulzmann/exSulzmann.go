package main

/*
import "sync"
import "time"
import "fmt"

*/

import (
	"os"
	"runtime"
	"sync"
	"time"
)

// NSC = No send on closed due to must-happens before relations
// FN = False negative
// FP = False positive

//////////////////////////////////////////////////////////////
// No send of closed due to (must) happens before relations.

// Synchronous channel.
// NSC.
func n1() {
	x := make(chan int)
	ch := make(chan int, 1)

	go func() {
		ch <- 1
		x <- 1
	}()

	<-x
	close(ch)

}

// Wait group
// NSC.
func n2() {
	ch := make(chan int, 1)
	var g sync.WaitGroup

	g.Add(1)

	func() {
		ch <- 1
		g.Done()
	}()

	g.Wait()
	close(ch)

}

// Once
// NSC.
func n3() {
	var once sync.Once
	ch := make(chan int, 1)
	setup := func() {
		ch <- 1
	}

	once.Do(setup)
	close(ch)

}

// RWMutex
// FN.
/*

T1 -> T2 -> T3 due to sleep statements

RU2 and RU1 sync with L

 => send <HB close


If we reorder critical sections,
we encounter send on closed.

*/
func n4() {
	var m sync.RWMutex
	ch := make(chan int, 1)

	// T1
	go func() {
		m.RLock()
		ch <- 1
		m.RUnlock() // RU1

	}()

	// T2
	go func() {
		time.Sleep(300 * time.Millisecond)
		m.RLock()
		m.RUnlock() // RU2

	}()

	// T3
	time.Sleep(1 * time.Second)
	m.Lock() // L
	close(ch)
	m.Unlock()

}

const N = 4

func main() {

	runtime.InitAtomics(0)

	defer func() {
		runtime.DisableTrace()

		file_name := "trace.log"
		os.Remove(file_name)
		file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		numRout := runtime.GetNumberOfRoutines()
		for i := 1; i <= numRout; i++ {
			dedegoChan := make(chan string)
			go func() {
				runtime.TraceToStringByIdChannel(i, dedegoChan)
				close(dedegoChan)
			}()
			for trace := range dedegoChan {
				if _, err := file.WriteString(trace); err != nil {
					panic(err)
				}
			}
			if _, err := file.WriteString("\n"); err != nil {
				panic(err)
			}
		}
	}()

	// ns := [N]func(){n1, n2, n3, n4}

	// for i := 0; i < N; i++ {
	// 	ns[i]()
	// }

	n4()
}
