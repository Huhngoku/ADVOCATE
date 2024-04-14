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

/*
 * Create a new trace from the given bug
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func RewriteTrace(bug bugs.Bug) error {
	var err error
	switch bug.Type {
	case bugs.SendOnClosed:
		err = rewriteClosedChannel(bug)
	case bugs.PosRecvOnClosed:
		err = rewriteClosedChannel(bug)
	case bugs.RecvOnClosed:
		err = errors.New("Actual receive on closed in trace. Therefore no rewrite is needed.")
	case bugs.CloseOnClosed:
		err = errors.New("Only actual close on close can be detected. Therefor no rewrite is needed.")
	case bugs.DoneBeforeAdd:
		err = rewriteWaitGroup(bug)
	case bugs.SelectWithoutPartner:
		err = errors.New("Rewriting trace for select without partner is not possible")
		// TODO: implement
	case bugs.ConcurrentRecv:
		err = errors.New("Rewriting trace for concurrent receive is not possible")
		// TODO: implement
	case bugs.MixedDeadlock:
		err = errors.New("Rewriting trace for mixed deadlock is not implemented yet")
		// TODO: implement
	case bugs.CyclicDeadlock:
		err = rewriteCyclicDeadlock(bug)
	case bugs.LeakUnbufChanPartner:
		err = rewriteUnbufChanLeak(bug)
	case bugs.LeakUnbufChanNoPartner:
		err = errors.New("No possible partner for stuck channel found. Cannot rewrite trace.")
	case bugs.LeakBufChan:
		err = LeakBufChan(bug)
	case bugs.LeakMutex:
		err = rewriteMutexLeak(bug)
	case bugs.LeakWaitGroup:
		err = rewriteWaitGroupLeak(bug)
	case bugs.LeakCond:
		err = rewriteCondLeak(bug)
	default:
		err = errors.New("For the given bug type no trace rewriting is implemented")
	}
	if err != nil {
		println("Error rewriting trace")
	}
	return err
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
