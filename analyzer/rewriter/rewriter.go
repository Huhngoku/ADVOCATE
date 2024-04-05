// Package rewriter provides functions for rewriting traces.
package rewriter

import (
	"analyzer/bugs"
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
	case bugs.CyclicDeadlockTwo:
		err = rewriteCyclicDeadlock(bug)
	case bugs.CyclicDeadlockMulti:
		err = rewriteCyclicDeadlock(bug)
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
