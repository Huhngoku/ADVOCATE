package rewriter

import (
	"analyzer/bugs"
	"analyzer/trace"
	"fmt"
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

	PrintTrace([]string{"W", "C"})
	println("\n\n")
	for _, trace := range *trace.GetTraces() {
		for _, elem := range trace {
			if elem.GetVC().GetSize() == 0 {
				continue
			}
			fmt.Println(elem.ToString(), " ", elem.GetVC())
		}
	}
	println("\n\n")
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
		trace.AddTraceElementReplay(minTime-1, true)
		trace.AddTraceElementReplay(maxTime+1, false)
	}

	PrintTrace([]string{"W", "C"})
	return nil
}
