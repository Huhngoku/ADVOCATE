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
		panic("Unknown bug type")
	}
}

/*
 * Rewrite the trace for a send on a closed channel
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteTraceSendOnClose(bug bugs.Bug) {
	skip := (*bug.TraceElement1).GetTSort() - (*bug.TraceElement2).GetTSort() + 1
	println("Skip: ", skip)
	trace.MoveTimeBack((*bug.TraceElement2).GetTSort(), skip, []int{(*bug.TraceElement1).GetRoutine()})
}

/*
 * Rewrite the trace for a receive on a closed channel
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteTraceRecvOnClose(bug bugs.Bug) {
	panic("Not implemented")
}
