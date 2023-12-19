package main

/*
import "sync"
import "time"
import "fmt"

*/

import (
	"advocate"
	"runtime"
	"sync"
	"time"
)

// FN = False negative
// FP = False positive
// TN = True negative
// TP = True positive

// =========== Send / Received to/from closed channel ===========

//////////////////////////////////////////////////////////////
// No send of closed due to (must) happens before relations.

// Synchronous channel.
// TN.
func n01() {
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
// TN
func n02() {
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
// TN
func n03() {
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
func n04() {
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

// TP send on closed
// TP recv on closed
func n05() {
	c := make(chan int)

	go func() {
		c <- 1
	}()

	go func() {
		<-c
	}()

	time.Sleep(300 * time.Millisecond) // prevent actual send on closed channel
	close(c)
}

// TP send on closed
// TP recv on closed
func n06() {
	c := make(chan int, 1)

	go func() {
		c <- 1
		<-c
	}()

	time.Sleep(300 * time.Millisecond) // prevent actual send on closed channel
	close(c)
}

// TN
func n07() {
	c := make(chan int)

	go func() {
		c <- 1
	}()

	<-c

	close(c)
}

// TP recv on closed
func n08() {
	c := make(chan int)

	go func() {
		time.Sleep(300 * time.Millisecond) // force actual recv on closed channel
		<-c
	}()

	close(c)
	time.Sleep(1 * time.Second) // prevent termination before receive
}

// TP send on closed
// TP recv on closed
func n09() {
	c := make(chan struct{}, 1)
	d := make(chan struct{}, 1)

	go func() {
		time.Sleep(300 * time.Millisecond) // prevent actual send on closed channel
		close(c)
		close(d)
	}()

	go func() {
		select {
		case c <- struct{}{}:
		default:
		}

		select {
		case <-d:
		default:
		}
	}()

	d <- struct{}{}
	<-c

	time.Sleep(1 * time.Second) // prevent termination before receive
}

// FN
func n10() {
	c := make(chan struct{}, 0)

	go func() {
		time.Sleep(200 * time.Millisecond) // prevent actual send on closed channel
		close(c)
	}()

	go func() {
		select {
		case c <- struct{}{}:
		default:
		}
	}()

	time.Sleep(300 * time.Millisecond) // make sure, that the default values are taken
}

// TP
func n11() {
	c := make(chan struct{}, 0)

	go func() {
		time.Sleep(200 * time.Millisecond) // prevent actual send on closed channel
		close(c)
	}()

	go func() {
		select {
		case c <- struct{}{}:
		default:
		}
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-c:
		default:
		}
	}()

	time.Sleep(300 * time.Millisecond) // make sure, that the default values are taken
}

// TN: no send to closed channel because of once
func n12() {
	c := make(chan int, 1)

	once := sync.Once{}

	go func() {
		once.Do(func() {
			c <- 1
		})
	}()

	go func() {
		once.Do(func() {
			close(c)
		})
	}()

	time.Sleep(100 * time.Millisecond)
}

// FN: potential send to closed channel not recorded because of once
func n13() {
	c := make(chan int, 1)

	once := sync.Once{}

	close(c)

	go func() {
		time.Sleep(100 * time.Millisecond) // prevent actual send on closed channel
		once.Do(func() {
			c <- 1
		})
	}()

	go func() {
		once.Do(func() {
			// do nothing
		})
	}()

	time.Sleep(100 * time.Millisecond)
}

// TP: potential send to closed channel recorded with once
func n14() {
	c := make(chan int, 1)

	once := sync.Once{}

	go func() {
		once.Do(func() {
			c <- 1
		})
	}()

	go func() {
		time.Sleep(100 * time.Millisecond) // prevent actual send on closed channel
		once.Do(func() {
			// do nothing
		})
	}()

	time.Sleep(200 * time.Millisecond)
	close(c)
}

// TN: no send possible
func n15() {
	c := make(chan int, 0)
	m := sync.Mutex{}

	go func() {
		time.Sleep(100 * time.Millisecond)
		t := m.TryLock()
		if t {
			c <- 1
			m.Unlock()
		}
	}()

	go func() {
		t := m.TryLock()
		if t {
			<-c
			m.Unlock()
		}
	}()

	time.Sleep(300 * time.Millisecond)
	close(c)
}

// TP
func n16() {
	c := make(chan int, 0)
	m := sync.Mutex{}

	go func() {
		t := m.TryLock()
		if t {
			c <- 1
			m.Unlock()
		}
	}()

	go func() {
		<-c
	}()

	time.Sleep(100 * time.Millisecond)
	close(c)
}

// FN
func n17() {
	c := make(chan int, 0)
	m := sync.Mutex{}

	go func() {
		time.Sleep(100 * time.Millisecond)
		t := m.TryLock()
		if t {
			c <- 1
			<-c
			m.Unlock()
		}
	}()

	m.Lock()
	time.Sleep(300 * time.Millisecond)
	close(c)
	m.Unlock()

	time.Sleep(100 * time.Millisecond)
}

// TP
func n18() {
	ch := make(chan int, 1)
	var g sync.WaitGroup

	g.Add(1)

	func() {
		g.Done()
		ch <- 1
	}()

	g.Wait()
	time.Sleep(100 * time.Millisecond)
	close(ch)
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

// TP
func n20() {
	ch := make(chan int, 2)

	f := func() {
		ch <- 1
	}

	go func() {
		f()
	}()

	time.Sleep(200 * time.Millisecond)
	close(ch)
}

// ============== Concurrent recv on same channel ==============

// TP
func n21() {
	x := make(chan int)

	go func() {
		<-x
	}()

	go func() {
		<-x
	}()

	x <- 1
	x <- 1

	time.Sleep(300 * time.Millisecond)
}

// TN
func n22() {
	x := make(chan int)

	go func() {
		x <- 1
	}()

	go func() {
		x <- 1
	}()

	<-x
	<-x

	time.Sleep(300 * time.Millisecond)
}

const n = 22

func main() {
<<<<<<< Updated upstream:examples/potentialBugs/constructed/potentialBugs.go

	runtime.InitAdvocate(0)
	defer advocate.CreateTrace("constructed.log")
=======
	if true {
		// init tracing
		runtime.InitAdvocate(0)
		defer advocate.CreateTrace("trace_constructed.log")
	} else {
		// init replay
		trace := advocate.ReadTrace("trace_constructed.log")
		runtime.EnableReplay(trace)
		defer runtime.WaitForReplayFinish()
	}
>>>>>>> Stashed changes:examples/constructedMain/constructed/potentialBugs.go

	ns := [n]func(){n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12,
		n13, n14, n15, n16, n17, n18, n19, n20, n21, n22}

	for i := 0; i < n; i++ {
		ns[i]()
<<<<<<< Updated upstream:examples/potentialBugs/constructed/potentialBugs.go
=======
		println("Done: ", i+1, " of ", n)
		// time.Sleep(300 * time.Millisecond)
>>>>>>> Stashed changes:examples/constructedMain/constructed/potentialBugs.go
	}
}
