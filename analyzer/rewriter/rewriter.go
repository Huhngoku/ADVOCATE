// Package rewriter provides functions for rewriting traces.
package rewriter

import (
	"analyzer/bugs"
	"analyzer/clock"
	"analyzer/trace"
	"analyzer/utils"
	"errors"
	"fmt"
	"sort"
)

var exitCode bool

const (
	exitCodeStuckFinish    = 10
	exitCodeStuckWaitElem  = 11
	exitCodeStuckNoElem    = 12
	exitCodeElemEmptyTrace = 13
	exitCodeLeakUnbuf      = 20
	exitCodeLeakBuf        = 21
	exitCodeLeakMutex      = 22
	exitCodeLeakCond       = 23
	exitCodeLeakWG         = 24
	exitSendClose          = 30
	exitRecvClose          = 31
	exitNegativeWG         = 32
	exitCodeCyclic         = 41
)

/*
 * Create a new trace from the given bug
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   bool: true if rewrite was needed, false otherwise (e.g. actual bug, warning)
 *   error: An error if the trace could not be created
 */
func RewriteTrace(bug bugs.Bug) (bool, error) {
	var err error
	rewriteNeeded := false
	switch bug.Type {
	case bugs.SendOnClosed:
		rewriteNeeded = true
		err = rewriteClosedChannel(bug, exitSendClose)
	case bugs.PosRecvOnClosed:
		rewriteNeeded = true
		err = rewriteClosedChannel(bug, exitRecvClose)
	case bugs.RecvOnClosed:
		err = errors.New("Actual receive on closed in trace. Therefore no rewrite is needed.")
	case bugs.CloseOnClosed:
		err = errors.New("Only actual close on close can be detected. Therefor no rewrite is needed.")
	case bugs.DoneBeforeAdd:
		rewriteNeeded = true
		err = rewriteWaitGroup(bug)
	case bugs.SelectWithoutPartner:
		err = errors.New("Rewriting trace for select without partner is not possible")
	case bugs.ConcurrentRecv:
		err = errors.New("Rewriting trace for concurrent receive is not possible")
	case bugs.MixedDeadlock:
		err = errors.New("Rewriting trace for mixed deadlock is not implemented yet")
	case bugs.CyclicDeadlock:
		rewriteNeeded = true
		err = rewriteCyclicDeadlock(bug)
	case bugs.LeakUnbufChanPartner:
		rewriteNeeded = true
		err = rewriteUnbufChanLeak(bug)
	case bugs.LeakUnbufChanNoPartner:
		err = errors.New("No possible partner for stuck channel found. Cannot rewrite trace.")
	case bugs.LeakBufChanPartner:
		rewriteNeeded = true
		err = rewriteBufChanLeak(bug)
	case bugs.LeakBufChanNoPartner:
		err = errors.New("No possible partner for stuck channel found. Cannot rewrite trace.")
	case bugs.LeakSelectPartnerUnbuf:
		rewriteNeeded = true
		err = rewriteUnbufChanLeak(bug)
	case bugs.LeakSelectPartnerBuf:
		rewriteNeeded = true
		err = rewriteBufChanLeak(bug)
	case bugs.LeakSelectNoPartner:
		err = errors.New("No possible partner for stuck select found. Cannot rewrite trace.")
	case bugs.LeakMutex:
		rewriteNeeded = true
		err = rewriteMutexLeak(bug)
	case bugs.LeakWaitGroup:
		rewriteNeeded = true
		err = rewriteWaitGroupLeak(bug)
	case bugs.LeakCond:
		rewriteNeeded = true
		err = rewriteCondLeak(bug)
	default:
		err = errors.New("For the given bug type no trace rewriting is implemented")
	}
	if rewriteNeeded && err != nil {
		println("Error rewriting trace")
	}
	return rewriteNeeded, err
}

/*
* Print the trace sorted by tPre
* Args:
*   types: types of the elements to print. If empty, all elements will be printed
*   clocks: if true, the clocks will be printed
* TODO: remove
 */
func PrintTrace(types []string, clocks bool) {
	elements := make([]struct {
		string
		int
		clock.VectorClock
	}, 0)
	for _, tra := range *trace.GetTraces() {
		for _, elem := range tra {
			elemStr := elem.ToString()
			if len(types) == 0 || utils.Contains(types, elemStr[0:1]) {
				elements = append(elements, struct {
					string
					int
					clock.VectorClock
				}{elemStr, elem.GetTPre(), elem.GetVC()})
			}
		}
	}

	// sort elements by timestamp
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].int < elements[j].int
	})

	for _, elem := range elements {
		if clocks {
			fmt.Println(elem.string, elem.VectorClock.ToString())
		} else {
			fmt.Println(elem.string)
		}
	}
}
