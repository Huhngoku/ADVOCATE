package trace

import "analyzer/clock"

// Interface for trace elements
type TraceElement interface {
	GetID() int
	GetTPre() int
	SetTPre(tPre int)
	getTpost() int
	GetTSort() int
	SetTSort(tSort int)
	SetT(time int)
	SetTWithoutNotExecuted(tSort int)
	GetRoutine() int
	GetPos() string
	GetTID() string
	ToString() string
	updateVectorClock()
	GetVC() clock.VectorClock
	Copy() TraceElement
}
