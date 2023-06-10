package main

import (
	"fmt"
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
Package: dedego-analyzer (Deadlock Detector in Go)
Project: Dynamic Analysis to detect potential deadlocks in
*/

/*
Interface for a trace element.
@signature PrintElement(): function to print the element
*/
type TraceElement interface {
	// GetTimestamp() uint32
	PrintElement()
}

// ==================== Routine =====================

/*
Struct for the creation of a new routine
@field routine uint64: id of the new routine
*/
type TraceSpawn struct {
	routine uint64
}

/*
Function to print the signal trace element
@receiver TraceSignal
*/
func (ts TraceSpawn) PrintElement() {
	fmt.Printf("G(%d)", ts.routine)
}

// ==================== Channel =====================

/*
Struct for a pre in the trace.
@field id uint32: id of the Chan
@field exec bool: true if the operation was fully executed, false otherwise
@field send bool: true if the operation is a send, false if it is a receive
@field oId uint64: id of the operation in the channel
@field position string: string describing the position of the creation in code
*/
type TraceChan struct {
	id       uint64
	exec     bool
	send     bool
	oId      uint64
	position string
}

/*
Function to print the pre trace element
@receiver TraceChan
*/
func (tc TraceChan) PrintElement() {
	operation := "R"
	if tc.send {
		operation = "S"
	}

	exec_str := "o"
	if tc.exec {
		exec_str = "e"
	}

	fmt.Printf("C(%d,%s,%d,%s)", tc.id, operation, tc.oId, exec_str)
}

/*
Struct for a close in the trace.
@field postion string: string describing the position of the creation in code
@field timestamp uint32: timestamp of the creation of the trace object
@field chanId uint32: id of the Chan
*/
type TraceClose struct {
	position string
	id       uint64
}

/*
Function to print the close trace element
@receiver *TraceClose
*/
func (tc TraceClose) PrintElement() {
	fmt.Printf("C(%d, %s)", tc.id, "C")
}

// ==================== Select =====================
/*
Struct for a select case .
@field id uint32: id of the Chan
@field send bool: true if the operation is a send, false if it is a receive
*/
type TraceSelectCase struct {
	id   uint64
	send bool
}

/*
Function to get a string representation of the select case
@receiver TraceSelectCase
@return string: string representation of the select case
*/
func (tsc TraceSelectCase) toString() string {
	op := "r"
	if tsc.send {
		op = "s"
	}
	return fmt.Sprintf("%d:%s", tsc.id, op)
}

/*
Struct for a preSelect in the trace.
@field id uint32: id of the select
@field cases []TraceSelectCase: cases of the select
@field def bool: true if the select has a default case, false otherwise
@field exec bool: true if the operation was executed, false otherwise
@field chosen uint64: index of the chosen case
@field oId uint64: id of the operation in the channel
@field position string: string describing the position of the creation in code
*/
type TraceSelect struct {
	id       uint64
	cases    []TraceSelectCase
	def      bool
	exec     bool
	chosen   uint64
	oId      uint64
	position string
}

/*
Function to print the preSelect trace element
@receiver *TracePreSelect
*/
func (tps TraceSelect) PrintElement() {
	casesStr := ""
	for i, c := range tps.cases {
		if i != 0 {
			casesStr += "."
		}
		casesStr += c.toString()
	}
	if tps.def {
		if len(tps.cases) != 0 {
			casesStr += "."
		}
		casesStr += "d"
	}

	exec := "o"
	if tps.exec {
		exec = "e"
	}

	fmt.Printf("S(%d,%s,%s,%d,%d)", tps.id, casesStr, exec,
		tps.chosen, tps.oId)

}

// ==================== Mutex =====================

/*
Struct for a lock in the trace.
@field postion string: string describing the position of the creation in code
@field lockId uint32: id of the Mutex
@field rLock bool: true if it is a rwlock, false otherwise
@field try bool: true if it is a try-lock, false otherwise
@field read bool: true if it is a r-lock, false otherwise
@field suc bool: true if the operation was successful, false otherwise (only try)
@field exec bool: true if the operation was executed, false otherwise
*/
type TraceLock struct {
	id       uint64
	rLock    bool
	try      bool
	read     bool
	suc      bool
	exec     bool
	position string
}

/*
Function to print the lock trace element
@receiver *TraceLock
*/
func (tl TraceLock) PrintElement() {
	rw := "-"
	if tl.rLock {
		rw = "R"
	}

	operation := ""
	if tl.try {
		if tl.read {
			operation = "TR"
		} else {
			operation = "T"
		}
	} else {
		if tl.read {
			operation = "LR"
		} else {
			operation = "L"
		}
	}

	exec := "o"
	if tl.exec {
		exec = "e"
	}

	suc := "f"
	if tl.suc {
		suc = "s"
	}
	fmt.Printf("M(%d,%s,%s,%s,%s)", tl.id, rw, operation, exec, suc)
}

/*
Struct for a unlock in the trace.
@field lockId uint32: id of the Mutex
@field rLock bool: true if it is a rwlock, false otherwise
@field postion string: string describing the position of the creation in code
*/
type TraceUnlock struct {
	id       uint64
	rLock    bool
	position string
}

/*
Function to print the unlock trace element
@receiver *TraceUnlock
*/
func (tu TraceUnlock) PrintElement() {
	operation := "U"
	if tu.rLock {
		operation = "UR"
	}
	fmt.Printf("M(%d,%s)", tu.id, operation)
}

// ==================== WaitGroup =====================
/*
Struct for a  wait group operation in the trace.
@field postion string: string describing the position of the creation in code
@field id uint64: id of the WaitGroup
@field wait bool: true if it is a wait, false if add or done
@field delta int64: delta of the operation
@field val int64: value of the WaitGroup after the operation
@field exec bool: true if the operation was executed, false otherwise
*/
type TraceWaitGroup struct {
	position string
	id       uint64
	wait     bool
	delta    int64
	val      int64
	exec     bool
}

/*
Function to print the wait group trace element
@receiver *TraceWaitGroup
*/
func (tw TraceWaitGroup) PrintElement() {
	operation := ""
	if tw.wait {
		operation = "W"
	} else {
		if tw.delta > 0 {
			operation = "A"
		} else {
			operation = "D"
		}
	}
	exec := "o"
	if tw.exec {
		exec = "e"
	}
	fmt.Printf("W(%d,%s,%d,%d,%s)", tw.id, operation, tw.delta, tw.val,
		exec)
}
