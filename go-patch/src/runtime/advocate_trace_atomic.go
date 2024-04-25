package runtime

import at "runtime/internal/atomic"

/*
 * Add an atomic operation to the trace
 * Args:
 * 	index: index of the atomic event in advocateAtomicMap
 */
func AdvocateAtomicPre(index uint64) {
	timer := GetNextTimeStep()

	elem := "A," + uint64ToString(timer) + "," + uint64ToString(index) + ",-"

	// // elem := advocateAtomicElement{index: index, timer: timer}
	insertIntoTrace(elem)
}

/*
 * Update the atomic operation in the trace after received
 * Args:
 * 	atomic: atomic operation to add to the trace
 */
func AdvocateAtomicPost(atomic at.AtomicElem) {
	lock(&advocateAtomicMapLock)
	advocateAtomicMap[atomic.Index] = advocateAtomicMapElem{
		addr:      atomic.Addr,
		operation: atomic.Operation,
	}
	unlock(&advocateAtomicMapLock)
}

/*
 * Add the id and operation to an atomic operation
 * Args:
 * 	elem: the atomic operation
 * Return:
 * 	the atomic operation with the id and operation
 */
func addAtomicInfo(elem string) string {
	split := splitStringAtCommas(elem, []int{2, 3}) // A,[tpre] - id - operation
	if split[2] != "-" {
		return elem
	}

	index := uint64(stringToInt(split[1]))

	lock(&advocateAtomicMapLock)
	mapElement := advocateAtomicMap[index]
	unlock(&advocateAtomicMapLock)
	lock(&advocateAtomicMapToIDLock)
	if _, ok := advocateAtomicMapToID[mapElement.addr]; !ok {
		advocateAtomicMapToID[mapElement.addr] = advocateAtomicMapIDCounter
		advocateAtomicMapIDCounter++
	}
	id := advocateAtomicMapToID[mapElement.addr]
	unlock(&advocateAtomicMapToIDLock)

	operation := ""
	switch mapElement.operation {
	case at.LoadOp:
		operation = "L"
	case at.StoreOp:
		operation = "S"
	case at.AddOp:
		operation = "A"
	case at.SwapOp:
		operation = "W"
	case at.CompSwapOp:
		operation = "C"
	default:
		operation = "U"
	}

	split[1] = uint64ToString(id)
	split[2] = operation

	return mergeString(split)
}
