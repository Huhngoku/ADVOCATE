package trace

// enum for happens before
type happensBefore int

const (
	Before happensBefore = iota
	Concurrent
	After
	None
)

// Interface for trace elements
type traceElement interface {
	getTpre() int
	getTpost() int
	getRoutine() int
	getSimpleString() string
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

// TODO: Implement this
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

	// TODO: remove panic if implemented
	panic("Not implemented")
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
