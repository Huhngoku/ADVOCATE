package rewriter

import (
	"analyzer/bugs"
	"analyzer/trace"
	"errors"
)

/*
* Create a new trace for send/recv on closed channel
* Let c be the close, a the send/recv, X a stop marker and T1, T2, T3 partial traces
* The trace before the rewrite looks as follows:
* 	T1 ++ [a] ++ T2 ++ [c] ++ T3
* We know, that a, c and all elements in T2 are concurrent. Otherwise the bug
* would not have been detected. We are also not interested in T2 and T3. We
* can therefore rewrite the trace as follows:
* 	T1 ++ [c, a, X]
* Args:
*   bug (Bug): The bug to create a trace for
* Returns:
*   error: An error if the trace could not be created
 */
func rewriteClosedChannel(bug bugs.Bug) error {
	if bug.TraceElement1 == nil { // close
		return errors.New("TraceElement1 is nil") // send/recv
	}
	if bug.TraceElement2[0] == nil {
		return errors.New("TraceElement2 is nil")
	}

	t1 := (*bug.TraceElement1).GetTSort()    // close
	t2 := (*bug.TraceElement2[0]).GetTSort() // send/recv

	if t1 < t2 { // actual close before send/recv
		return errors.New("Close is before send/recv")
	}

	// shorten routine with send. After this, t1 and t2 are not in the trace anymore
	trace.ShortenTrace(t2, false)

	// switch the times of close and send/recv and add them at the end of the trace
	(*bug.TraceElement1).SetTPre(t2)
	(*bug.TraceElement2[0]).SetTPre(t1)

	trace.AddElementToTrace(*bug.TraceElement1)
	trace.AddElementToTrace(*bug.TraceElement2[0])

	// add a start and stop marker
	trace.AddTraceElementReplay(t1-1, true)
	trace.AddTraceElementReplay(t2+1, false)

	return nil
}
