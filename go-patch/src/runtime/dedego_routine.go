// DEDEGO_FILE_START

package runtime

var DedegoRoutines map[uint64]*[]dedegoTraceElement
var DedegoRoutinesLock *mutex

type DedegoRoutine struct {
	id    uint64
	G     *g
	Trace []dedegoTraceElement
}

/*
 * Create a new dedego routine
 * Params:
 * 	g: the g struct of the routine
 * Return:
 * 	the new dedego routine
 */
func newDedegoRoutine(g *g) *DedegoRoutine {
	routine := &DedegoRoutine{id: GetDedegoRoutineId(), G: g, Trace: make([]dedegoTraceElement, 0)}

	if DedegoRoutinesLock == nil {
		DedegoRoutinesLock = &mutex{}
	}

	lock(DedegoRoutinesLock)

	if DedegoRoutines == nil {
		DedegoRoutines = make(map[uint64]*[]dedegoTraceElement)
	}

	DedegoRoutines[routine.id] = &routine.Trace // Todo: causes warning in race detector

	unlock(DedegoRoutinesLock)

	return routine
}

/*
 * Add an element to the trace of the current routine
 * Params:
 * 	elem: the element to add
 * Return:
 * 	the index of the element in the trace
 */
func (gi *DedegoRoutine) addToTrace(elem dedegoTraceElement) int {
	// TODO: find better way to distinguish between internal and external calls
	if elem.getFile() != "/home/erikkassubek/Uni/dedego/go-patch/bin/main.go" {
		return -1
	}

	// never needed in actual code, without it the compiler tests fail
	if gi == nil {
		return -1
	}
	if gi.Trace == nil {
		gi.Trace = make([]dedegoTraceElement, 0)
	}
	gi.Trace = append(gi.Trace, elem)
	return len(gi.Trace) - 1
}

/*
 * Get the current routine
 * Return:
 * 	the current routine
 */
func currentGoRoutine() *DedegoRoutine {
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

// DEDEGO-FILE-END
