package trace

import (
	"errors"
	"math"
	"strconv"

	"analyzer/analysis"
	"analyzer/logging"
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
 * TraceElementMutex is a trace element for a mutex
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the mutex
 *   rw (bool): Whether the mutex is a read-write mutex
 *   opM (opMutex): The operation on the mutex
 *   suc (bool): Whether the operation was successful (only for trylock else always true)
 *   pos (string): The position of the mutex operation in the code
 *   tID (string): The id of the trace element, contains the position and the tpre
 *   partner (*TraceElementMutex): The partner of the mutex operation
 */
type TraceElementMutex struct {
	routine int
	tPre    int
	tPost   int
	id      int
	rw      bool
	opM     opMutex
	suc     bool
	pos     string
	tID     string
	partner *TraceElementMutex
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

	tIDStr := pos + "@" + strconv.Itoa(tPreInt)

	elem := TraceElementMutex{
		routine: routine,
		tPre:    tPreInt,
		tPost:   tPostInt,
		id:      idInt,
		rw:      rwBool,
		opM:     opMInt,
		suc:     sucBool,
		pos:     pos,
		tID:     tIDStr,
	}

	return addElementToTrace(&elem)
}

/*
 * Get the id of the element
 * Returns:
 *   int: The id of the element
 */
func (mu *TraceElementMutex) GetID() int {
	return mu.id
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (mu *TraceElementMutex) GetRoutine() int {
	return mu.routine
}

/*
 * Get the tpre of the element.
 * Returns:
 *   int: The tpre of the element
 */
func (mu *TraceElementMutex) getTpre() int {
	return mu.tPre
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (mu *TraceElementMutex) getTpost() int {
	return mu.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (mu *TraceElementMutex) GetTSort() int {
	if mu.tPost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return mu.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The position of the element
 */
func (mu *TraceElementMutex) GetPos() string {
	return mu.pos
}

/*
 * Get the tID of the element.
 * Returns:
 *   string: The tID of the element
 */
func (mu *TraceElementMutex) GetTID() string {
	return mu.tID
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tsort (int): The timer of the element
 */
func (mu *TraceElementMutex) SetTsort(tSort int) {
	mu.tPost = tSort
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tsort (int): The timer of the element
 */
func (mu *TraceElementMutex) SetTsortWithoutNotExecuted(tSort int) {
	if mu.tPost != 0 {
		mu.tPost = tSort
	}
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (mu *TraceElementMutex) ToString() string {
	res := "M,"
	res += strconv.Itoa(mu.tPre) + "," + strconv.Itoa(mu.tPost) + ","
	res += strconv.Itoa(mu.id) + ","

	if mu.rw {
		res += "R,"
	} else {
		res += "-,"
	}

	switch mu.opM {
	case LockOp:
		res += "L"
	case RLockOp:
		res += "R"
	case TryLockOp:
		res += "T"
	case TryRLockOp:
		res += "Y"
	case UnlockOp:
		res += "U"
	case RUnlockOp:
		res += "N"
	}

	if mu.suc {
		res += ",t"
	} else {
		res += ",f"
	}
	res += "," + mu.pos
	return res
}

// mutex operations, for which no partner has been found yet
var mutexNoPartner []*TraceElementMutex

/*
 * Update the vector clock of the trace and element
 */
func (mu *TraceElementMutex) updateVectorClock() {
	switch mu.opM {
	case LockOp:
		analysis.Lock(mu.routine, mu.id, currentVCHb, currentVCWmhb, mu.tID, mu.tPost)
		analysis.AnalysisDeadlockMutexLock(mu.id, mu.routine, mu.rw, false, currentVCHb[mu.routine], mu.tPost)
	case RLockOp:
		analysis.RLock(mu.routine, mu.id, currentVCHb, currentVCWmhb, mu.tID, mu.tPost)
		analysis.AnalysisDeadlockMutexLock(mu.id, mu.routine, mu.rw, true, currentVCHb[mu.routine], mu.tPost)
	case TryLockOp:
		if mu.suc {
			analysis.Lock(mu.routine, mu.id, currentVCHb, currentVCWmhb, mu.tID, mu.tPost)
			analysis.AnalysisDeadlockMutexLock(mu.id, mu.routine, mu.rw, false, currentVCHb[mu.routine], mu.tPost)
		}
	case TryRLockOp:
		if mu.suc {
			analysis.RLock(mu.routine, mu.id, currentVCHb, currentVCWmhb, mu.tID, mu.tPost)
			analysis.AnalysisDeadlockMutexLock(mu.id, mu.routine, mu.rw, true, currentVCHb[mu.routine], mu.tPost)
		}
	case UnlockOp:
		analysis.Unlock(mu.routine, mu.id, currentVCHb, mu.tPost)
		analysis.AnalysisDeadlockMutexUnLock(mu.id, mu.routine, mu.tPost)
	case RUnlockOp:
		analysis.RUnlock(mu.routine, mu.id, currentVCHb, mu.tPost)
		analysis.AnalysisDeadlockMutexUnLock(mu.id, mu.routine, mu.tPost)
	default:
		err := "Unknown mutex operation: " + mu.ToString()
		logging.Debug(err, logging.ERROR)
	}
}

func (mu *TraceElementMutex) updateVectorClockAlt() {
	currentVCHb[mu.routine] = currentVCHb[mu.routine].Inc(mu.routine)
}
