// DEDEGO_FILE_START

package runtime

type DedegoRoutine struct {
	id    uint64
	G     *g
	Trace []dedegoTraceElement
}

func newDedegoRoutine(g *g) *DedegoRoutine {
	return &DedegoRoutine{id: GetDedegoRoutineId(), G: g, Trace: make([]dedegoTraceElement, 0)}
}

func (gi *DedegoRoutine) addToTrace(elem dedegoTraceElement) int {
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

func currentGoRoutine() *DedegoRoutine {
	return getg().goInfo
}

/*
 * Get the id of the current routine
 * Return:
 * 	id of the current routine
 */
func GetRoutineId() uint64 {
	return currentGoRoutine().id
}

// DEDEGO-FILE-END
