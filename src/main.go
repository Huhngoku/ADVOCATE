package main

import "flag"

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
Package: dedego-instrumenter
Project: Dynamic Analysis to detect potential deadlocks in concurrent Go programs
*/

// TODO: remove this
var traces [][]TraceElement
var chanSize []int

type PreObj struct {
	id           uint32
	chanCreation string
	receive      bool
}

// END TODO

/*
main.go
main file to run the analyser
*/
func main() {
	traceFileNamePtr := flag.String("t", "", "trace file name")

	flag.Parse()

	traceFileName := "dedego.log"
	if *traceFileNamePtr != "" {
		traceFileName = *traceFileNamePtr
	}

	traces = make([][]TraceElement, 0)

	createTrace(traceFileName, &traces)
	printTrace(&traces)
}
