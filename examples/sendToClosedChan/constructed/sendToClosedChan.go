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

// FN = False negative
// FP = False positive
// TN = True negative
// TP = True positive

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
// NSC.
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

// TN: no send to closed channel because of once
func n11() {
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
func n12() {
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

const N = 12

func main() {

	runtime.InitAtomics(0)

	defer func() {
		runtime.DisableTrace()

		file_name := "constructed.log"
		os.Remove(file_name)
		file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		numRout := runtime.GetNumberOfRoutines()
		for i := 1; i <= numRout; i++ {
			cobufiChan := make(chan string)
			go func() {
				runtime.TraceToStringByIdChannel(i, cobufiChan)
				close(cobufiChan)
			}()
			for trace := range cobufiChan {
				if _, err := file.WriteString(trace); err != nil {
					panic(err)
				}
			}
			if _, err := file.WriteString("\n"); err != nil {
				panic(err)
			}
		}
	}()

	ns := [N]func(){n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12}

	for i := 0; i < N; i++ {
		ns[i]()
	}
}

/* Expected:
Possible send on closed channel:
	close: .../sendToClosedChan.go:129
	send : .../sendToClosedChan.go:121
Possible receive on closed channel:
	close: .../sendToClosedChan.go:129
	recv : .../sendToClosedChan.go:125
Possible send on closed channel:
	close: .../sendToClosedChan.go:143
	send : .../sendToClosedChan.go:138
Possible receive on closed channel:
	close: .../sendToClosedChan.go:143
	recv : .../sendToClosedChan.go:139
Receive on closed channel:
	close: .../sendToClosedChan.go:168
	recv : .../sendToClosedChan.go:165
*/
