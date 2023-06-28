package main

import (
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	runtime.EnableTrace()
	defer func() {
		file_name := "dedego.log"
		os.Remove(file_name)
		file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		runtime.DisableTrace()
		numRout := runtime.GetNumberOfRoutines()
		for i := 0; i < numRout; i++ {
			c := make(chan string)
			go func() {
				runtime.TraceToStringByIdChannel(i, c)
				close(c)
			}()
			for trace := range c {
				if _, err := file.WriteString(trace); err != nil {
					panic(err)
				}
			}
			if _, err := file.WriteString("\n"); err != nil {
				panic(err)
			}
		}
		file.Close()
		println("With: ", time.Since(start).Seconds())
	}()

	c := make(chan int)
	d := make(chan int)
	m := sync.Mutex{}

	max := 100000

	for i := 0; i < max; i++ {
		time.Sleep(1 * time.Millisecond)
		go func() {
			m.Lock()
			c <- 1
			m.Unlock()
		}()
	}

	for i := 0; i < max; i++ {
		select {
		case <-c:
		case <-d:
		}
	}

	println("Without: ", time.Since(start).Seconds())
}
