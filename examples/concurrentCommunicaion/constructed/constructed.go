package main

/*
import "sync"
import "time"
import "fmt"

*/

import (
	"os"
	"runtime"
	"time"
)

// NSC = No send on closed due to must-happens before relations
// FN = False negative
// FP = False positive
// TP = True positive

func n1() {
	x := make(chan int)

	go func() {
		x <- 1
	}()

	go func() {
		x <- 1
	}()

	<-x
	<-x
}

// Wait group
// NSC.
func n2() {
	x := make(chan int, 2)

	go func() {
		x <- 1
		x <- 1
	}()

	<-x
	<-x
}

func n3() {
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

// Wait group
// NSC.
func n4() {
	x := make(chan int)

	go func() {

		<-x
		<-x
	}()

	x <- 1
	x <- 1
}

const N = 4

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
			advocateChan := make(chan string)
			go func() {
				runtime.TraceToStringByIdChannel(i, advocateChan)
				close(advocateChan)
			}()
			for trace := range advocateChan {
				if _, err := file.WriteString(trace); err != nil {
					panic(err)
				}
			}
			if _, err := file.WriteString("\n"); err != nil {
				panic(err)
			}
		}
	}()

	ns := [N]func(){n1, n2, n3, n4}

	for i := 0; i < N; i++ {
		ns[i]()
	}
}
