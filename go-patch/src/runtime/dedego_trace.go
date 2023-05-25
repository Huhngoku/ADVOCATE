// DEDEGO-FILE-START

package runtime

type operation int // enum for operation

const (
	opLock operation = iota
	opTryLock
	opUnlock
	opSend
	opRecv
	opClode
)

type prePost int // enum for pre/post
const (
	pre prePost = iota
	post
	none
)

type dedegoTraceElement interface {
	isDedegoTraceElement()
	toString() string
}

// struct to save the trace og one routine
type dedegoRoutineTrace []dedegoTraceElement

// TODO: make routine local
var dedegoTrace map[uint64]dedegoRoutineTrace = make(map[uint64]dedegoRoutineTrace)
var dedegoTraceLock *mutex = new(mutex)

var tracePrintLock *mutex = new(mutex)

/*
 * Return a string representation of the trace
 * Return:
 * 	string representation of the trace
 */
func (t *dedegoRoutineTrace) ToString() string {
	res := "["
	for i, elem := range *t {
		if i != 0 {
			res += ", "
		}
		res += elem.toString()
	}
	res += "]"

	return res
}

/*
 * Add a mutex operation to the trace
 * Args:
 *  elem: element to add to the trace
 */
func addToTrace(elem dedegoTraceElement) {
	lock(dedegoTraceLock)
	defer unlock(dedegoTraceLock)
	routineId := GetRoutineId()
	dedegoTrace[routineId] = append(dedegoTrace[routineId], elem)
}

func PrintTrace() {
	lock(tracePrintLock)
	defer unlock(tracePrintLock)
	for routineId, trace := range dedegoTrace {
		println("Routine", routineId, ":", trace.ToString())
	}
}

// ============================= Mutex =============================

// type to save in the trace for mutexe
type dedegoTraceMutexElement struct {
	id        uint32    // id of the mutex
	operation operation // operation
	suc       bool      // success of the operation, only for tryLock
}

func (elem dedegoTraceMutexElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element
 */
func (elem dedegoTraceMutexElement) toString() string {
	res := "M{" + uint64ToString(uint64(elem.id))

	switch elem.operation {
	case opLock:
		res += ", L"
	case opTryLock:
		res += ", T"
	case opUnlock:
		res += ", U"
	}
	if elem.operation == opTryLock {
		if elem.suc {
			res += ", true"
		} else {
			res += ", false"
		}
	}
	res += "}"
	return res
}

/*
 * Add a mutex lock to the trace
 * Args:
 * 	id: id of the mutex
 */
func DedegoLock(id uint32) {
	elem := dedegoTraceMutexElement{id, opLock, true}
	addToTrace(elem)
}

/*
 * Add a mutex trylock to the trace
 * Args:
 * 	id: id of the mutex
 * 	suc: success of the trylock
 */
func DedegoTryLock(id uint32, suc bool) {
	elem := dedegoTraceMutexElement{id, opTryLock, suc}
	addToTrace(elem)
}

/*
 * Add a mutex unlock to the trace
 * Args:
 * 	id: id of the mutex
 */
func DedegoUnlock(id uint32) {
	elem := dedegoTraceMutexElement{id, opUnlock, true}
	addToTrace(elem)
}

// ============================= Channel =============================

type dedegoTraceChannelElement struct {
	id        uint32    // id of the channel
	prePost   prePost   // pre/post
	operation operation // operation
	partnerId uint32    // id of the channel with wich the communication took place
}

func (elem dedegoTraceChannelElement) isDedegoTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element
 */
func (elem dedegoTraceChannelElement) toString() string {
	res := "C{" + uint64ToString(uint64(elem.id))

	switch elem.operation {
	case opSend:
		res += ", S"
	case opRecv:
		res += ", R"
	case opClode:
		res += ", C"
	}

	switch elem.prePost {
	case pre:
		res += ", pre"
	case post:
		res += ", post"
	}

	res += ", " + uint64ToString(uint64(elem.partnerId)) + "}"
	return res
}

/*
 * Add a channel send to the trace
 * Args:
 * 	id: id of the channel
 * 	isPre: true if pre event, post otherwise
 */
// TODO: add into channel implementation
func DedegoSend(id uint32, isPre bool) {
	p := pre
	if !isPre {
		p = post
	}
	elem := dedegoTraceChannelElement{id, p, opSend, 0}
	addToTrace(elem)
}

/*
 * Add a channel recv to the trace
 * Args:
 * 	id: id of the channel
 * 	isPre: true if pre event, post otherwise
 * 	partnerId: id of the channel with wich the communication took place
 */
// TODO: add into channel implementation
func DedegoRecv(id uint32, isPre bool, partnerId uint32) {
	p := pre
	if !isPre {
		p = post
	}
	elem := dedegoTraceChannelElement{id, p, opRecv, partnerId}
	addToTrace(elem)
}

/*
 * Add a channel close to the trace
 * Args:
 * 	id: id of the channel
 */
// TODO: is there a better way to get the channel id?
func DedegoClose(id uint32) {
	elem := dedegoTraceChannelElement{id, none, opClode, 0}
	addToTrace(elem)
}

// DEDUGO-FILE-END
