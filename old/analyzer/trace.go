package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

/*
Copyright (c) 2023, Erik Kassubek
All rights reserved.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: main
Project: Dynamic Analysis to detect potential deadlocks in concurrent Go programs
*/

/*
tracerRead.go
Read the trace file and build the trace
*/

/*
Read the trace file and build the trace
@param fileName string: name of the trace file
@return traces *[][]TraceElement: pointer tp the trace
*/
func createTrace(fileName string, traces *([][]TraceElement)) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	// traverse all lines
	for _, line := range lines {
		// split line into elements
		elements := strings.Split(line, ";")
		trace := make([]TraceElement, 0) // declare trace slice
		for _, element := range elements {
			// split into fields
			fields := strings.Split(element, ",")
			var elem TraceElement

			if len(fields) <= 1 {
				continue
			}

			switch fields[0] {
			case "G": // spawn new routine
				routine, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					panic(err)
				}

				elem = TraceSpawn{routine: routine}
			case "M": // mutex
				id, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					panic(err)
				}

				rw := false
				if fields[2] == "R" {
					rw = true
				}

				try := false
				if fields[3] == "T" || fields[3] == "TR" {
					try = true
				}

				read := false
				if fields[3] == "LR" || fields[3] == "TR" {
					read = true
				}

				exec := false
				if fields[4] == "e" {
					exec = true
				}

				suc := true
				if fields[5] == "f" {
					suc = false
				}

				if fields[3] == "L" || fields[3] == "LR" || fields[3] == "TR" ||
					fields[3] == "T" {
					elem = TraceLock{id: id, rLock: rw, try: try,
						read: read, position: fields[6], exec: exec, suc: suc}
				} else if fields[3] == "U" || fields[3] == "UR" {
					elem = TraceUnlock{id: id, rLock: rw,
						position: fields[6]}
				} else {
					panic("Unknown mutex operation")
				}
			case "W": // waitGroup
				id, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					panic(err)
				}

				wait := false
				if fields[2] == "W" {
					wait = true
				}

				exec := false
				if fields[3] == "e" {
					exec = true
				}

				delta, err := strconv.ParseInt(fields[4], 10, 64)
				if err != nil {
					panic(err)
				}

				val, err := strconv.ParseInt(fields[5], 10, 64)
				if err != nil {
					panic(err)
				}

				elem = TraceWaitGroup{id: id, wait: wait, exec: exec,
					delta: delta, val: val, position: fields[6]}

			case "C": // channel
				id, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					panic(err)
				}

				exec := false
				if fields[3] == "e" {
					exec = true
				}

				oId, err := strconv.ParseUint(fields[4], 10, 64)
				if err != nil {
					panic(err)
				}

				if fields[2] == "S" {
					elem = TraceChan{id: id, exec: exec, send: true, oId: oId,
						position: fields[5]}
				} else if fields[2] == "R" {
					elem = TraceChan{id: id, exec: exec, send: false, oId: oId,
						position: fields[5]}
				} else if fields[2] == "C" {
					elem = TraceClose{id: id, position: fields[5]}
				} else {
					panic("Unknown channel operation")
				}
			case "S": // select
				id, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					panic(err)
				}

				cases := make([]TraceSelectCase, 0)
				def := false
				c := strings.Split(fields[2], ".")
				for _, elem := range c {
					if elem == "d" {
						def = true
						continue
					}
					chanId, err := strconv.ParseUint(elem[:len(elem)-1], 10, 64)
					if err != nil {
						panic(err)
					}

					if elem[len(elem)-1] == 's' {
						cases = append(cases, TraceSelectCase{id: chanId,
							send: true})
					} else if elem[len(elem)-1] == 'r' {
						cases = append(cases, TraceSelectCase{id: chanId,
							send: false})
					} else {
						panic("Unknown select case")
					}
				}

				exec := false
				if fields[3] == "e" {
					exec = true
				}

				chosen, err := strconv.ParseUint(fields[4], 10, 64)
				if err != nil {
					panic(err)
				}

				oId, err := strconv.ParseUint(fields[5], 10, 64)
				if err != nil {
					panic(err)
				}

				elem = TraceSelect{id: id, cases: cases, def: def, exec: exec,
					chosen: chosen, oId: oId, position: fields[6]}
			default:
				msg := fmt.Sprintf("Unknown trace element %s.", element)
				panic(msg)
			}

			if elem == nil {
				error_msg := fmt.Sprintf("nil element in %s", element)
				panic(error_msg)
			}

			trace = append(trace, elem) // add elem to trace
		}
		*traces = append(*traces, trace) // add trace to traces

	}
}

/*
Print the trace
@param traces *[][]TraceElement: trace to be printed
*/
func printTrace(traces *([][]TraceElement)) {
	for _, trace := range *traces {
		for i, elem := range trace {
			if i != 0 {
				fmt.Print(";")
			}
			elem.PrintElement()
		}
		fmt.Println("")
	}
}
