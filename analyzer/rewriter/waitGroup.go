package rewriter

import (
	"analyzer/bugs"
	"analyzer/trace"
	"errors"
	"sort"
)

/*
 * Create a new trace for a negative wait group counter (done before add)
 * Args:
 *   bug (Bug): The bug to create a trace for
 */
func rewriteWaitGroup(bug bugs.Bug) error {
	// get all concurrent add and done operations
	ops := make(map[int][]*trace.TraceElementWait, 0)

	for _, elem := range bug.TraceElement2 {
		if elem == nil {
			continue
		}

		routine := (*elem).GetRoutine()

		ops[routine] = append(ops[routine], (*elem).(*trace.TraceElementWait))

	}

	// sort adds and dones by time
	for _, add := range ops {
		sort.Slice(add, func(i, j int) bool {
			return add[i].GetTSort() > add[j].GetTSort()
		})
	}

	// TODO: continue implementation

	return errors.New("Rewriting trace for negative waitgroup counter is not implemented yet")
}
