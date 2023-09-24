package trace

// Interface for trace elements
type traceElement interface {
	getTpre() int
	getTpost() int
	getVpre() *vectorClock
	getVpost() *vectorClock
	getRoutine() int
	toString() string
}

/*
 * Get the relationship between two trace elements in the recorded run
 * Args:
 *   first (traceElement): The first trace element
 *   second (traceElement): The second trace element
 * Returns:
 *   happensBefore: The relationship between the two trace elements
 */
func GetRelationshipInRecordedRun(first traceElement, second traceElement) happensBefore {
	if first.getTpost() < second.getTpre() {
		return Before
	} else if first.getTpre() > second.getTpost() {
		return After
	} else {
		return Concurrent
	}
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
func GetHappensBefore(first traceElement, second traceElement) happensBefore {
	// if elements are in the same routine
	if first.getRoutine() == second.getRoutine() {
		return getHappensBeforeSameRoutine(first, second)
	}

	return getHappensBefore(first.getVpre(), first.getVpost(), second.getVpre(),
		second.getVpost())
}

/*
* Return a given happens-before relationship (befor, after, concurrent),
* given two trace elements. The function assumes, that both elements are in the
* same routine.
* Args:
*   first (traceElement): The first trace element
*   second (traceElement): The second trace element
* Returns:
*   happensBefore: The relationship between the two trace elements
 */
func getHappensBeforeSameRoutine(first traceElement, second traceElement) happensBefore {
	return GetRelationshipInRecordedRun(first, second)
}
