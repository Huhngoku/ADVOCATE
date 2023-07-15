# Patch version of the Go Programming language for Dedego

## How
The go-patch folder contains a modified version of the go compiler and runtime.
With this modified version it is possible to save a trace as described in
'Trace Structure'.

To build the new runtime, run the 'all.bash' or 'all.bat' file in the 'src'
directory. This will create a 'bin' directory containing a 'go' executable.
This executable can be used as your new go envirement e.g. with
`./go run main.go` or `./go build`.

In some cases it is necessary to set the GOROOT environment variable to the 
path of the `./go` executable, e.g. with 
```
export GOROOT=$HOME/dedego/go-patch/
```

To create a trace, add

```go
import (
  "runtime",
  "io/ioutil"
  "os"
)

runtime.InitAtomics()

defer func() {
	runtime.DisableTrace()

	file_name := "dedego.log"
	os.Remove(file_name)
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	numRout := runtime.GetNumberOfRoutines()
	for i := 0; i <= numRout; i++ {
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
	file.Close()

}()

```

Autocompletion often includes "std/runtime" instead of "runtime". Make sure to 
include the correct one.

After that run the program with `./go run main.go` or `./go build && ./main`,
using the new runtime.

## Trace structure

The following is the structure of the trace T in BNF.
```
T := L\nta | ""                                                 (trace)
t := L\nt  | ""                                                 (trace without atomics)
a := "" | {A";"}A                                               (trace of atomics)
L := "" | {E";"}E                                               (routine local trace)
E := G | M | W | C | S                                          (trace element)
G := "G,"tpre","id                                              (element for creation of new routine)
A := "A,"tpre","addr                                            (element for atomic operation)
M := "M,"tpre","tpost","id","rw","opM","exec","suc","pos        (element for operation on sync (rw)mutex)
W := "W,"tpre","tpost","id","opW","exec","delta","val","pos     (element for operation on sync wait group)
C := "C,"tpre","tpost","id","opC","exec","oId","pos             (element for operation on channel)
S := "S,"tpre","tpost","id","cases","exec","chosen","oId","pos  (element for select)
tpre := ℕ                                                       (timer when the operation is started)
tpost := ℕ                                                      (timer when the operation has finished)
addr := ℕ                                                       (pointer to the atomic variable, used as id)
id := ℕ                                                         (unique id of the underling object)
rw := "R" | "-"                                                 ("R" if the mutex is an RW mutex, "-" otherwise)
opM := "L" | "LR" | "T" | "TR" | "U" | "UR"                     (operation on the mutex, L: lock, LR: rLock, T: tryLock, TR: tryRLock, U: unlock, UR: rUnlock)
opW := "A" | "W"                                                (operation on the wait group, A: add (delta > 0) or done (delta < 0), W: wait)
opC := "S" | "R" | "C"                                          (operation on the channel, S: send, R: receive, C: close)
exec := "e" | "f"                                               (e: the operation was fully executed, o: the operation was not fully executed, e.g. a mutex was still waiting at a lock operation when the program was terminated or a channel never found an communication partner)
suc := "s" | "f"                                                (the mutex lock was successful ("s") or it failed ("f", only possible for try(r)lock))
pos := file":"line                                              (position in the code, where the operation was executed)
file := 𝕊                                                       (file path of pos)
line := ℕ                                                       (line number of pos)
delta := ℕ                                                      (change of the internal counter of wait group, normally +1 for add, -1 for done)
val := ℕ                                                        (internal counter of the wait group after the operation)
oId := ℕ                                                        (identifier for an communication on the channel, the send and receive (or select) that have communicated share the same oId)
cases := case | {case"."}case                                   (list of cases in select, seperated by .)
case := cId""("r" | "s") | "d"                                  (case in select, consisting of channel id and "r" for receive or "s" for send. "d" shows an existing default case)  
cId := ℕ                                                        (id of channel in select case)
chosen := ℕ0 | "-1"                                             (index of the chosen case in cases, -1 for default case)    
```

Info: 
- \n: newline
- ℕ: natural number not including 0
- ℕ0: natural number including 0
- 𝕊: string containing 1 or more characters
- The tracer contains a global timer for all routines that is incremented every time an timer element (tpre/tpost) is recorded.

## Changed files
Added files:

- src/runtime/dedego_routine.go
- src/runtime/dedego_trace.go
- src/runtime/dedego_util.go
- src/runtime/internal/atomic/dedego_atomic.go

Changed files (marked with DEDEGO-ADD):

- src/cmd/cgo/internal/test/testx.go
- src/runtime/proc.go
- src/runtime/runtime2.go
- src/runtime/chan.go
- src/runtime/select.go
- src/runtime/internal/atomic/doc.go
- src/runtime/internal/atomic/atomic_amd64.go
- src/runtime/internal/atomic/atomic_amd64.s
- src/sync/mutex.go
- src/sync/rwmutex.go
- src/sync/waitgroup.go

Disabled Tests (files contain disabled tests, marked with DEDEGO-REMOVE_TEST): 

- src/cmd/cgo/internal/test/cgo_test.go
- src/cmd/dist/test.go
- src/cmd/go/stript_test.go
- src/cmd/compile/internal/types2/sizeof_test.go
- src/context/x_test.go
- src/crypto/internal/nistec/nistec_test.go
- src/crypto/tls/tls_test.go
- src/go/build/deps_test.
- src/go/types/sizeof_test.go
- src/internal/intern/inter_test.go
- src/log/slog/text_handler_test.go
- src/net/netip/netip_test.go
- src/runtime/crash_cgo_test.go
- src/runtime/sizeof_test.go
- src/runtime/align_test.go
- src/runtime/metrics_test.go
- src/net/tcpsock_test.go
- src/reflect/all_test.go
- src/os/signal/signal_test.go


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
	runtime.InitAtomics()

	defer func() {
		runtime.DisableTrace()

		file_name := "dedego.log"
		os.Remove(file_name)
		file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		numRout := runtime.GetNumberOfRoutines()
		for i := 0; i <= numRout; i++ {
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
		file.Close()

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
}
```

Running this leads to the following trace (indented lines are in the same line 
as the previous line, only for better readability):

```

G,1,2;G,2,3;G,3,4;C,1,4,9,R,e,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgc.go:180;C,1,10,11,R,e,2,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgc.go:181;G,12,5;C,2,13,13,C,o,0,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/proc.go:256;G,14,6;G,15,7;C,4,16,20,R,e,1,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:74;C,4,21,22,R,e,2,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:74;C,4,23,30,R,e,3,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:74;C,4,31,32,R,e,4,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:74

C,1,7,8,S,e,2,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgcsweep.go:279
C,1,5,6,S,e,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/mgcscavenge.go:652

C,3,33,34,R,e,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:198;C,3,35,42,R,e,2,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:198;C,3,43,44,R,e,3,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:198;C,3,45,0,R,o,4,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/dedego_trace.go:198
C,4,17,18,S,e,1,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:63;C,4,19,24,S,e,2,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:63;C,4,25,26,S,e,3,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:63;C,4,27,27,C,o,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:65;A,28,824633794968;C,3,29,36,S,e,1,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedegoAtomic.go:34;A,37,824633794968;C,3,38,39,S,e,2,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedegoAtomic.go:34;A,40,824633794972;C,3,41,46,S,e,3,/home/erikkassubek/Uni/dedego/go-patch/src/runtime/internal/atomic/dedegoAtomic.go:34


```

The trace includes both the concurrent operations of the program it self, as well
as internal operations used by the go runtime. The elements from the
program are in file .../worst/main.go
