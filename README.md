# Patch version of the Go programming language to record concurrent event

## What
We want to statically analyze concurrent Go programs. For this we first have
to record the program. To do this, we adapt the go runtime and compiler
to automatically record a program while it runs. The modified runtime can 
be found in the `go-patch` directory. Running a program with this modified 
go runtime will create a trace of the program including 

- spawning of new routines
- atomic operations
- mutex operations
- channel operations
- select operations
- wait group operations

The following is a short explanation about how to build and run the 
new runtime and create tge trace. A full explanation of the created trace can be found in the 
`doc` directory. 

## Warning
The recording of atomic operations only works with `amd64`.

## How
The go-patch folder contains a modified version of the go compiler and runtime.
With this modified version it is possible to save a trace of the program.

To build the new runtime, run the 'all.bash' or 'all.bat' file in the 'src'
directory. This will create a 'bin' directory containing a 'go' executable.
This executable can be used as your new go envirement e.g. with
`./go run main.go` or `./go build`.

It is necessary to set the GOROOT environment variable to the path of `go-patch`, e.g. with 
```
export GOROOT=$HOME/dedego/go-patch/
```

To create a trace, add

```go
runtime.InitAtomics(0)

defer func() {
	runtime.DisableTrace()

	file_name := "dedego.log"
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
```

at the beginning of the main function.
Also include the following imports 
```go
runtime
io/ioutil
os
```

Autocompletion often includes "std/runtime" instead of "runtime". Make sure to include the correct one.

For some reason, `fmt.Print` and similar can lead to `fatal error: schedule: holding lock`. In this case increase the argument in `runtime.InitAtomics(0)`
until the problem disappears.

After that run the program with `./go run main.go` or `./go build && ./main`,
using the new runtime.

## Example
Let's create the trace for the following program:

```go
package main

import (
	"sync/atomic"
	"time"
)

func main() {
	c := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			c <- 1
		}
		var a int32
		var b int32
		atomic.AddInt32(&a, 1)
		atomic.AddInt32(&a, 1)
		atomic.AddInt32(&b, 1)
	}()

	for i := 0; i < 3; i++ {
			<-c
	}

	time.Sleep(time.Second)
}
```

After adding the preamble, we get 

```go
package main

import (
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {
	runtime.InitAtomics(0)

	defer func() {
		runtime.DisableTrace()

		file_name := "dedego.log"
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

	c := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			c <- 1
		}
		var a int32
		var b int32
		atomic.AddInt32(&a, 1)
		atomic.AddInt32(&a, 1)
		atomic.AddInt32(&b, 1)
	}()

	for i := 0; i < 3; i++ {
			<-c
	}

	time.Sleep(time.Second)
}
```

Running this leads to the following trace (indented lines are in the same line 
as the previous line, only for better readability):

```txt
G,1,2;G,2,3;G,3,4;C,4,9,1,R,t,1,2,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgc.go:180;C,10,11,1,R,t,2,2,1,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgc.go:181;G,12,5;C,13,13,2,C,t,0,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/proc.go:256;G,14,6;G,15,7;C,16,20,4,R,t,1,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:56;C,21,22,4,R,t,2,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:56;C,23,43,4,R,t,3,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:56

C,5,6,1,S,t,1,2,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgcsweep.go:279
C,7,8,1,S,t,2,2,0,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgcscavenge.go:652

C,27,33,3,R,t,1,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:201;C,34,35,3,R,t,2,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:201;C,36,41,3,R,t,3,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:201;C,42,0,3,R,f,4,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:201
C,17,18,4,S,t,1,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:46;C,19,24,4,S,t,2,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:46;C,25,26,4,S,t,3,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:46;A,28,824634294432,A;C,29,30,3,S,t,1,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedego_atomic.go:104;A,31,824634294432,A;C,32,37,3,S,t,2,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedego_atomic.go:104;A,38,824634294436,A;C,39,40,3,S,t,3,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedego_atomic.go:104
```

The trace includes both the concurrent operations of the program it self, as well
as internal operations used by the go runtime. An explanation of the trace 
file including the explanations for all elements can be found in the `doc`
directory.
