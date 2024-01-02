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
	println("Start rewriting trace for send on closed channel...")
	routine1 := (*bug.TraceElement1).GetRoutine()    // close
	routine2 := (*bug.TraceElement2[0]).GetRoutine() // send

	// shorten routine with send
	trace.ShortenTrace(routine2, (*bug.TraceElement2[0]))
	// shorten routine with close
	trace.ShortenTrace(routine1, (*bug.TraceElement1))

	// switch the timer of send and close
	trace.SwitchTimer(bug.TraceElement1, bug.TraceElement2[0])
}
