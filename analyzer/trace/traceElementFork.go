package trace

import (
	vc "analyzer/vectorClock"
	"errors"
	"strconv"
)

/*
 * traceElementFork is a trace element for a go statement
 * Fields:
 *   routine (int): The routine id
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the new go statement
 */
type traceElementFork struct {
	routine int
	tpost   int
	id      int
}

/*
 * Create a new go statement trace element
 * Args:
 *   routine (int): The routine id
 *   numberOfRoutines (int): The number of routines in the trace
 *   tpost (string): The timestamp at the end of the event
 *   id (string): The id of the new routine
 */
func AddTraceElementFork(routine int, numberOfRoutines int, tpost string,
	id string) error {
	tpost_int, err := strconv.Atoi(tpost)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	elem := traceElementFork{
		routine: routine,
		tpost:   tpost_int,
		id:      id_int}
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
	return ro.tpost
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (ro *traceElementFork) getTpost() int {
	return ro.tpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (ro *traceElementFork) getTsort() int {
	return ro.tpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (ro *traceElementFork) toString() string {
	return "G" + "," + strconv.Itoa(ro.tpost) + "," + strconv.Itoa(ro.id)
}

/*
 * Update and calculate the vector clock of the element
 * TODO: implement
 */
func (ro *traceElementFork) updateVectorClock() {
	vc.Fork(ro.routine, ro.id, currentVectorClocks)
}
