package main

/*
import "sync"
import "time"
import "fmt"

*/

import (
	"cobufi"
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

const N = 20

func main() {

	runtime.InitCobufi(0)
	defer cobufi.CreateTrace("constructedTest.log")

	// defer func() {
	// 	runtime.DisableTrace()

	// 	file_name := "constructed.log"
	// 	os.Remove(file_name)
	// 	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer file.Close()

	// 	numRout := runtime.GetNumberOfRoutines()
	// 	for i := 1; i <= numRout; i++ {
	// 		cobufiChan := make(chan string)
	// 		go func() {
	// 			runtime.TraceToStringByIdChannel(i, cobufiChan)
	// 			close(cobufiChan)
	// 		}()
	// 		for trace := range cobufiChan {
	// 			if _, err := file.WriteString(trace); err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 		if _, err := file.WriteString("\n"); err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }()

	ns := [N]func(){n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12,
		n13, n14, n15, n16, n17, n18, n19, n20}

	for i := 0; i < N; i++ {
		ns[i]()
	}
}

/*
==================== Summary ====================

-------------------- Critical -------------------
Possible send on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:129
	send : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:121
Possible send on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:143
	send : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:138
Possible send on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:180
	send : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:186
Possible send on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:181
	send : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:196
Possible send on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:313
	send : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:301
Possible send on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:360
	send : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:350
Possible send on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:436
	send : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:428
-------------------- Warning --------------------
Found concurrent Send on same channel:
	send: /home/erik/Uni/HiWi/CoBuFi-Go/go-patch/src/runtime/mgcsweep.go:279
	send : /home/erik/Uni/HiWi/CoBuFi-Go/go-patch/src/runtime/mgcscavenge.go:652
Possible receive on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:129
	recv : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:125
Possible receive on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:143
	recv : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:139
Receive on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:168
	recv : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:165
Possible receive on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:180
	recv : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:197
Possible receive on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:181
	recv : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:191
Receive on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:339
	recv : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:333
Possible receive on closed channel:
	close: /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:360
	recv : /home/erik/Uni/HiWi/CoBuFi-Go/examples/sendToClosedChan/constructed/sendToClosedChan.go:356

=================================================
Total runtime: 9.958979ms
=================================================
*/
