package runtime

type advocateAtomicElement struct {
	timer     uint64 // global timer
	index     uint64 // index of the atomic event in advocateAtomicMap
	operation int    // type of operation
}

func (elem advocateAtomicElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "A,'addr'"
 *    'addr' (number): address of the atomic variable
 */
// enum for atomic operation, must be the same as in advocate_atomic.go
const (
	LoadOp = iota
	StoreOp
	AddOp
	SwapOp
	CompSwapOp
)

func (elem advocateAtomicElement) toString() string {
	lock(&advocateAtomicMapLock)
	mapElement := advocateAtomicMap[elem.index]
	unlock(&advocateAtomicMapLock)
	lock(&advocateAtomicMapToIDLock)
	if _, ok := advocateAtomicMapToID[mapElement.addr]; !ok {
		advocateAtomicMapToID[mapElement.addr] = advocateAtomicMapIDCounter
		advocateAtomicMapIDCounter++
	}
	id := advocateAtomicMapToID[mapElement.addr]
	unlock(&advocateAtomicMapToIDLock)

	res := "A," + uint64ToString(elem.timer) + "," +
		uint64ToString(id) + ","
	switch mapElement.operation {
	case LoadOp:
		res += "L"
	case StoreOp:
		res += "S"
	case AddOp:
		res += "A"
	case SwapOp:
		res += "W"
	case CompSwapOp:
		res += "C"
	default:
		res += "U"
	}
	return res
}

/*
 * Get the operation
 */
func (elem advocateAtomicElement) getOperation() Operation {
	return OperationAtomic
}

/*
 * Get the file
 */
func (elem advocateAtomicElement) getFile() string {
	return ""
}

/*
 * Get the line
 */
func (elem advocateAtomicElement) getLine() int {
	return 0
}

/*
 * Add an atomic operation to the trace
 * Args:
 * 	index: index of the atomic event in advocateAtomicMap
 */
func AdvocateAtomic(index uint64) {
	timer := GetAdvocateCounter()
	elem := advocateAtomicElement{index: index, timer: timer}
	insertIntoTrace(elem)
}
