// DEDEGO_FILE_START

package runtime

var DedegoRoutines map[uint64]*[]dedegoTraceElement
var DedegoRoutinesLock *mutex

var projectPath string

type DedegoRoutine struct {
	id    uint64
	G     *g
	Trace []dedegoTraceElement
}

/*
 * set the project path
 * Params:
 * 	path: the path to the project
 */
func DedegoInit(path string) {
	projectPath = path
}

/*
 * Create a new dedego routine
 * Params:
 * 	g: the g struct of the routine
 * Return:
 * 	the new dedego routine
 */
func newDedegoRoutine(g *g) *DedegoRoutine {
	routine := &DedegoRoutine{id: GetDedegoRoutineId(), G: g,
		Trace: make([]dedegoTraceElement, 0)}

	if DedegoRoutinesLock == nil {
		DedegoRoutinesLock = &mutex{}
	}

	lock(DedegoRoutinesLock)
	defer unlock(DedegoRoutinesLock)

	if DedegoRoutines == nil {
		DedegoRoutines = make(map[uint64]*[]dedegoTraceElement)
	}

	DedegoRoutines[routine.id] = &routine.Trace // Todo: causes warning in race detector

	return routine
}

/*
 * Add an element to the trace of the current routine
 * Params:
 * 	elem: the element to add
 * 	checkInternal: if true, only insert into trace if not internal
 * Return:
 * 	the index of the element in the trace
 */
func (gi *DedegoRoutine) addToTrace(elem dedegoTraceElement,
	checkInternal bool) int {
	if checkInternal && doNotCollectForTrace(elem.getFile()) {
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
