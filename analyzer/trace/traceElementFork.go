package trace

import (
	"analyzer/analysis"
	"errors"
	"strconv"
)

/*
* TraceElementFork is a trace element for a go statement
* Fields:
*   routine (int): The routine id
*   tpost (int): The timestamp at the end of the event
*   id (int): The id of the new go statement
*  pos (string): The position of the trace element in the file
*  tID (string): The id of the trace element, contains the position and the tpost
 */
type TraceElementFork struct {
	routine int
	tPost   int
	id      int
	pos     string
	tID     string
}

/*
 * Create a new go statement trace element
 * Args:
 *   routine (int): The routine id
 *   tPost (string): The timestamp at the end of the event
 *   id (string): The id of the new routine
 *   pos (string): The position of the trace element in the file
 */
func AddTraceElementFork(routine int, tPost string, id string, pos string) error {
	tPostInt, err := strconv.Atoi(tPost)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	tIDStr := pos + "@" + strconv.Itoa(tPostInt)

	elem := TraceElementFork{
		routine: routine,
		tPost:   tPostInt,
		id:      idInt,
		pos:     pos,
		tID:     tIDStr,
	}
	return AddElementToTrace(&elem)
}

/*
 * Get the id of the element
 * Returns:
 *   int: The id of the element
 */
func (fo *TraceElementFork) GetID() int {
	return fo.id
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (fo *TraceElementFork) GetRoutine() int {
	return fo.routine
}

/*
 * Get the tpre of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpre of the element
 */
func (fo *TraceElementFork) getTpre() int {
	return fo.tPost
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (fo *TraceElementFork) getTpost() int {
	return fo.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (fo *TraceElementFork) GetTSort() int {
	return fo.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The position of the element
 */
func (fo *TraceElementFork) GetPos() string {
	return fo.pos
}

/*
 * Get the tID of the element.
 * Returns:
 *   string: The tID of the element
 */
func (fo *TraceElementFork) GetTID() string {
	return fo.tID
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tsort (int): The timer of the element
 */
func (fo *TraceElementFork) SetTsort(tpost int) {
	fo.tPost = tpost
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tsort (int): The timer of the element
 */
func (fo *TraceElementFork) SetTSortWithoutNotExecuted(tsort int) {
	if fo.tPost != 0 {
		fo.tPost = tsort
	}
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (fo *TraceElementFork) ToString() string {
	return "G" + "," + strconv.Itoa(fo.tPost) + "," + strconv.Itoa(fo.id) +
		"," + fo.pos
}

/*
 * Update and calculate the vector clock of the element
 */
func (fo *TraceElementFork) updateVectorClock() {
	analysis.Fork(fo.routine, fo.id, currentVCHb, currentVCWmhb)
}
