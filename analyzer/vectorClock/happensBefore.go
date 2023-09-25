package vectorClock

type HappensBefore int

const (
	Before HappensBefore = iota
	After
	Concurrent
	None
)
