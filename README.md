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
new runtime. A full explanation of the created trace can be found in the 
`doc` directory. 

## Warning
The modified runtime is currently only implemented and tested fot `amd64`. 
For `arm64` an untested implementation exists, but there are no guaranties, that
this implementation is runnable.

## How
The go-patch folder contains a modified version of the go compiler and runtime.
With this modified version it is possible to save a trace as described in
'Trace Structure'.

To build the new runtime, run the 'all.bash' or 'all.bat' file in the 'src'
directory. This will create a 'bin' directory containing a 'go' executable.
This executable can be used as your new go envirement e.g. with
`./go run main.go` or `./go build`.

It is necessary to set the GOROOT environment variable to the path of the `./go` executable, e.g. with 
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
		close(c)
		var a int32
		var b int32
		atomic.AddInt32(&a, 1)
		atomic.AddInt32(&a, 1)
		atomic.AddInt32(&b, 1)
	}()

	for a := range c {
		_ = a
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
		file.Close()

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
		close(c)
		var a int32
		var b int32
		atomic.AddInt32(&a, 1)
		atomic.AddInt32(&a, 1)
		atomic.AddInt32(&b, 1)
	}()

	for a := range c {
		_ = a
	}
	time.Sleep(time.Second)
}
```

Running this leads to the following trace (indented lines are in the same line 
as the previous line, only for better readability):

```

G,1,2;G,2,3;G,3,4;C,1,4,9,R,e,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgc.go:180;C,1,10,11,R,e,2,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgc.go:181;G,12,5;C,2,13,13,C,o,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/proc.go:256;G,14,6;G,15,7;C,4,16,20,R,e,1,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:56;C,4,21,22,R,e,2,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:56;C,4,23,30,R,e,3,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:56;C,4,31,32,R,e,4,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:56

C,1,5,6,S,e,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgcsweep.go:279
C,1,7,8,S,e,2,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgcscavenge.go:652

C,3,33,34,R,e,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:196;C,3,35,42,R,e,2,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:196;C,3,43,44,R,e,3,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:196;C,3,45,0,R,o,4,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:196
C,4,17,18,S,e,1,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:46;C,4,19,24,S,e,2,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:46;C,4,25,26,S,e,3,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:46;C,4,27,27,C,o,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:48;A,28,824634851496;C,3,29,36,S,e,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedego_atomic.go:40;A,37,824634851496;C,3,38,39,S,e,2,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedego_atomic.go:40;A,40,824634851500;C,3,41,46,S,e,3,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedego_atomic.go:40


```

The trace includes both the concurrent operations of the program it self, as well
as internal operations used by the go runtime. The elements from the
program are in file .../worst/main.go
