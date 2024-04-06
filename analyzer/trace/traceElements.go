package trace

// Interface for trace elements
type TraceElement interface {
	GetID() int
	GetTPre() int
	SetTPre(tPre int)
	getTpost() int
	GetTSort() int
	SetTSort(tSort int)
	SetTSortWithoutNotExecuted(tSort int)
	GetRoutine() int
	GetPos() string
	GetTID() string
	ToString() string
	updateVectorClock()
}
