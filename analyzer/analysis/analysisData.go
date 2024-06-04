package analysis

import "analyzer/clock"

type VectorClockTID struct {
	Vc      clock.VectorClock
	TID     string
	Routine int
}

type VectorClockTID2 struct {
	routine  int
	id       int
	vc       clock.VectorClock
	tID      string
	typeVal  int
	val      int
	buffered bool
	sel      bool
	selID    int
}

type VectorClockTID3 struct {
	Routine int
	TID     string
	Vc      clock.VectorClock
	Val     int
}

type allSelectCase struct {
	selectID int            // select id
	chanID   int            // channel id
	vcTID    VectorClockTID // vector clock and tID
	send     bool           // true: send, false: receive
	buffered bool           // true: buffered, false: unbuffered
	partner  bool           // true: partner found, false: no partner found
}

var (
	// analysis cases to run
	analysisCases = make(map[string]bool)

	// vc of close on channel
	closeData = make(map[int]VectorClockTID3) // id -> vcTID3 val = objID

	// last receive for each routine and each channel
	lastRecvRoutine = make(map[int]map[int]VectorClockTID) // routine -> id -> vcTID

	// most recent send, used for detection of send on closed
	hasSend        = make(map[int]bool)                    // id -> bool
	mostRecentSend = make(map[int]map[int]VectorClockTID3) // routine -> id -> vcTID

	// most recent send, used for detection of received on closed
	hasReceived       = make(map[int]bool)                    // id -> bool
	mostRecentReceive = make(map[int]map[int]VectorClockTID3) // routine -> id -> vcTID3, val = objID

	// vector clock for each buffer place in vector clock
	// the map key is the channel id. The slice is used for the buffer positions
	bufferedVCs = make(map[int]([]bufferedVC))
	// the current buffer position
	bufferedVCsCount = make(map[int]int)

	// add on waitGroup
	wgAdd = make(map[int]map[int][]VectorClockTID) // id -> routine -> []vcTID

	// done on waitGroup
	wgDone = make(map[int]map[int][]VectorClockTID) // id -> routine -> []vcTID

	// wait on waitGroup
	// wgWait = make(map[int]map[int][]VectorClockTID) // id -> routine -> []vcTID

	// last acquire on mutex for each routine
	lockSet                = make(map[int]map[int]string)         // routine -> id -> string
	mostRecentAcquire      = make(map[int]map[int]VectorClockTID) // routine -> id -> vcTID  // TODO: do we need to store the operation?
	mostRecentAcquireTotal = make(map[int]VectorClockTID3)        // id -> vcTID

	// vector clocks for last release times
	relW = make(map[int]clock.VectorClock) // id -> vc
	relR = make(map[int]clock.VectorClock) // id -> vc

	// for leak check
	leakingChannels = make(map[int][]VectorClockTID2) // id -> vcTID

	// for check of select without partner
	// store all select cases
	selectCases = make([]allSelectCase, 0)
)

// InitAnalysis initializes the analysis cases
func InitAnalysis(analysisCasesMap map[string]bool) {
	analysisCases = analysisCasesMap
}
