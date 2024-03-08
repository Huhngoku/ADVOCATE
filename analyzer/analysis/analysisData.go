package analysis

type VectorClockTID struct {
	vc  VectorClock
	tID string
}

type VectorClockTID2 struct {
	vc      VectorClock
	tID     string
	typeVal int
	val     int
}

type VectorClockTID3 struct {
	vc       VectorClock
	tID      string
	buffered bool
}

type untriggeredSelectCase struct {
	id    int            // channel id
	vcTID VectorClockTID // vector clock and tID
}

var (
	// analysis cases to run
	analysisCases = make(map[string]bool)

	// vc of close on channel
	closeData = make(map[int]VectorClockTID)
	closeRout = make(map[int]int)

	// last receive for each routine and each channel
	lastRecvRoutine = make(map[int]map[int]VectorClockTID) // routine -> id -> vcTID

	// most recent send, used for detection of send on closed
	hasSend        = make(map[int]bool)           // id -> bool
	mostRecentSend = make(map[int]VectorClockTID) // id -> vcTID

	// most recent send, used for detection of received on closed
	hasReceived       = make(map[int]bool)           // id -> bool
	mostRecentReceive = make(map[int]VectorClockTID) // id -> vcTID

	// vector clock for each buffer place in vector clock
	// the map key is the channel id. The slice is used for the buffer positions
	bufferedVCs = make(map[int]([]bufferedVC))
	// the current buffer position
	bufferedVCsCount = make(map[int]int)

	// add on waitGroup
	addWait = make(map[int]map[int][]VectorClockTID) // id -> routine -> []vcTID

	// done on waitGroup
	doneWait = make(map[int]map[int][]VectorClockTID) // id -> routine -> []vcTID

	// last acquire on mutex for each routine
	lockSet           = make(map[int]map[int]string)         // routine -> id -> string
	mostRecentAcquire = make(map[int]map[int]VectorClockTID) // routine -> id -> vcTID  // TODO: do we need to store the operation?

	// vector clocks for last release times
	relW = make(map[int]VectorClock) // id -> vc
	relR = make(map[int]VectorClock) // id -> vc

	// for leak check
	leakingChannels = make(map[int][]VectorClockTID2) // id -> vcTID

	// for check of select without partner
	// not triggered
	selectCasesSend = make(map[int][]VectorClockTID3) // chanID -> []vcTID3
	selectCasesRecv = make(map[int][]VectorClockTID3) // chanID -> []vcTID3
)

// InitAnalysis initializes the analysis cases
func InitAnalysis(analysisCasesMap map[string]bool) {
	analysisCases = analysisCasesMap
}
