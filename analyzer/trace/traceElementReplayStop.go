package trace

import (
	"strconv"
)

/*
 * Struct to save an atomic event in the trace
 * Fields:
 *   tpost (int): The timestamp of the event
 */
type TraceElementReplayStop struct {
	tPost int
}

/*
 * Create a new atomic trace element
 * Args:
 *   tpost (string): The timestamp of the event
 */
func AddTraceElementReplayStop(tPost int) error {
	elem := TraceElementAtomic{
		tPost: tPost,
	}

	return AddElementToTrace(&elem)
}

/*
 * Get the id of the element
 * Returns:
 *   int: The id of the element
 */
func (at *TraceElementReplayStop) GetID() int {
	return 0
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (at *TraceElementReplayStop) GetRoutine() int {
	return 1
}

/*
 * Get the tpost of the element.
 *   int: The tpost of the element
 */
func (at *TraceElementReplayStop) getTpre() int {
	return at.tPost
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (at *TraceElementReplayStop) getTpost() int {
	return at.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (at *TraceElementReplayStop) GetTSort() int {
	return at.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The file of the element
 */
func (at *TraceElementReplayStop) GetPos() string {
	return ""
}

/*
 * Get the tID of the element.
 * Returns:
 *   string: The tID of the element
 */
func (at *TraceElementReplayStop) GetTID() string {
	return ""
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tSort (int): The timer of the element
 */
func (at *TraceElementReplayStop) SetTsort(tSort int) {
	at.tPost = tSort
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tSort (int): The timer of the element
 */
func (at *TraceElementReplayStop) SetTSortWithoutNotExecuted(tSort int) {
	at.tPost = tSort
}

/*
 * Get the simple string representation of the element.
 * Returns:
 *   string: The simple string representation of the element
 */
func (at *TraceElementReplayStop) ToString() string {
	res := "X," + strconv.Itoa(at.tPost)
	return res
}

/*
 * Update and calculate the vector clock of the element
 */
func (at *TraceElementReplayStop) updateVectorClock() {
	// nothing to do
}
