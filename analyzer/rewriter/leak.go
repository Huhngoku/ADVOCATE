package rewriter

import (
	"analyzer/bugs"
	"analyzer/clock"
	"analyzer/trace"
	"errors"
)

/*
 * Rewrite a trace where a leaking routine was found.
 * Different to most other rewrites, we don not try to get the program to run
 * into a potential bug, but to take an actual leak (we only detect actual leaks,
 * not potential leaks) and rewrite them in such a way, that the routine
 * gets unstuck, meaning is not leaking any more.
 * We detect leaks, that are stuck because of the following conditions:
 *  - channel operation without a potential  partner (may be in select)
 *  - channel operation with a potential partner, but no communication (may be in select)
 *  - mutex operation without a post event
 *  - waitgroup operation without a post event
 *  - cond operation without a post event
 * TODO:
 *  - implement rewriteChannelLeak
 *  - implement rewriteWaitGroupLeak -> not possible????
 *  - implement rewriteCondLeak
 *  - look at stuck select
 */

// =============== Channel ====================

/*
 * Rewrite a trace where a leaking channel with possible partner was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteChannelLeak(bug bugs.Bug) error {
	return errors.New("Rewriting trace for routine leak with partner is not implemented yet")

	// println("Start rewriting trace for channel leak...")
	// println((*bug.TraceElement1[0]).ToString()) // stuck
	// println((*bug.TraceElement2[0]).ToString()) // potential partner

	// // get the original partner of the potential partner, and set its post and oId to 0
	// originalPartner := (*bug.TraceElement2[0]).(*trace.TraceElementChannel).GetPartner()
	// println(originalPartner.ToString())
	// originalPartner.SetTPost(0)
	// originalPartner.SetOID(0)

	// // set the oId of the stuck operation to the oId of the potential partner
	// (*bug.TraceElement1[0]).(*trace.TraceElementChannel).SetOID((*bug.TraceElement2[0]).(*trace.TraceElementChannel).GetOID())

	// // TODO: shift correctly
	// distance := (*bug.TraceElement2[0]).GetTPre() - (*bug.TraceElement1[0]).GetTPre()
	// trace.ShiftRoutine((*bug.TraceElement1[0]).GetRoutine(), (*bug.TraceElement1[0]).GetTPre(), distance)

	// println((*bug.TraceElement1[0]).ToString())
	// println((*bug.TraceElement2[0]).ToString())
	// println(originalPartner.ToString())

}

// ================== Mutex ====================

/*
 * Rewrite a trace where a leaking mutex was found.
 * The trace can only be rewritten, if the stuck lock operation is concurrent
 * with the last lock operation on this mutex. If it is not concurrent, the
 * rewrite fails. If a rewrite is possible, we try to run the stock lock operation
 * before the last lock operation, so that the mutex is not blocked anymore.
 * We therefore rewrite the trace from
 *   T_1 + [l'] + T_2 + [l] + T_3
 * to
 *   T_1' + T_2' + [X_s, l, X_e]
 * where l is the stuck lock, l' is the last lock, T_1, T_2, T_3 are the traces
 * before, between and after the locks, T_1' and T_2' are the elements from T_1 and T_2, that
 * are before (HB) l, X_s is the start and X_e is the stop signal, that releases the program from the
 * guided replay.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteMutexLeak(bug bugs.Bug) error {
	println("Start rewriting trace for mutex leak...")

	// get l and l'
	lockOp := (*bug.TraceElement1[0]).(*trace.TraceElementMutex)
	lastLockOp := (*bug.TraceElement2[0]).(*trace.TraceElementMutex)

	hb := clock.GetHappensBefore(lockOp.GetVC(), lastLockOp.GetVC())
	if hb != clock.Concurrent {
		return errors.New("The stuck mutex lock is not concurrent with the prior lock. Cannot rewrite trace.")
	}

	// remove T_3 -> T_1 + [l'] + T_2 + [l]
	trace.ShortenTrace(lockOp.GetTSort(), true)

	// remove all elements, that are concurrent with l. This includes l'
	// -> T_1' + T_2' + [l]
	trace.RemoveConcurrent(bug.TraceElement1[0])

	// set tpost of l to non zero
	lockOp.SetTSort(lockOp.GetTPre())

	// add the start and stop signal after l -> T_1' + T_2' + [X_s, l, X_e]
	trace.AddTraceElementReplay(lockOp.GetTPre()-1, true)
	trace.AddTraceElementReplay(lockOp.GetTPre()+1, false)

	return nil
}

// ================== WaitGroup ====================

/*
 * Rewrite a trace where a leaking waitgroup was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteWaitGroupLeak(bug bugs.Bug) error {
	// println("Start rewriting trace for waitgroup leak...")
	return errors.New("Rewrite for leaking waitgroup not possible")
}

// ================== Cond ====================

/*
 * Rewrite a trace where a leaking cond was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteCondLeak(bug bugs.Bug) error {
	println("Start rewriting trace for cond leak...")
	return errors.New("Rewriting trace for routine leak with cond is not implemented yet")
}
