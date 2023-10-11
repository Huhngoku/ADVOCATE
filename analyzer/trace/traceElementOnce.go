package trace

import (
	"errors"
	"math"
	"strconv"

	vc "analyzer/vectorClock"
)

/*
 * traceElementMutex is a trace element for a once
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the mutex
 *   suc (bool): Whether the operation was successful
 *   pos (string): The position of the mutex operation in the code
 */
type traceElementOnce struct {
	routine int
	tpre    int
	tpost   int
	id      int
	suc     bool
	pos     string
}

/*
 * Create a new mutex trace element
 * Args:
 *   routine (int): The routine id
 *   numberOfRoutines (int): The number of routines in the trace
 *   tpre (string): The timestamp at the start of the event
 *   tpost (string): The timestamp at the end of the event
 *   id (string): The id of the mutex
 *   suc (string): Whether the operation was successful (only for trylock else always true)
 *   pos (string): The position of the mutex operation in the code
 */
func AddTraceElementOnce(routine int, numberOfRoutines int, tpre string,
	tpost string, id string, suc string, pos string) error {
	tpre_int, err := strconv.Atoi(tpre)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	tpost_int, err := strconv.Atoi(tpost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	suc_bool, err := strconv.ParseBool(suc)
	if err != nil {
		return errors.New("suc is not a boolean")
	}

	elem := traceElementOnce{
		routine: routine,
		tpre:    tpre_int,
		tpost:   tpost_int,
		id:      id_int,
		suc:     suc_bool,
		pos:     pos}

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (on *traceElementOnce) getRoutine() int {
	return on.routine
}

/*
 * Get the tpre of the element.
 * Returns:
 *   int: The tpre of the element
 */
func (on *traceElementOnce) getTpre() int {
	return on.tpre
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (on *traceElementOnce) getTpost() int {
	return on.tpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (on *traceElementOnce) getTsort() int {
	if on.tpost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return on.tpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (on *traceElementOnce) toString() string {
	return "O" + "," + strconv.Itoa(on.tpre) + "," + strconv.Itoa(on.tpost) +
		strconv.Itoa(on.id) + "," + strconv.FormatBool(on.suc) + "," +
		on.pos
}

/*
 * Update the vector clock of the trace and element
 */
func (on *traceElementOnce) updateVectorClock() {
	if on.suc {
		vc.DoSuc(on.routine, on.id, currentVectorClocks)
	} else {
		vc.DoFail(on.routine, on.id, currentVectorClocks)
	}
}
