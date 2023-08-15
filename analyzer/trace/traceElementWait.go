package trace

import (
	"errors"
	"strconv"
)

// enum for opW
type opW int

const (
	ChangeOp opW = iota
	WaitOp
)

/*
* traceElementWait is a trace element for a wait group statement
* Fields:
*   tpre (int): The timestamp at the start of the event
*   tpost (int): The timestamp at the end of the event
*   id (int): The id of the wait group
*   opW (opW): The operation on the wait group
*   exec (bool): The execution status of the operation
*   delta (int): The delta of the wait group
*   val (int): The value of the wait group
*   pos (string): The position of the wait group in the code
 */
type traceElementWait struct {
	routine int
	tpre    int
	tpost   int
	id      int
	opW     opW
	exec    bool
	delta   int
	val     int
	pos     string
}

func AddTraceElementWait(routine int, tpre string, tpost string, id string,
	opW string, exec string, delta string, val string, pos string) error {
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

	opW_op := ChangeOp
	if opW == "W" {
		opW_op = WaitOp
	} else if opW != "A" {
		return errors.New("op is not a valid operation")
	}

	exec_bool, err := strconv.ParseBool(exec)
	if err != nil {
		return errors.New("exec is not a boolean")
	}

	delta_int, err := strconv.Atoi(delta)
	if err != nil {
		return errors.New("delta is not an integer")
	}

	val_int, err := strconv.Atoi(val)
	if err != nil {
		return errors.New("val is not an integer")
	}

	elem := traceElementWait{routine, tpre_int, tpost_int, id_int, opW_op,
		exec_bool, delta_int, val_int, pos}

	return addElementToTrace(routine, elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (elem traceElementWait) getRoutine() int {
	return elem.routine
}

/*
 * Get the timestamp at the start of the event
 * Returns:
 *   int: The timestamp at the start of the event
 */
func (elem traceElementWait) getTpre() int {
	return elem.tpre
}

/*
 * Get the timestamp at the end of the event
 * Returns:
 *   int: The timestamp at the end of the event
 */
func (elem traceElementWait) getTpost() int {
	return elem.tpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (elem traceElementWait) getSimpleString() string {
	return "W" + strconv.Itoa(elem.id) + "," + strconv.Itoa(elem.tpre) + "," +
		strconv.Itoa(elem.tpost) + "," + "," +
		strconv.Itoa(elem.delta)
}
