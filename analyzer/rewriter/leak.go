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
 * into a possible bug, but to take an actual leak (we only detect actual leaks,
 * not possible leaks) and rewrite them in such a way, that the routine
 * gets unstuck, meaning is not leaking any more.
 * We detect leaks, that are stuck because of the following conditions:
 *  - channel operation without a possible  partner (may be in select)
 *  - channel operation with a possible partner, but no communication (may be in select)
 *  - mutex operation without a post event
 *  - waitgroup operation without a post event
 *  - cond operation without a post event
 */

// =============== Channel/Select ====================
// MARK: Channel/Select

/*
 * Rewrite a trace where a leaking unbuffered channel/select with possible partner was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteUnbufChanLeak(bug bugs.Bug) error {
	// check if one or both of the bug elements are select
	t1Sel := false
	t2Sel := false
	switch (*bug.TraceElement1[0]).(type) {
	case *trace.TraceElementSelect:
		t1Sel = true
	}
	switch (*bug.TraceElement2[0]).(type) {
	case *trace.TraceElementSelect:
		t2Sel = true
	}

	if !t1Sel && !t2Sel { // both are channel operations
		return rewriteUnbufChanLeakChanChan(bug)
	} else if !t1Sel && t2Sel { // first is channel operation, second is select
		return rewriteUnbufChanLeakChanSel(bug)
	} else if t1Sel && !t2Sel { // first is select, second is channel operation
		return rewriteUnbufChanLeakSelChan(bug)
	} // both are select
	return rewriteUnbufChanLeakSelSel(bug)
}

/*
 * Rewrite a trace for a leaking buffered channel
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func LeakBufChan(bug bugs.Bug) error {
	// check if one or both of the bug elements are select
	t1Sel := false
	t2Sel := false
	switch (*bug.TraceElement1[0]).(type) {
	case *trace.TraceElementSelect:
		t1Sel = true
	}
	switch (*bug.TraceElement2[0]).(type) {
	case *trace.TraceElementSelect:
		t2Sel = true
	}

	if !t1Sel && !t2Sel { // both are channel operations
		return rewriteUnbufChanLeakChanChan(bug)
	} else if !t1Sel && t2Sel { // first is channel operation, second is select
		return rewriteUnbufChanLeakChanSel(bug)
	} else if t1Sel && !t2Sel { // first is select, second is channel operation
		return rewriteUnbufChanLeakSelChan(bug)
	} // both are select
	return rewriteUnbufChanLeakSelSel(bug)
}

/*
 * Rewrite a trace where a leaking unbuffered channel/select with possible partner was found
 * if both elements are channel operations.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteUnbufChanLeakChanChan(bug bugs.Bug) error {
	stuck := (*bug.TraceElement1[0]).(*trace.TraceElementChannel)
	possiblePartner := (*bug.TraceElement2[0]).(*trace.TraceElementChannel)
	possiblePartnerPartner := possiblePartner.GetPartner()

	hb := clock.GetHappensBefore(possiblePartnerPartner.GetVC(), stuck.GetVC())
	if hb != clock.Concurrent {
		return errors.New("The actual partner of the potential partner is not HB " +
			"concurrent to the stuck element. Cannot rewrite trace.")
	}

	// T = T1 ++ [f] ++ T2 ++ [g] ++ T3 ++ [e]

	// remove the potential partner partner from the trace
	trace.RemoveElementFromTrace(possiblePartnerPartner.GetTID())

	// T = T1 ++ [f] ++ T2 ++ T3 ++ [e]

	if stuck.Operation() == trace.Recv { // Case 3
		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement1[0], possiblePartner.GetTSort()) // bug.TraceElement1[0] = stuck

		// T = T1 ++ [f] ++ T2' ++ T3' ++ [e]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]

		// add replay signals
		trace.AddTraceElementReplay(stuck.GetTSort() + 1)

	} else { // Case 4
		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement1[0], possiblePartnerPartner.GetTSort()) // bug.TraceElement1[0] = stuck

		// T = T1 ++ T2' ++ T3' ++ [e] ++ T4 ++ [f]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]
		// and T4 = [h in T4 | h >= e]

		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement2[0], stuck.GetTSort()) // bug.TraceElement2[0] = possiblePartner

		// T = T1 ++ T2' ++ T3' ++ [e] ++ T4' ++ [f]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]
		// and T4' = [h in T4 | h >= e and h < f]

	}

	// add replay signal
	trace.AddTraceElementReplay(possiblePartner.GetTSort() + 1)

	return nil
}

/*
 * Rewrite a trace where a leaking unbuffered channel/select with possible partner was found
 * if both elements are channel operations.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteUnbufChanLeakChanSel(bug bugs.Bug) error {
	stuck := (*bug.TraceElement1[0]).(*trace.TraceElementChannel)
	possiblePartner := (*bug.TraceElement2[0]).(*trace.TraceElementSelect)
	possiblePartnerPartner := possiblePartner.GetPartner()

	hb := clock.GetHappensBefore(possiblePartnerPartner.GetVC(), stuck.GetVC())
	if hb != clock.Concurrent {
		return errors.New("The actual partner of the potential partner is not HB " +
			"concurrent to the stuck element. Cannot rewrite trace.")
	}

	// T = T1 ++ [f] ++ T2 ++ [g] ++ T3 ++ [e]

	// remove the potential partner partner from the trace
	trace.RemoveElementFromTrace(possiblePartnerPartner.GetTID())

	// T = T1 ++ [f] ++ T2 ++ T3 ++ [e]

	if stuck.Operation() == trace.Recv { // Case 3
		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement1[0], possiblePartner.GetTSort()) // bug.TraceElement1[0] = stuck

		// T = T1 ++ [f] ++ T2' ++ T3' ++ [e]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]

		// add replay signal
		trace.AddTraceElementReplay(stuck.GetTSort() + 1)

	} else { // Case 4
		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement1[0], possiblePartnerPartner.GetTSort()) // bug.TraceElement1[0] = stuck

		// T = T1 ++ T2' ++ T3' ++ [e] ++ T4 ++ [f]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]
		// and T4 = [h in T4 | h >= e]

		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement2[0], stuck.GetTSort()) // bug.TraceElement2[0] = possiblePartner

		// T = T1 ++ T2' ++ T3' ++ [e] ++ T4' ++ [f]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]
		// and T4' = [h in T4 | h >= e and h < f]

	}

	// add replay signal
	trace.AddTraceElementReplay(possiblePartner.GetTSort() + 1)

	return nil
}

/*
 * Rewrite a trace where a leaking unbuffered channel/select with possible partner was found
 * if both elements are channel operations.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteUnbufChanLeakSelChan(bug bugs.Bug) error {
	stuck := (*bug.TraceElement1[0]).(*trace.TraceElementSelect)
	possiblePartner := (*bug.TraceElement2[0]).(*trace.TraceElementChannel)
	possiblePartnerPartner := possiblePartner.GetPartner()

	hb := clock.GetHappensBefore(possiblePartnerPartner.GetVC(), stuck.GetVC())
	if hb != clock.Concurrent {
		return errors.New("The actual partner of the potential partner is not HB " +
			"concurrent to the stuck element. Cannot rewrite trace.")
	}

	// T = T1 ++ [f] ++ T2 ++ [g] ++ T3 ++ [e]

	// remove the potential partner partner from the trace
	trace.RemoveElementFromTrace(possiblePartnerPartner.GetTID())

	// T = T1 ++ [f] ++ T2 ++ T3 ++ [e]

	if possiblePartner.Operation() == trace.Recv { // Case 4
		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement1[0], possiblePartnerPartner.GetTSort()) // bug.TraceElement1[0] = stuck

		// T = T1 ++ T2' ++ T3' ++ [e] ++ T4 ++ [f]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]
		// and T4 = [h in T4 | h >= e]

		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement2[0], stuck.GetTSort()) // bug.TraceElement2[0] = possiblePartner

		// T = T1 ++ T2' ++ T3' ++ [e] ++ T4' ++ [f]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]
		// and T4' = [h in T4 | h >= e and h < f]

	} else { // Case 3
		trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement1[0], possiblePartner.GetTSort()) // bug.TraceElement1[0] = stuck

		// T = T1 ++ [f] ++ T2' ++ T3' ++ [e]
		// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]

		// add replay signal
		trace.AddTraceElementReplay(stuck.GetTSort() + 1)

	}

	// add replay signals
	trace.AddTraceElementReplay(possiblePartner.GetTSort() + 1)

	return nil
}

/*
 * Rewrite a trace where a leaking unbuffered channel/select with possible partner was found
 * if both elements are channel operations.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteUnbufChanLeakSelSel(bug bugs.Bug) error {
	stuck := (*bug.TraceElement1[0]).(*trace.TraceElementSelect)
	possiblePartner := (*bug.TraceElement2[0]).(*trace.TraceElementSelect)
	possiblePartnerPartner := possiblePartner.GetPartner()

	hb := clock.GetHappensBefore(possiblePartnerPartner.GetVC(), stuck.GetVC())
	if hb != clock.Concurrent {
		return errors.New("The actual partner of the potential partner is not HB " +
			"concurrent to the stuck element. Cannot rewrite trace.")
	}

	// T = T1 ++ [f] ++ T2 ++ [g] ++ T3 ++ [e]

	// remove the potential partner partner from the trace
	trace.RemoveElementFromTrace(possiblePartnerPartner.GetTID())

	// find communication
	for _, c := range stuck.GetCases() {
		for _, d := range possiblePartner.GetCases() {
			if c.GetID() != d.GetID() {
				continue
			}

			if c.Operation() == d.Operation() {
				continue
			}

			// T = T1 ++ [f] ++ T2 ++ T3 ++ [e]

			if c.Operation() == trace.Recv { // Case 3
				trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement1[0], possiblePartner.GetTSort()) // bug.TraceElement1[0] = stuck

				// T = T1 ++ [f] ++ T2' ++ T3' ++ [e]
				// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]

				// add replay signal
				trace.AddTraceElementReplay(stuck.GetTSort() + 1)
				return nil
			}

			// Case 4
			trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement1[0], possiblePartner.GetTSort()) // bug.TraceElement1[0] = stuck

			// T = T1 ++ T2' ++ T3' ++ [e] ++ T4 ++ [f]
			// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]
			// and T4 = [h in T4 | h >= e]

			trace.ShiftConcurrentOrAfterToAfterStartingFromElement(bug.TraceElement2[0], stuck.GetTSort()) // bug.TraceElement2[0] = possiblePartner

			// T = T1 ++ T2' ++ T3' ++ [e] ++ T4' ++ [f]
			// where T2' = [h in T2 | h < e] and T3' = [h in T3 | h < e]
			// and T4' = [h in T4 | h >= e and h < f]

			// add replay signals
			trace.AddTraceElementReplay(possiblePartner.GetTSort() + 1)
		}
	}

	return errors.New("Could not establish communication between two selects. Cannot rewrite trace.")
}

// ================== Mutex ====================
// MARK: Mutex

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
	trace.AddTraceElementReplay(lockOp.GetTPre() + 1)

	return nil
}

// ================== WaitGroup ====================
// MARK: WaitGroup

/*
 * Rewrite a trace where a leaking waitgroup was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteWaitGroupLeak(bug bugs.Bug) error {
	println("Start rewriting trace for waitgroup leak...")

	wait := bug.TraceElement1[0]

	trace.ShiftConcurrentOrAfterToAfter(wait)

	trace.AddTraceElementReplay((*wait).GetTPre() + 1)

	return nil
}

// ================== Cond ====================
// MARK: Cond

/*
 * Rewrite a trace where a leaking cond was found.
 * Args:
 *   bug (Bug): The bug to create a trace for
 * Returns:
 *   error: An error if the trace could not be created
 */
func rewriteCondLeak(bug bugs.Bug) error {
	println("Start rewriting trace for cond leak...")

	couldRewrite := false

	wait := bug.TraceElement1[0]

	res := trace.GetConcurrentWaitgroups(wait)

	// possible signals to release the wait
	if len(res["signal"]) > 0 {
		couldRewrite = true

		(*wait).SetTSort((*wait).GetTPre())

		// move the signal after the wait
		trace.ShiftConcurrentOrAfterToAfter(wait)

		// TODO: Problem: locks create a happens before relation -> currently only works with -c
	}

	// possible broadcasts to release the wait
	for _, broad := range res["broadcast"] {
		couldRewrite = true
		trace.ShiftConcurrentToBefore(broad)
	}

	(*wait).SetTSort((*wait).GetTPre())

	trace.AddTraceElementReplay((*wait).GetTPre() + 1)

	if couldRewrite {
		return nil
	}

	return errors.New("Could not rewrite trace for cond leak")

}
