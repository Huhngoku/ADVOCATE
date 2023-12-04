package trace

// Interface for trace elements
type TraceElement interface {
	getTpre() int
	getTpost() int
	GetTSort() int
	SetTsortWithoutNotExecuted(tsort int)
	SetTsort(tsort int)
	GetRoutine() int
	GetPos() string
	ToString() string
}
