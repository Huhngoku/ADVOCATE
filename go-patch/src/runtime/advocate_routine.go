// ADVOCATE-FILE_START

package runtime

var AdvocateRoutines map[uint64]*[]advocateTraceElement
var AdvocateRoutinesLock mutex = mutex{}

var projectPath string

type AdvocateRoutine struct {
	id    uint64
	G     *g
	Trace []advocateTraceElement
	lock  *mutex
}

/*
 * Create a new advocate routine
 * Params:
 * 	g: the g struct of the routine
 * Return:
 * 	the new advocate routine
 */
func newAdvocateRoutine(g *g) *AdvocateRoutine {
	routine := &AdvocateRoutine{id: GetAdvocateRoutineId(), G: g,
		Trace: make([]advocateTraceElement, 0),
		lock:  &mutex{}}

	lock(&AdvocateRoutinesLock)
	defer unlock(&AdvocateRoutinesLock)

	if AdvocateRoutines == nil {
		AdvocateRoutines = make(map[uint64]*[]advocateTraceElement)
	}

	AdvocateRoutines[routine.id] = &routine.Trace // Todo: causes warning in race detector

	return routine
}

/*
 * Add an element to the trace of the current routine
 * Params:
 * 	elem: the element to add
 * Return:
 * 	the index of the element in the trace
 */
func (gi *AdvocateRoutine) addToTrace(elem advocateTraceElement) int {
	// do nothing if tracer disabled
	if advocateDisabled {
		return -1
	}
	// never needed in actual code, without it the compiler tests fail
	if gi == nil {
		return -1
	}
	lock(gi.lock)
	defer unlock(gi.lock)
	if gi.Trace == nil {
		gi.Trace = make([]advocateTraceElement, 0)
	}
	gi.Trace = append(gi.Trace, elem)
	return len(gi.Trace) - 1
}

func (gi *AdvocateRoutine) getElement(index int) advocateTraceElement {
	lock(gi.lock)
	defer unlock(gi.lock)
	return gi.Trace[index]
}

/*
 * Update an element in the trace of the current routine
 * Params:
 * 	index: the index of the element to update
 * 	elem: the new element
 */
func (gi *AdvocateRoutine) updateElement(index int, elem advocateTraceElement) {
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

	lock(gi.lock)
	defer unlock(gi.lock)
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
 * Get the id of the current routine
 * Return:
 * 	id of the current routine, 0 if current routine is nil
 */
func GetRoutineId() uint64 {
	if currentGoRoutine() == nil {
		return 0
	}
	return currentGoRoutine().id
}

// ADVOCATE-FILE-END
