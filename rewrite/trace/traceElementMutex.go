package trace

import (
	"errors"
	"math"
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
 *   partner (*traceElementMutex): The partner of the mutex operation
 */
type traceElementMutex struct {
	routine int
	tPre    int
	tPost   int
	id      int
	rw      bool
	opM     opMutex
	suc     bool
	pos     string
	partner *traceElementMutex
}

/*
 * Create a new mutex trace element
 * Args:
 *   routine (int): The routine id
 *   tPre (string): The timestamp at the start of the event
 *   tPost (string): The timestamp at the end of the event
 *   id (string): The id of the mutex
 *   rw (string): Whether the mutex is a read-write mutex
 *   opM (string): The operation on the mutex
 *   suc (string): Whether the operation was successful (only for trylock else always true)
 *   pos (string): The position of the mutex operation in the code
 */
func AddTraceElementMutex(routine int, tPre string,
	tPost string, id string, rw string, opM string, suc string,
	pos string) error {
	tPreInt, err := strconv.Atoi(tPre)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	tPostInt, err := strconv.Atoi(tPost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	rwBool := false
	if rw == "R" {
		rwBool = true
	}

	var opMInt opMutex
	switch opM {
	case "L":
		opMInt = LockOp
	case "R":
		opMInt = RLockOp
	case "T":
		opMInt = TryLockOp
	case "Y":
		opMInt = TryRLockOp
	case "U":
		opMInt = UnlockOp
	case "N":
		opMInt = RUnlockOp
	default:
		return errors.New("opM is not a valid operation")
	}

	sucBool, err := strconv.ParseBool(suc)
	if err != nil {
		return errors.New("suc is not a boolean")
	}

	elem := traceElementMutex{
		routine: routine,
		tPre:    tPreInt,
		tPost:   tPostInt,
		id:      idInt,
		rw:      rwBool,
		opM:     opMInt,
		suc:     sucBool,
		pos:     pos}

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (mu *traceElementMutex) getRoutine() int {
	return mu.routine
}

/*
 * Get the tpre of the element.
 * Returns:
 *   int: The tpre of the element
 */
func (mu *traceElementMutex) getTpre() int {
	return mu.tPre
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (mu *traceElementMutex) getTpost() int {
	return mu.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (mu *traceElementMutex) getTsort() int {
	if mu.tPost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return mu.tPost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (mu *traceElementMutex) toString() string {
	return "M" + "," + strconv.Itoa(mu.tPre) + "," + strconv.Itoa(mu.tPost) + "," +
		strconv.Itoa(mu.id) + "," + strconv.FormatBool(mu.rw) + "," +
		strconv.Itoa(int(mu.opM)) + "," + strconv.FormatBool(mu.suc) + "," +
		mu.pos
}

// mutex operations, for which no partner has been found yet
var mutexNoPartner []*traceElementMutex
