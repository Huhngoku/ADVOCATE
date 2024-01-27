package rewriter

import (
	"analyzer/bugs"
	"analyzer/trace"
)

/*
 * Create a new trace from the given bug
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func RewriteTrace(bug bugs.Bug) {
	switch bug.Type {
	case bugs.SendOnClosed:
		println("Start rewriting trace for send on closed channel...")
		rewriteTraceSingle(bug)
	case bugs.RecvOnClosed:
		println("Start rewriting trace for receive on closed channel...")
		rewriteTraceSingle(bug)
	case bugs.DoneBeforeAdd:
		println("Start rewriting trace for negative waitgroup counter...")
		rewriteTraceMultiple(bug)
	default:
		println("For the given bug type no trace rewriting is implemented")
	}
}

/*
 * Create a new trace from the given bug, given TraceElement2 has only one element
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteTraceSingle(bug bugs.Bug) {
	routine1 := (*bug.TraceElement1).GetRoutine()    // close
	routine2 := (*bug.TraceElement2[0]).GetRoutine() // send

	// shorten routine with send
	trace.ShortenTrace(routine2, (*bug.TraceElement2[0]))
	// shorten routine with close
	trace.ShortenTrace(routine1, (*bug.TraceElement1))

	// switch the timer of send and close
	trace.SwitchTimer(bug.TraceElement1, bug.TraceElement2[0])
}

/*
 * Create a new trace from the given bug, given TraceElement2 has multiple elements
 * In this case, all elements in TraceElement2 should come directly after TraceElement1
 * The necessary before order should be kept
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteTraceMultiple(bug bugs.Bug) {
	// get the smallest tSort in TraceElement1 or TraceElement2
	minTSort := (*bug.TraceElement1).GetTSort()
	maxTSort := (*bug.TraceElement1).GetTSort()
	for _, elem := range bug.TraceElement2 {
		if elem == nil {
			continue
		}
		if (*elem).GetTSort() < minTSort {
			minTSort = (*elem).GetTSort()
		}
		if (*elem).GetTSort() > maxTSort {
			maxTSort = (*elem).GetTSort()
		}
	}
	trace.ShiftTrace(minTSort, maxTSort-minTSort+1)
	(*bug.TraceElement1).SetTsortWithoutNotExecuted(minTSort)
}
