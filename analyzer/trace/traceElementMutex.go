package trace

import (
	"errors"
	"math"
	"strconv"

	"analyzer/debug"
	vc "analyzer/vectorClock"
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
 *   partner (*traceElementMutex): The partner of the mutex operation
 */
type traceElementMutex struct {
	routine int
	tpre    int
	tpost   int
	// vpre    vc.VectorClock
	vpost   vc.VectorClock
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
func AddTraceElementMutex(routine int, numberOfRoutines int, tpre string,
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
		// vpre:    vc.NewVectorClock(numberOfRoutines),
		vpost: vc.NewVectorClock(numberOfRoutines),
		id:    id_int,
		rw:    rw_bool,
		opM:   opM_int,
		suc:   suc_bool,
		pos:   pos}

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
	return mu.tpre
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (mu *traceElementMutex) getTpost() int {
	return mu.tpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (mu *traceElementMutex) getTsort() int {
	if mu.tpost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return mu.tpost
}

/*
 * Get the vector clock at the begin of the event
 * Returns:
 *   vectorClock: The vector clock at the begin of the event
 */
// func (mu *traceElementMutex) getVpre() *vc.VectorClock {
// 	return &mu.vpre
// }

/*
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (mu *traceElementMutex) getVpost() *vc.VectorClock {
	return &mu.vpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (mu *traceElementMutex) toString() string {
	return "M" + "," + strconv.Itoa(mu.tpre) + "," + strconv.Itoa(mu.tpost) +
		strconv.Itoa(mu.id) + "," + strconv.FormatBool(mu.rw) + "," +
		strconv.Itoa(int(mu.opM)) + "," + strconv.FormatBool(mu.suc) + "," +
		mu.pos
}

// mutex operations, for which no partner has been found yet
var mutexNoPartner []*traceElementMutex

/*
 * Update the vector clock of the trace and element
 */
func (mu *traceElementMutex) updateVectorClock() {
	switch mu.opM {
	case LockOp:
		mu.vpost = vc.Lock(mu.routine, mu.id, currentVectorClocks)
	case RLockOp:
		mu.vpost = vc.RLock(mu.routine, mu.id, currentVectorClocks)
	case TryLockOp:
		if mu.suc {
			mu.vpost = vc.Lock(mu.routine, mu.id, currentVectorClocks)
		}
	case TryRLockOp:
		if mu.suc {
			mu.vpost = vc.RLock(mu.routine, mu.id, currentVectorClocks)
		}
	case UnlockOp:
		mu.vpost = vc.Unlock(mu.routine, mu.id, currentVectorClocks)
	case RUnlockOp:
		mu.vpost = vc.RUnlock(mu.routine, mu.id, currentVectorClocks)
	default:
		err := "Unknown mutex operation: " + mu.toString()
		debug.Log(err, debug.ERROR)
	}
}
