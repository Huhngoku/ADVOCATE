package trace

type traceElementType int

const (
	Atomic traceElementType = iota
	Chan
	Mutex
	Routine
	Select
	Wait
)

/*
 * traceElement is a trace element to save a pre event
 * Fields:
 *   routine (int): The routine id
 *   t (int): The timestamp of the event
 *   vc (vectorClock): The vector clock at the end of the event
 *   elem (traceElement): The corresponding post element
 */
type traceElementPre struct {
	elem     traceElement
	elemType traceElementType
}

/*
 * Get the tpre of the element.
 * Returns:
 *   int: The tpre of the element
 */
func (elem *traceElementPre) getTpre() int {
	return elem.elem.getTpre()
}

/*
 * Get the tpost of the element.
 * Returns:
 *   int: The tpost of the element
 */
func (elem *traceElementPre) getTpost() int {
	return elem.elem.getTpost()
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   float32: The timer of the element
 */
func (elem *traceElementPre) getTsort() float32 {
	return float32(elem.elem.getTpre())
}

/*
 * Get the vector clock at the begin of the event
 * Returns:
 *   vectorClock: The vector clock at the begin of the event
 */
func (elem *traceElementPre) getVpre() *vectorClock {
	return elem.elem.getVpre()
}

/*
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (elem *traceElementPre) getVpost() *vectorClock {
	return elem.elem.getVpost()
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (elem *traceElementPre) getRoutine() int {
	return elem.elem.getRoutine()
}

/*
 * Get the string representation of the element
 * Returns:
 *   string: The string representation of the element
 */
func (elem *traceElementPre) toString() string {
	return "Pre{" + elem.elem.toString() + "}"
}

/*
 * Update the vector clock of the element
 * Params:
 *   vc (vectorClock): The current vector clocks
 */
func (elem *traceElementPre) calculateVectorClock(vc *[]vectorClock) {

}
