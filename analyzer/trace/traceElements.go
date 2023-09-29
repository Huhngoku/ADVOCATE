package trace

import vc "analyzer/vectorClock"

// Interface for trace elements
type traceElement interface {
	getTpre() int
	getTpost() int
	getTsort() int
	// getVpre() *vc.VectorClock
	getVpost() *vc.VectorClock
	getRoutine() int
	toString() string
	updateVectorClock()
}

/*
* Return a given happens-before relationship (befor, after, concurrent),
* given two trace elements. This relationship must hold even with a valid
* reordering of the trace.
* Args:
*   first (traceElement): The first trace element
*   second (traceElement): The second trace element
* Returns:
*   happensBefore: The relationship between the two trace elements
 */
func GetHappensBefore(first traceElement, second traceElement) vc.HappensBefore {
	return vc.GetHappensBefore(*first.getVpost(), *second.getVpost())
}
