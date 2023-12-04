package rewriter

import "rewrite/bugs"

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
	panic("Not implemented")
}

/*
 * Rewrite the trace for a receive on a closed channel
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteTraceRecvOnClose(bug bugs.Bug) {
	panic("Not implemented")
}
