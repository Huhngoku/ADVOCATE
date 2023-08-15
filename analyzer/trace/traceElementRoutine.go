package trace

import (
	"errors"
	"strconv"
)

/*
 * traceElementRoutine is a trace element for a go statement
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   id (int): The id of the new go statement
 */
type traceElementRoutine struct {
	routine int
	tpre    int
	id      int
}

func AddTraceElementRoutine(routine int, tpre string, id string) error {
	tpre_int, err := strconv.Atoi(tpre)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	elem := traceElementRoutine{routine, tpre_int, id_int}
	return addElementToTrace(routine, elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (elem traceElementRoutine) getRoutine() int {
	return elem.routine
}

/*
 * Get the tpre of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpre of the element
 */
func (elem traceElementRoutine) getTpre() int {
	return elem.tpre
}

/*
 * Get the tpost of the element. For atomic elements, tpre and tpost are the same
 * Returns:
 *   int: The tpost of the element
 */
func (elem traceElementRoutine) getTpost() int {
	return elem.tpre
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (elem traceElementRoutine) getSimpleString() string {
	return "G" + "," + strconv.Itoa(elem.tpre) + "," + strconv.Itoa(elem.id)
}
