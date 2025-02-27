// ADVOCATE-FILE_START

package runtime

import (
	"fmt"
)

var AdvocateRoutines map[uint64]*AdvocateRoutine
var AdvocateRoutinesLock = mutex{}

var projectPath string

var atomicRecordingDisabled = false

/*
 * AdvocateRoutine is a struct to store the trace of a routine
 * id: the id of the routine
 * G: the g struct of the routine
 * Trace: the trace of the routine
 */
type AdvocateRoutine struct {
	id      uint64
	G       *g
	Trace   []string
	Atomics []string
	// lock    *mutex
	newEvents []string
}

/*
 * Create a new advocate routine
 * Params:
 * 	g: the g struct of the routine
 * Return:
 * 	the new advocate routine
 */
func newAdvocateRoutine(g *g) *AdvocateRoutine {
	routine := &AdvocateRoutine{id: GetAdvocateRoutineID(), G: g,
		Trace:     make([]string, 0),
		Atomics:   make([]string, 0),
		newEvents: make([]string, 0)}

	lock(&AdvocateRoutinesLock)
	defer unlock(&AdvocateRoutinesLock)

	if AdvocateRoutines == nil {
		AdvocateRoutines = make(map[uint64]*AdvocateRoutine)
	}

	AdvocateRoutines[routine.id] = routine

	return routine
}

/*
 * Add an element to the trace of the current routine
 * Params:
 * 	elem: the element to add
 * Return:
 * 	the index of the element in the trace
 */
func (gi *AdvocateRoutine) addToTrace(elem string) int {
	// do nothing if tracer disabled
	if advocateDisabled {
		return -1
	}

	// do nothing while trace writing disabled
	// this is used to avoid writing to the trace, while the trace is written
	// to the file in case of a too high memory usage
	// for advocateTraceWritingDisabled {
	// 	slowExecution()
	// }

	// never needed in actual code, without it the compiler tests fail
	if gi == nil {
		return -1
	}

	gi.Trace = append(gi.Trace, elem)
	return len(gi.Trace) - 1
}

// MY_CHANGES

func (gi *AdvocateRoutine) addToTrace_newEvents(elem string) int {
	// do nothing if tracer disabled
	if advocateDisabled {
		return -1
	}

	// do nothing while trace writing disabled
	// this is used to avoid writing to the trace, while the trace is written
	// to the file in case of a too high memory usage
	// for advocateTraceWritingDisabled {
	// 	slowExecution()
	// }

	// never needed in actual code, without it the compiler tests fail
	if gi == nil {
		return -1
	}

	gi.Trace = append(gi.newEvents, elem)
	return len(gi.newEvents) - 1
}

func (gi *AdvocateRoutine) writeLockEvent(l *mutex) string {
	event := ""
	if l.lastUnlock > 0 {
		event = fmt.Sprintf("lock, key: %d, lastRelease: %d", &l.key, &l.lastUnlock)
	} else {
		event = fmt.Sprintf("lock, key: %d, lastRelease: None")
	}
	return event
}

func (gi *AdvocateRoutine) writeUnLockEvent(l *mutex) string {
	event := ""
	if l.lastUnlock > 0 {
		event = fmt.Sprintf("Unlock, key: %d, lastRelease: %d", &l.key, &l.lastUnlock)
	} else {
		throw("Unlocked with no release time")
	}
	return event
}

//MY_CHANGES

/*
 * Ignore the atomic operations. Use if not enough memory is available.
 */
func IgnoreAtomicOperations() {
	atomicRecordingDisabled = true
	sum := 0
	lock(&AdvocateRoutinesLock)
	for _, routine := range AdvocateRoutines {
		println("Delete ", len(routine.Atomics), " atomic operations")
		sum += len(routine.Atomics)
		routine.Atomics = nil
	}
	unlock(&AdvocateRoutinesLock)
	println("Deleted ", sum, " atomic operations")
	GC() // run the garbage collector
}

/*
 * Get if atomic operations are ignored
 */
func GetIgnoreAtomicOperations() bool {
	return atomicRecordingDisabled
}

/*
 * Add an atomic operation to the trace of the current routine
 * Params:
 * 	elem: the element to add
 */
func (gi *AdvocateRoutine) addAtomicToTrace(elem string) {
	if advocateDisabled {
		return
	}

	if gi == nil {
		return
	}

	// delete atomic operations if disabled
	if atomicRecordingDisabled {
		// if gi.Atomics != nil {
		// 	println("Delete ", len(gi.Atomics), " atomic operations")
		// }
		// gi.Atomics = nil
		return
	}

	gi.Atomics = append(gi.Atomics, elem)
}

func (gi *AdvocateRoutine) getElement(index int) string {
	return gi.Trace[index]
}

/*
 * Update an element in the trace of the current routine
 * Params:
 * 	index: the index of the element to update
 * 	elem: the new element
 */
func (gi *AdvocateRoutine) updateElement(index int, elem string) {
	if advocateDisabled {
		return
	}

	if gi == nil {
		return
	}

	if gi.Trace == nil {
		panic("Tried to update element in nil trace")
	}

	if index >= len(gi.Trace) {
		panic("Tried to update element out of bounds")
	}

	gi.Trace[index] = elem
}

/*
 * Get the current routine
 * Return:
 * 	the current routine
 */
func currentGoRoutine() *AdvocateRoutine {
	return getg().goInfo
}

/*
 * GetRoutineID gets the id of the current routine
 * Return:
 * 	id of the current routine, 0 if current routine is nil
 */
func GetRoutineID() uint64 {
	if currentGoRoutine() == nil {
		return 0
	}
	return currentGoRoutine().id
}

/*
 * DisableAtomicRecording disables the recording of atomic operations
 */
func DisableAtomicRecording() {
	atomicRecordingDisabled = true
}

// ADVOCATE-FILE-END
