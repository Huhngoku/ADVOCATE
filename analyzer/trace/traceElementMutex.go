package trace

import (
	"errors"
	"strconv"
)

// enum for opM
type opMutex int

const (
	LockOp opMutex = iota
	RLockOp
	TryLockOp
	TryRLockOp
	UnlockOp
	RUnlockOp
)

/*
 * traceElementMutex is a trace element for a mutex
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the mutex
 *   rw (bool): Whether the mutex is a read-write mutex
 *   opM (opMutex): The operation on the mutex
 *   suc (bool): Whether the operation was successful (only for trylock else always true)
 *   pos (string): The position of the mutex operation in the code
 */
type traceElementMutex struct {
	routine int
	tpre    int
	tpost   int
	id      int
	rw      bool
	opM     opMutex
	suc     bool
	pos     string
}

func AddTraceElementMutex(routine int, tpre string, tpost string, id string,
	rw string, opM string, suc string, pos string) error {
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

	rw_bool := false
	if rw == "R" {
		rw_bool = true
	}

	var opM_int opMutex = 0
	switch opM {
	case "L":
		opM_int = LockOp
	case "R":
		opM_int = RLockOp
	case "T":
		opM_int = TryLockOp
	case "Y":
		opM_int = TryRLockOp
	case "U":
		opM_int = UnlockOp
	case "N":
		opM_int = RUnlockOp
	default:
		return errors.New("opM is not a valid operation")
	}

	suc_bool, err := strconv.ParseBool(suc)
	if err != nil {
		return errors.New("suc is not a boolean")
	}

	elem := traceElementMutex{
		routine: routine,
		tpre:    tpre_int,
		tpost:   tpost_int,
		id:      id_int,
		rw:      rw_bool,
		opM:     opM_int,
		suc:     suc_bool,
		pos:     pos}

	return addElementToTrace(routine, elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (elem traceElementMutex) getRoutine() int {
	return elem.routine
}

/*
 * Get the tpre of the element.
 * Returns:
 *   int: The tpre of the element
 */
func (elem traceElementMutex) getTpre() int {
	return elem.tpre
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (elem traceElementMutex) getTpost() int {
	return elem.tpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (elem traceElementMutex) toString() string {
	return "M" + "," + strconv.Itoa(elem.tpre) + "," + strconv.Itoa(elem.tpost) +
		strconv.Itoa(elem.id) + "," + strconv.FormatBool(elem.rw) + "," +
		strconv.Itoa(int(elem.opM)) + "," + strconv.FormatBool(elem.suc) + "," +
		elem.pos
}
