package trace

import (
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
 */
type TraceElementFork struct {
	routine int
	tPost   int
	id      int
	pos     string
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

	elem := TraceElementFork{
		routine: routine,
		tPost:   tPostInt,
		id:      idInt,
		pos:     pos}
	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (ro *TraceElementFork) GetRoutine() int {
	return ro.routine
}

/*
 * Get the tpre of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpre of the element
 */
func (ro *TraceElementFork) getTpre() int {
	return ro.tPost
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (ro *TraceElementFork) getTpost() int {
	return ro.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (ro *TraceElementFork) GetTSort() int {
	return ro.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The position of the element
 */
func (at *TraceElementFork) GetPos() string {
	return at.pos
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tsort (int): The timer of the element
 */
func (te *TraceElementFork) SetTsort(tpost int) {
	te.tPost = tpost
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tsort (int): The timer of the element
 */
func (te *TraceElementFork) SetTsortWithoutNotExecuted(tsort int) {
	if te.tPost != 0 {
		te.tPost = tsort
	}
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (ro *TraceElementFork) ToString() string {
	return "G" + "," + strconv.Itoa(ro.tPost) + "," + strconv.Itoa(ro.id) +
		"," + ro.pos
}
