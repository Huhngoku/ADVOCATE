package rewriter

import (
	"analyzer/bugs"
	"analyzer/trace"
	"errors"
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
		println("Start rewriting trace for send on closed channel...")
		err = rewriteClosedChannel(bug)
	case bugs.PosRecvOnClosed:
		println("Start rewriting trace for receive on closed channel...")
		err = rewriteClosedChannel(bug)
	case bugs.RecvOnClosed:
		println("Actual receive on closed in trace. Therefor no rewrite is needed.")
	case bugs.CloseOnClosed:
		println("Only actual close on close can be detected. Therefor no rewrite is needed.")
	case bugs.DoneBeforeAdd:
		println("Start rewriting trace for negative waitgroup counter...")
		rewriteWaitGroup(bug)
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
		err = errors.New("Rewriting trace for cyclic deadlock is not implemented yet")
		// TODO: implement
	case bugs.RoutineLeakPartner:
		err = errors.New("Rewriting trace for routine leak with partner is not implemented yet")
		// TODO: implement
	case bugs.RoutineLeakNoPartner:
		err = errors.New("Rewriting trace for routine leak without partner is not implemented yet")
		// TODO: implement
	case bugs.RoutineLeakMutex:
		err = errors.New("Rewriting trace for routine leak with mutex is not implemented yet")
		// TODO: implement
	case bugs.RoutineLeakWaitGroup:
		err = errors.New("Rewriting trace for routine leak with waitgroup is not implemented yet")
		// TODO: implement
	case bugs.RoutineLeakCond:
		err = errors.New("Rewriting trace for routine leak with cond is not implemented yet")
		// TODO: implement
	default:
		err = errors.New("For the given bug type no trace rewriting is implemented")
	}
	if err != nil {
		println("Error rewriting trace")
	}
	return err
}

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
	(*bug.TraceElement1).SetTSortWithoutNotExecuted(t2)
	(*bug.TraceElement2[0]).SetTSortWithoutNotExecuted(t1)

	trace.AddElementToTrace(*bug.TraceElement1)
	trace.AddElementToTrace(*bug.TraceElement2[0])

	// add a stop marker
	trace.AddTraceElementReplayStop(t2 + 1)

	return nil
}

/*
 * Create a new trace from the given bug, given TraceElement2 has multiple elements
 * In this case, all elements in TraceElement2 should come directly after TraceElement1
 * The necessary before order should be kept
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteWaitGroup(bug bugs.Bug) {
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
	(*bug.TraceElement1).SetTSortWithoutNotExecuted(minTSort) // TODO: rewrite based on tpre
}
