package runtime

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
 * Add an atomic operation to the trace
 * Args:
 * 	index: index of the atomic event in advocateAtomicMap
 */
func AdvocateAtomic(index uint64) {
	timer := GetNextTimeStep()
	elem := advocateAtomicElement{index: index, timer: timer}
	insertIntoTrace(elem)
}
