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
	println("Start creating new trace...")
	switch bug.Type {
	case bugs.SendOnClosed:
		rewriteTraceSendOnClose(bug)
	case bugs.RecvOnClosed:
		rewriteTraceRecvOnClose(bug)
	default:
		println("Unknown bug type. Cannot rewrite trace.")
	}
}

/*
 * Rewrite the trace for a send on a closed channel
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteTraceSendOnClose(bug bugs.Bug) {
	println("Start rewriting trace for send on closed channel...")
	routineClose := (*bug.TraceElement1).GetRoutine()
	routineSend := (*bug.TraceElement2).GetRoutine()

	// shorten routine with send
	trace.ShortenTrace(routineSend, (*bug.TraceElement2))
	// shorten routine with close
	trace.ShortenTrace(routineClose, (*bug.TraceElement1))

	// switch the timer of send and close
	trace.SwitchTimer(bug.TraceElement1, bug.TraceElement2)
}

/*
 * Rewrite the trace for a receive on a closed channel
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteTraceRecvOnClose(bug bugs.Bug) {
	panic("Not implemented")
}

// skip := (*bug.TraceElement1).GetTSort() - (*bug.TraceElement2).GetTSort() + 1
// println("Skip: ", skip)
// trace.MoveTimeBack((*bug.TraceElement2).GetTSort(), skip, []int{(*bug.TraceElement1).GetRoutine()})
