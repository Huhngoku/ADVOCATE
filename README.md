# AdvocateGo: Automated Detection and Verification Of Concurrency bugs through Analysis of Trace recordings of program Executions in Go."

> [!WARNING]
> This program is still under development and may return no or wrong results.

## What
We want to analyze concurrent Go programs to automatically find potential concurrency bug. The different analysis scenarios can be found in `doc/Analysis.md`.

We also implement a trace replay mechanism, to replay a trace as recorded.

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

> [!WARNING]
> The recording of atomic operations is only tested with `amd64`. For `arm64` an untested implementation exists. 

### How
The go-patch folder contains a modified version of the go runtime.
With this modified version it is possible to save a trace of the program.

To build the new runtime, run the `make.bash` or `make.bat` file in the `src`
directory. This will create a `bin` directory containing a `go` executable.
This executable can be used as your new go environment e.g. with
`./go run main.go` or `./go build`. Please make sure, that the program expects 
go version 1.21 or earlier.

<!-- WARNING: It can currently happen, that `make.bash` command result in a `fatal error: runtime: releaseSudog with non-nil gp.param`. It can normally be fixed by just running `make.bash` again. I'm working on fixing it. -->

It is necessary to set the GOROOT environment variable to the path of `go-patch`, e.g. with 
```
export GOROOT=$HOME/ADVOCATE/go-patch/
```

To create a trace, add

```go
advocate.InitTracing(0)
defer advocate.Finish()
```

at the beginning of the main function.
Also include the following imports 
```go
advocate
```


In some cases, we can get a `fatal error: schedule: holding lock`. In this case increase the argument in `runtime.InitAtomics(0)` until the problem disappears.

After that run the program with `./go run main.go` or `./go build && ./main`,
using the new runtime.


### Example
Let's create the trace for the following program:

```go
package main

import (
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
	"advocate"
	"time"
)

func main() {
	// ======= Preamble Start =======
		advocate.InitTracing(0)
		defer advocate.Finish()
	// ======= Preamble End =======

	c := make(chan int, 0)

	go func() {
		c <- 1  // line 17
	}()

	go func() {
		<-c  // line 21
	}()

	time.Sleep(10 * time.Millisecond)
	close(c)  // line 25
}
```

Running this create a `trace` folder containing one trace file for each routine.
The non-empty trace files contain the following trace:

```txt
G,20,16,/home/.../main.go:16;G,21,17,/home/.../main.go:20;C,26,26,4,C,f,0,0,/home/.../main.go:25

C,23,24,4,S,f,1,0,/home/.../main.go:17

C,22,25,4,R,f,1,0,/home/.../main.go:21
```
In this example the file paths are shortened for readability. In the real trace, the full path is given.

The trace includes both the concurrent operations of the program it self, as well
as internal operations used by the go runtime. An explanation of the trace 
file including the explanations for all elements can be found in the `doc`
directory.

## Analysis and Reorder

We can now analyze the created file using the program in the `analyzer`
folder. For now we only support the search for potential send on a closed channel, but we plan to expand the use cases in the future.
The analyzer can also create a new reordered trace, in which a detected possible bug actually occurs. This new trace can then used in the replay, to confirm and simplify the removal of the bug.

> [!WARNING]
> The Reorder is still in development and may result in incorrect traces

The analyzer can take the following command line arguments:

- -c	Ignore happens before relations of critical sections (default false)
- -d int
    	Debug Level, 0 = silent, 1 = errors, 2 = info, 3 = debug (default 1) (default 1)
- -f	Assume a FIFO ordering for buffered channels (default false)
- -i int
    	Index of the result to use for the reordered trace file. Only needed if -n is set. 1 based (default -1)
- -n	Create a reordered trace file from a given analysis result without running the analysis. -r and -i are required. If not set, a rewritten trace can be created from the current analysis results
- -p	Do not print the results to the terminal (default false). Automatically set -x to true
- -r string
    	Path to where the result file should be saved. If not set, it is saved in the trace folder
- -t string
    	Path to the trace folder to analyze or rewrite
  -w	Do not print warnings (default false)
  -x	Do not ask to create a reordered trace file after the analysis (default false)

If we assume the trace from our example is saved in file `trace.go` and run the analyzer with
```
./analyzer -x -t /trace
```
it will create the following result, show it in the terminal and print it into an `result_readable.log` file: 
```txt
==================== Summary ====================

-------------------- Critical -------------------
1 Possible send on closed channel:
	close: /home/.../main.go:25@26
	send : /home/.../main.go:17@23
-------------------- Warning --------------------
2 Possible receive on closed channel:
	close: /home/.../main.go:25@26
	recv : /home/.../main.go:21@22
```
The send can cause a panic of the program, if it occurs. It is therefor an error message (in terminal red).

A receive on a closed channel does not cause a panic, but returns a default value. It can therefor be a desired behavior. For this reason it is only considered a warning (in terminal orange, can be silenced with -w).


## Trace Replay
The trace replay reruns a given program as given in the recorded trace. Please be aware, 
that only the recorded elements are considered for the trace replay. This means, that 
the order of non-recorded operations between two or more routines can still very. 

<!-- The implementation of the trace replay is not finished yet. The following is a short overview over the current state.
- order enforcement for most elements.
	- The operations are started in same global order as in the recorded trace. 
	- This is not yet implemented for the spawn of new routines and atomic operations
- correct channel partner
	- Communication partner of (most) channel operations are identical to the partners in the trace. For selects this cannot be guarantied yet. -->

### How
To start the replay, add the following header at the beginning of the 
main function:

```go
advocate.EnableReplay()
defer advocate.WaitForReplayFinish()
```

Also include the following import:
```go
"advocate"
```
Now the program can be run with the modified go routine, identical to the recording of the trace (remember to export the new gopath). 

It is important that the program is not changed between recording and replay.
This is especially true for the positions of operations on the code. For this 
reason it can be beneficial to add the following header, instead the two separate 
ones for recording and replay:

```go
if true {
	// init tracing
	advocate.InitTracing(0)
	defer advocate.Finish()
} else {
	// init replay
	advocate.EnableReplay()
	defer advocate.WaitForReplayFinish()
}
```

With changing `true` to `false` one can switch between recording and replay.

### Warning:
It is the users responsibility of the user to make sure, that the input to 
the program, including e.g. API calls are equal for the recording and the 
tracing. Otherwise the replay is likely to get stuck.

Do not change the program code between trace recording and replay. The identification of the operations is based on the file names and lines, where the operations occur. If they get changed, the program will most likely block without terminating. If you need to change the program, you must either rerun the trace recording or change the effected trace elements in the recorded trace.
This also includes the adding of the replay header. Make sure, that it is already in the program (but commented out), when you run the recording.
