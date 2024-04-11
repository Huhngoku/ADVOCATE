package trace

import (
	"analyzer/clock"
	"strconv"
)

/*
 * Struct to save an atomic event in the trace
 * Fields:
 *   tpost (int): The timestamp of the event
 *   start (int): True if before the important event, false if after
 */
type TraceElementReplay struct {
	tPost int
	start bool
}

/*
 * Create a new atomic trace element
 * Args:
 *   tpost (string): The timestamp of the event
 *   start (int): True if before the important event, false if after
 */
func AddTraceElementReplay(tPost int, start bool) error {
	elem := TraceElementReplay{
		tPost: tPost,
		start: start,
	}

	return AddElementToTrace(&elem)
}

/*
 * Get the id of the element
 * Returns:
 *   int: The id of the element
 */
func (at *TraceElementReplay) GetID() int {
	return 0
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (at *TraceElementReplay) GetRoutine() int {
	return 1
}

/*
 * Get the tpost of the element.
 *   int: The tpost of the element
 */
func (at *TraceElementReplay) GetTPre() int {
	return at.tPost
}

/*
 * Set the tpre of the element.
 * Args:
 *   tPre (int): The tpre of the element
 */
func (mu *TraceElementReplay) SetTPre(tPre int) {
	mu.tPost = tPre
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (at *TraceElementReplay) getTpost() int {
	return at.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (at *TraceElementReplay) GetTSort() int {
	return at.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The file of the element
 */
func (at *TraceElementReplay) GetPos() string {
	return ""
}

/*
 * Get the tID of the element.
 * Returns:
 *   string: The tID of the element
 */
func (at *TraceElementReplay) GetTID() string {
	return ""
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tSort (int): The timer of the element
 */
func (at *TraceElementReplay) SetTSort(tSort int) {
	at.SetTPre(tSort)
	at.tPost = tSort
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tSort (int): The timer of the element
 */
func (at *TraceElementReplay) SetTSortWithoutNotExecuted(tSort int) {
	at.SetTPre(tSort)
	at.tPost = tSort
}

/*
 * Get the simple string representation of the element.
 * Returns:
 *   string: The simple string representation of the element
 */
func (at *TraceElementReplay) ToString() string {
	res := "X," + strconv.Itoa(at.tPost) + ","
	if at.start {
		res += "s"
	} else {
		res += "e"
	}
	return res
}

/*
 * Update and calculate the vector clock of the element
 */
func (at *TraceElementReplay) updateVectorClock() {
	// nothing to do
}

/*
 * Dummy function to implement the interface
 * Returns:
 *   VectorClock: The vector clock of the element
 */
func (at *TraceElementReplay) GetVC() clock.VectorClock {
	return clock.VectorClock{}
}
