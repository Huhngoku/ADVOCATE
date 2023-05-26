// DEDEGO_FILE_START

package runtime

type GoInfo struct {
	G     *g
	Trace []dedegoTraceElement
}

func NewGoInfo(g *g) *GoInfo {
	return &GoInfo{G: g, Trace: make([]dedegoTraceElement, 0)}
}

func (gi *GoInfo) AddToTrace(elem dedegoTraceElement) {
	gi.Trace = append(gi.Trace, elem)
}

func CurrentGoInfo() *GoInfo {
	return getg().goInfo
}

/*
 * Get the id of the current routine
 * Return:
 * 	id of the current routine
 */
func GetRoutineId() uint64 {
	return CurrentGoInfo().G.goid
}

// DEDEGO-FILE-END
