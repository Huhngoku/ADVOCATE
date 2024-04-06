package rewriter

import (
	"analyzer/bugs"
	"errors"
)

/*
 * Create a new trace for a negative wait group counter (done before add)
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteWaitGroup(bug bugs.Bug) error {
	for i := range bug.TraceElement1 {
		elem1 := bug.TraceElement1[i]
		elem2 := bug.TraceElement2[i]
	}

	return errors.New("Rewriting trace for negative waitgroup counter is not implemented yet")

}
