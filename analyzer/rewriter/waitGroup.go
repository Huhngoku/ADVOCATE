package rewriter

import (
	"analyzer/bugs"
	"analyzer/trace"
)

/*
 * Create a new trace for a negative wait group counter (done before add)
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteWaitGroup(bug bugs.Bug) error {
	// TODO: does not work yet -> gets stuck -> do not just shift  routine
	// TODO: check if the pairs must be sorted

	println("Start rewriting trace for negative waitgroup counter...")
	// for each pair of element, move the add after the done
	minTime := -1
	maxTime := -1
	for i := range bug.TraceElement1 {
		elem1 := bug.TraceElement1[i] // add
		elem2 := bug.TraceElement2[i] // done

		if minTime == -1 || (*bug.TraceElement2[i]).GetTPre() < minTime {
			minTime = (*bug.TraceElement2[i]).GetTPre()
		}

		shift := (*elem2).GetTPre() - (*elem1).GetTPre() + 1
		if maxTime == -1 || (*elem2).GetTPre()+shift > maxTime {
			maxTime = (*elem2).GetTPre() + shift
		}

		trace.ShiftRoutine((*elem1).GetRoutine(), (*elem1).GetTPre(), shift)
	}

	trace.AddTraceElementReplay(minTime, true)
	trace.AddTraceElementReplay(maxTime, false)

	return nil
}
