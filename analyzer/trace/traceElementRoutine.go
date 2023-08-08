package trace

import (
	"errors"
	"strconv"
)

/*
 * traceElementRoutine is a trace element for a go statement
 * Fields:
 *   tpre (int): The timestamp at the start of the event
 *   id (int): The id of the new go statement
 */
type traceElementRoutine struct {
	tpre int
	id   int
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

	elem := traceElementRoutine{tpre_int, id_int}
	return addElementToTrace(routine, elem)
}

func (elem traceElementRoutine) getSimpleString() string {
	return "G" + "," + strconv.Itoa(elem.tpre) + "," + strconv.Itoa(elem.id)
}
