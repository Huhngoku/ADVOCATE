// Package rewriter provides functions for rewriting traces.
package rewriter

import (
	"analyzer/bugs"
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
		err = errors.New("Rewriting trace for select without partner is not implemented yet")
		// TODO: implement
	case bugs.ConcurrentRecv:
		err = errors.New("Rewriting trace for concurrent receive is not implemented yet")
		// TODO: implement
	case bugs.MixedDeadlock:
		err = errors.New("Rewriting trace for mixed deadlock is not implemented yet")
		// TODO: implement
	case bugs.CyclicDeadlock:
		err = rewriteCyclicDeadlock(bug)
	case bugs.RoutineLeakPartner:
		err = rewriteRoutineLeak(bug)
	case bugs.RoutineLeakNoPartner:
		err = errors.New("No possible partner for stuck channel found. Cannot rewrite trace.")
	case bugs.RoutineLeakMutex:
		err = rewriteMutexLeak(bug)
	case bugs.RoutineLeakWaitGroup:
		err = rewriteWaitGroupLeak(bug)
	case bugs.RoutineLeakCond:
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
* TODO: remove
 */
func PrintTrace(types []string) {
	elements := make([]struct {
		string
		int
	}, 0)
	for _, tra := range *trace.GetTraces() {
		for _, elem := range tra {
			elemStr := elem.ToString()
			if len(types) == 0 || utils.Contains(types, elemStr[0:1]) {
				elements = append(elements, struct {
					string
					int
				}{elemStr, elem.GetTPre()})
			}
		}
	}

	// sort elements by timestamp
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].int < elements[j].int
	})

	for _, elem := range elements {
		fmt.Println(elem.string)
	}
}
