package trace

import (
	"errors"
	"strconv"

	"analyzer/debug"
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
 *   vpre (vectorClock): The vector clock at the start of the event
 *   vpost (vectorClock): The vector clock at the end of the event
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
	vpre    vectorClock
	vpost   vectorClock
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
 *   numberOfRoutines (int): The number of routines in the trace
 *   tpre (string): The timestamp at the start of the event
 *   tpost (string): The timestamp at the end of the event
 *   id (string): The id of the mutex
 *   rw (string): Whether the mutex is a read-write mutex
 *   opM (string): The operation on the mutex
 *   suc (string): Whether the operation was successful (only for trylock else always true)
 *   pos (string): The position of the mutex operation in the code
 */
func addTraceElementMutex(routine int, numberOfRoutines int, tpre string,
	tpost string, id string, rw string, opM string, suc string,
	pos string) error {
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
		vpre:    newVectorClock(numberOfRoutines),
		vpost:   newVectorClock(numberOfRoutines),
		id:      id_int,
		rw:      rw_bool,
		opM:     opM_int,
		suc:     suc_bool,
		pos:     pos}

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (elem *traceElementMutex) getRoutine() int {
	return elem.routine
}

/*
 * Get the tpre of the element.
 * Returns:
 *   int: The tpre of the element
 */
func (elem *traceElementMutex) getTpre() int {
	return elem.tpre
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (elem *traceElementMutex) getTpost() int {
	return elem.tpost
}

/*
 * Get the vector clock at the begin of the event
 * Returns:
 *   vectorClock: The vector clock at the begin of the event
 */
func (elem *traceElementMutex) getVpre() *vectorClock {
	return &elem.vpre
}

/*
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (elem *traceElementMutex) getVpost() *vectorClock {
	return &elem.vpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (elem *traceElementMutex) toString() string {
	return "M" + "," + strconv.Itoa(elem.tpre) + "," + strconv.Itoa(elem.tpost) +
		strconv.Itoa(elem.id) + "," + strconv.FormatBool(elem.rw) + "," +
		strconv.Itoa(int(elem.opM)) + "," + strconv.FormatBool(elem.suc) + "," +
		elem.pos
}

// mutex operations, for which no partner has been found yet
var mutexNoPartner []*traceElementMutex

/*
 * Find pairs of lock and unlock operations. If a partner is found, the partner
 * is set in the element.
 * The functions assumes, that the trace list is sorted by tpost
 */
func (elem *traceElementMutex) findPartner() {
	// check if the element should have a partner
	if elem.tpost == 0 || !elem.suc {
		debug.Log("Mutex operation "+elem.toString()+" has not executed", 3)
		return
	}

	found := false
	if elem.opM == LockOp || elem.opM == RLockOp || elem.opM == TryLockOp {
		debug.Log("Add mutex lock operations "+elem.toString()+" to mutexNoPartner", 3)
		// add lock operations to list of locks without partner
		mutexNoPartner = append(mutexNoPartner, elem)
		found = true // set to true to prevent panic
	} else if elem.opM == UnlockOp || elem.opM == RUnlockOp {
		// for unlock operations, check find the last lock operation
		// on the same mutex
		for i := len(mutexNoPartner) - 1; i >= 0; i-- {
			lock := mutexNoPartner[i]
			if elem.id != lock.id {
				continue
			}
			if lock.opM == UnlockOp || lock.opM == RUnlockOp {
				debug.Log("Two consecutive lock on the same channel without unlock in between: "+elem.toString()+lock.toString(), 1)
			}
			debug.Log("Found partner for mutex operation "+lock.toString()+" <-> "+elem.toString(), 3)
			elem.partner = lock
			lock.partner = elem
			debug.Log("Remove mutex lock operation "+lock.toString()+" from mutexNoPartner", 3)
			mutexNoPartner = append(mutexNoPartner[:i], mutexNoPartner[i+1:]...)
			found = true
			break
		}
	} else {
		panic("Unknown mutex operation")
	}

	if !found {
		debug.Log("Unlock "+elem.toString()+" without prior lock", 1)
	}
}
