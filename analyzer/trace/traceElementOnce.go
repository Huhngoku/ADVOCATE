package trace

import (
	"errors"
	"math"
	"strconv"

	"analyzer/analysis"
)

/*
 * traceElementMutex is a trace element for a once
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the mutex
 *   suc (bool): Whether the operation was successful
 *   pos (string): The position of the mutex operation in the code
 */
type TraceElementOnce struct {
	routine int
	tPre    int
	tPost   int
	id      int
	suc     bool
	pos     string
}

/*
 * Create a new mutex trace element
 * Args:
 *   routine (int): The routine id
 *   tPre (string): The timestamp at the start of the event
 *   tPost (string): The timestamp at the end of the event
 *   id (string): The id of the mutex
 *   suc (string): Whether the operation was successful (only for trylock else always true)
 *   pos (string): The position of the mutex operation in the code
 */
func AddTraceElementOnce(routine int, tPre string,
	tPost string, id string, suc string, pos string) error {
	tPreInt, err := strconv.Atoi(tPre)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	tPostInt, err := strconv.Atoi(tPost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	sucBool, err := strconv.ParseBool(suc)
	if err != nil {
		return errors.New("suc is not a boolean")
	}

	elem := TraceElementOnce{
		routine: routine,
		tPre:    tPreInt,
		tPost:   tPostInt,
		id:      idInt,
		suc:     sucBool,
		pos:     pos}

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (on *TraceElementOnce) GetRoutine() int {
	return on.routine
}

/*
 * Get the tpre of the element.
 * Returns:
 *   int: The tpre of the element
 */
func (on *TraceElementOnce) getTpre() int {
	return on.tPre
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (on *TraceElementOnce) getTpost() int {
	return on.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (on *TraceElementOnce) GetTSort() int {
	if on.tPost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return on.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The position of the element
 */
func (on *TraceElementOnce) GetPos() string {
	return on.pos
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tSort (int): The timer of the element
 */
func (on *TraceElementOnce) SetTsort(tSort int) {
	on.tPost = tSort
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tSort (int): The timer of the element
 */
func (on *TraceElementOnce) SetTsortWithoutNotExecuted(tSort int) {
	if on.tPost != 0 {
		on.tPost = tSort
	}
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (on *TraceElementOnce) ToString() string {
	res := "O,"
	res += strconv.Itoa(on.tPre) + ","
	res += strconv.Itoa(on.tPost) + ","
	res += strconv.Itoa(on.id) + ","
	if on.suc {
		res += "t"
	} else {
		res += "f"
	}
	res += "," + on.pos
	return res
}

/*
 * Update the vector clock of the trace and element
 */
func (on *TraceElementOnce) updateVectorClock() {
	if on.suc {
		analysis.DoSuc(on.routine, on.id, currentVectorClocks)
	} else {
		analysis.DoFail(on.routine, on.id, currentVectorClocks)
	}
}
