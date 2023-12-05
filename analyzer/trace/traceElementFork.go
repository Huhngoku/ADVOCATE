package trace

import (
	"analyzer/analysis"
	"errors"
	"strconv"
)

/*
* traceElementFork is a trace element for a go statement
* Fields:
*   routine (int): The routine id
*   tpost (int): The timestamp at the end of the event
*   id (int): The id of the new go statement
*  pos (string): The position of the trace element in the file
 */
type traceElementFork struct {
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

	elem := traceElementFork{
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
func (ro *traceElementFork) getRoutine() int {
	return ro.routine
}

/*
 * Get the tpre of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpre of the element
 */
func (ro *traceElementFork) getTpre() int {
	return ro.tPost
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (ro *traceElementFork) getTpost() int {
	return ro.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (ro *traceElementFork) getTsort() int {
	return ro.tPost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (ro *traceElementFork) toString() string {
	return "G" + "," + strconv.Itoa(ro.tPost) + "," + strconv.Itoa(ro.id) +
		"," + ro.pos
}

/*
 * Update and calculate the vector clock of the element
 */
func (ro *traceElementFork) updateVectorClock() {
	analysis.Fork(ro.routine, ro.id, currentVectorClocks)
}
