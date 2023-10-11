package trace

// Interface for trace elements
type traceElement interface {
	getTpre() int
	getTpost() int
	getTsort() int
	getRoutine() int
	toString() string
	updateVectorClock()
}
