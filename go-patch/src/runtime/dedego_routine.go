// DEDEGO_FILE_START

package runtime

import "sync/atomic"

type DedegoRoutine struct {
	G       *g
	Trace   []dedegoTraceElement
	counter int32
}

func newDedegoRoutine(g *g) *DedegoRoutine {
	return &DedegoRoutine{G: g, Trace: make([]dedegoTraceElement, 0), counter: 0}
}

func (gi *DedegoRoutine) addToTrace(elem dedegoTraceElement) int {
	gi.Trace = append(gi.Trace, elem)
	return len(gi.Trace) - 1
}

func currentGoRoutine() *DedegoRoutine {
	return getg().goInfo
}

func updateCounter() int32 {
	return atomic.AddInt32(&currentGoRoutine().counter, 1)
}

/*
 * Get the id of the current routine
 * Return:
 * 	id of the current routine
 */
func GetRoutineId() uint64 {
	return currentGoRoutine().G.goid
}

// DEDEGO-FILE-END
