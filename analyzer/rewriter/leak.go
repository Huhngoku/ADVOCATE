package rewriter

import (
	"analyzer/bugs"
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
 *  - implement rewriteRoutineLeak
 *  - implement rewriteMutexLeak
 *  - implement rewriteWaitGroupLeak
 *  - implement rewriteCondLeak
 */

/*
 * Rewrite a trace where a leaking channel with possible partner was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteRoutineLeak(bug bugs.Bug) error {
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

/*
 * Rewrite a trace where a leaking mutex was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteMutexLeak(bug bugs.Bug) error {
	println("Start rewriting trace for mutex leak...")
	return errors.New("Rewriting trace for routine leak with mutex is not implemented yet")
}

/*
 * Rewrite a trace where a leaking waitgroup was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteWaitGroupLeak(bug bugs.Bug) error {
	println("Start rewriting trace for waitgroup leak...")
	return errors.New("Rewriting trace for routine leak with waitgroup is not implemented yet")
}

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
