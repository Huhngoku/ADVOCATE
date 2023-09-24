package trace

import (
	"errors"
	"strconv"
)

/*
 * traceElementRoutine is a trace element for a go statement
 * Fields:
 *   routine (int): The routine id
 *   tpost (int): The timestamp at the end of the event
 *   vpost (vectoClock): The vector clock at the end of the event
 *   id (int): The id of the new go statement
 */
type traceElementRoutine struct {
	routine int
	tpost   int
	vpost   vectorClock
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
func addTraceElementRoutine(routine int, numberOfRoutines int, tpost string,
	id string) error {
	tpost_int, err := strconv.Atoi(tpost)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	elem := traceElementRoutine{
		routine: routine,
		tpost:   tpost_int,
		vpost:   newVectorClock(numberOfRoutines),
		id:      id_int}
	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (elem *traceElementRoutine) getRoutine() int {
	return elem.routine
}

/*
 * Get the tpre of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpre of the element
 */
func (elem *traceElementRoutine) getTpre() int {
	return elem.tpost
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (elem *traceElementRoutine) getTpost() int {
	return elem.tpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (elem *traceElementRoutine) getTsort() float32 {
	return float32(elem.tpost)
}

/*
 * Get the vector clock at the begin of the event. It is equal to the vector clock
 * at the end of the event.
 * Returns:
 *   vectorClock: The vector clock at the begin of the event
 */
func (elem *traceElementRoutine) getVpre() *vectorClock {
	return &elem.vpost
}

/*
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (elem *traceElementRoutine) getVpost() *vectorClock {
	return &elem.vpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (elem *traceElementRoutine) toString() string {
	return "G" + "," + strconv.Itoa(elem.tpost) + "," + strconv.Itoa(elem.id)
}

/*
 * Update and calculate the vector clock of the element
 * Args:
 *   vc (vectorClock): The current vector clocks
 * TODO: implement
 */
func (elem *traceElementRoutine) calculateVectorClock(vc *[]vectorClock) {
}
