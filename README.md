# Analysis of Go programs for the automatic detection of potential concurrency bugs

This program is still under development and may return no or wrong results.

## What
We want to analyze concurrent Go programs to automatically find potential concurrency bug. For now we only support the search for potential send/receive on a closed channel, but we plan to expand the use cases in the future. 

## Recording
To analyze the program, we first need
to record it. To do this, we modify the go runtime
to automatically record a program while it runs. The modified runtime can 
be found in the `go-patch` directory. Running a program with this modified 
go runtime will create a trace of the program including 

- spawning of new routines
- atomic operations
- mutex operations
- channel operations
- select operations
- wait group operations
- once operations

The following is a short explanation about how to build and run the 
new runtime and create the trace. A full explanation of the created trace can be found in the 
`doc` directory. 

### Warning
The recording of atomic operations is only tested with `amd64`. For `arm64` an untested implementation exists. 

### How
The go-patch folder contains a modified version of the go runtime.
With this modified version it is possible to save a trace of the program.

To build the new runtime, run the 'make.bash' or 'make.bat' file in the 'src'
directory. This will create a 'bin' directory containing a 'go' executable.
This executable can be used as your new go environment e.g. with
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

	file_name := "tracer.log"
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


### Example
Let's create the trace for the following program:

```go
package main

import (
	"sync/atomic"
	"time"
)

func main() {
	c := make(chan int, 0)

	go func() {
		c <- 1
	}()

	go func() {
		<-c
	}()

	time.Sleep(10 * time.Millisecond)
	close(c)
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
	// ======= Preamble Start =======
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
	// ======= Preamble End =======

	c := make(chan int, 0)

	go func() {
		c <- 1  // line 48
	}()

	go func() {
		<-c  // line 52
	}()

	time.Sleep(10 * time.Millisecond)
	close(c)  // line 56
}
```

Running this leads to the following trace (indented lines are in the same line 
as the previous line, only for better readability):

```txt
G,1,2;G,2,3;G,3,4;C,4,9,1,R,f,1,2,.../go-patch/src/runtime/mgc.go:180;C,10,11,1,R,f,2,2,.../go-patch/src/runtime/mgc.go:181;G,12,5;C,13,13,2,C,f,0,0,.../go-patch/src/runtime/proc.go:256;G,14,6;G,15,7;G,16,8;C,21,21,4,C,f,0,0,.../main.go:56

C,7,8,1,S,f,2,2,.../go-patch/src/runtime/mgcsweep.go:279
C,5,6,1,S,f,1,2,.../go-patch/src/runtime/mgcscavenge.go:652


C,18,19,4,S,f,1,0,.../main.go:48
C,17,20,4,R,f,1,0,.../main.go:52
```
In this example the file paths are shortened for readability. In the real trace, the full path is given.

The trace includes both the concurrent operations of the program it self, as well
as internal operations used by the go runtime. An explanation of the trace 
file including the explanations for all elements can be found in the `doc`
directory.

## Analysis


We can now analyze the created file using the program in the `analyzer`
folder. For now we only support the search for potential send on a closed channel, but we plan to expand the use cases in the future. 

The analyzer can take the following command line arguments:

- -l [file]: path to the log file, default: ./trace.log
- -d [level]: output level, 0: silent, 1: results, 2: errors, 3: info, 4: debug, default: 2
- -b [buffer_size]: if the trace file is to big, it can be necessary to increase the size of the reader buffer. The size is given in MB, default: 25
- -f: if set, the analyzer assumes a fifo ordering of messages in the buffer of buffered channels. This is not part of the [Go Memory Mode](https://go.dev/ref/mem), but should follow from the implementation. For this reason, it is only an optional addition.
- -o [file_name]: set the name of the output file. If it is not set, or set to "", no output file will be created.
- -r: show the result immediately when found (default false)
- -s: do not show the summary at the end

If we assume the trace from our example is saved in file `trace.go` and run the analyzer with
```
./analyzer -f -l "trace.log" -o "result.log"
```
it will create the following result, show it in the terminal and print it into 
an `result.log` file: 
```txt
==================== Summary ====================

-------------------- Critical -------------------
Possible send on closed channel:
	close: .../main.go:56
	send: .../main.go:48
-------------------- Warning --------------------
Possible receive on closed channel:
	close: .../main.go:56
	recv: .../main.go:42
=================================================
Total runtime: Total runtime: 5.464833ms
=================================================
```
The send can cause a panic of the program, if it occurs. It is therefor an error message (in terminal red).

A receive on a closed channel does not cause a panic, but returns a default value. It can therefor be a desired behavior. For this reason it is only considered a warning (in terminal orange, can be silenced with -w).

