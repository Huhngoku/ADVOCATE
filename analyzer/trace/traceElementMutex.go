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
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the mutex
 *   rw (bool): Whether the mutex is a read-write mutex
 *   opM (opMutex): The operation on the mutex
 *   exec (bool): The execution status of the operation
 *   suc (bool): Whether the operation was successful (only for trylock else always true)
 *   pos (string): The position of the mutex operation in the code
 */
type traceElementMutex struct {
	tpre  int
	tpost int
	id    int
	rw    bool
	opM   opMutex
	exec  bool
	suc   bool
	pos   string
}

func AddTraceElementMutex(routine int, tpre string, tpost string, id string,
	rw string, opM string, exec string, suc string, pos string) error {
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

	exec_bool, err := strconv.ParseBool(exec)
	if err != nil {
		return errors.New("exec is not a boolean")
	}

	suc_bool, err := strconv.ParseBool(suc)
	if err != nil {
		return errors.New("suc is not a boolean")
	}

	elem := traceElementMutex{tpre_int, tpost_int, id_int, rw_bool, opM_int, exec_bool, suc_bool, pos}

	return addElementToTrace(routine, elem)
}

func (elem traceElementMutex) getSimpleString() string {
	return "M" + "," + strconv.Itoa(elem.tpre) + "," + strconv.Itoa(elem.tpost) +
		strconv.Itoa(elem.id) + "," + strconv.FormatBool(elem.rw) + "," +
		strconv.Itoa(int(elem.opM))
}
