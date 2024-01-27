package analysis

var (
	// vc of close on channel
	closeVC  = make(map[int]VectorClock)
	closePos = make(map[int]string)

	// last receive for each routine and each channel
	lastRecvRoutine    = make(map[int]map[int]VectorClock) // routine -> id -> vc
	lastRecvRoutinePos = make(map[int]map[int]string)      // routine -> id -> pos

	// most recent send, used for detection of send on closed
	hasSend                = make(map[int]bool)        // id -> bool
	mostRecentSend         = make(map[int]VectorClock) // id -> vc
	mostRecentSendPosition = make(map[int]string)      // id -> pos

	// most recent send, used for detection of received on closed
	hasReceived               = make(map[int]bool)        // id -> bool
	mostRecentReceive         = make(map[int]VectorClock) // id -> vc
	mostRecentReceivePosition = make(map[int]string)      // id -> pos

	// vector clock for each buffer place in vector clock
	// the map key is the channel id. The slice is used for the buffer positions
	bufferedVCs = make(map[int]([]bufferedVC))
	// the current buffer position
	bufferedVCsCount = make(map[int]int)

	// add on waitGroup
	addVcs = make(map[int]map[int][]VectorClock) // id -> routine -> []vc
	addPos = make(map[int]map[int][]string)      // id -> routine -> []pos

	// done on waitGroup
	doneVcs = make(map[int]map[int][]VectorClock) // id -> routine -> []vc
	donePos = make(map[int]map[int][]string)      // id > routine -> []pos

	// last acquire on mutex for each routine
	lockSet           = make(map[int]map[int]string)      // routine -> id -> string
	mostRecentAcquire = make(map[int]map[int]VectorClock) // routine -> id -> vc  // TODO: do we need to store the operation?

	// vector clocks for last release times
	relW = make(map[int]VectorClock) // id -> vc
	relR = make(map[int]VectorClock) // id -> vc
)
