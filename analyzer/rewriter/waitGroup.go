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
	println("Start rewriting trace for negative waitgroup counter...")

	minTime := -1
	maxTime := -1

	for i := range bug.TraceElement1 {
		elem2 := bug.TraceElement2[i] // done

		trace.ShiftConcurrentOrAfterToAfter(elem2)

		if minTime == -1 || (*elem2).GetTPre() < minTime {
			minTime = (*elem2).GetTPre()
		}
		if maxTime == -1 || (*elem2).GetTPre() > maxTime {
			maxTime = (*elem2).GetTPre()
		}

	}

	// add start and end
	if !(minTime == -1 && maxTime == -1) {
		trace.AddTraceElementReplay(maxTime + 1)
	}

	return nil
}
