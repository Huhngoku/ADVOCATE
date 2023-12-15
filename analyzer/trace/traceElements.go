package trace

// Interface for trace elements
type traceElement interface {
	getTpre() int
	getTpost() int
<<<<<<< Updated upstream
	getTsort() int
	getRoutine() int
	toString() string
=======
	GetTSort() int
	SetTsort(tsort int)
	SetTSortWithoutNotExecuted(tsort int)
	GetRoutine() int
	GetPos() string
	ToString() string
>>>>>>> Stashed changes
	updateVectorClock()
}
