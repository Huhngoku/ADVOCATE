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
	case bugs.RecvOnClosed:
		println("Start rewriting trace for receive on closed channel...")
		err = rewriteClosedChannel(bug)
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
 * Create a new trace from the given bug, given TraceElement2 has only one element
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteClosedChannel(bug bugs.Bug) error {
	if bug.TraceElement1 == nil {
		return errors.New("TraceElement1 is nil")
	}
	if bug.TraceElement2[0] == nil {
		return errors.New("TraceElement2 is nil")
	}

	routine1 := (*bug.TraceElement1).GetRoutine()    // close
	routine2 := (*bug.TraceElement2[0]).GetRoutine() // send

	// shorten routine with send
	err := trace.ShortenTrace(routine2, (*bug.TraceElement2[0]))
	if err != nil {
		return err
	}

	// shorten routine with close
	err = trace.ShortenTrace(routine1, (*bug.TraceElement1))
	if err != nil {
		return err
	}

	// switch the timer of send and close
	trace.SwitchTimer(bug.TraceElement1, bug.TraceElement2[0])

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
	(*bug.TraceElement1).SetTsortWithoutNotExecuted(minTSort) // TODO: rewrite based on tpre
}
