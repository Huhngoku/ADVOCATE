package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
)

func checkForConcurrentRecv(routine int, id int, pos string, vc map[int]clock.VectorClock) {
	if _, ok := lastRecvRoutine[routine]; !ok {
		lastRecvRoutine[routine] = make(map[int]VectorClockTID)
	}

	lastRecvRoutine[routine][id] = VectorClockTID{vc[routine].Copy(), pos, routine}

	for r, elem := range lastRecvRoutine {
		if r == routine {
			continue
		}

		if elem[id].Vc.GetClock() == nil {
			continue
		}

		happensBefore := clock.GetHappensBefore(elem[id].Vc, vc[routine])
		if happensBefore == clock.Concurrent {
			found := "Found concurrent Recv on same channel:\n"
			found += "\trecv: " + pos + "\n"
			found += "\trecv: " + lastRecvRoutine[r][id].TID
			logging.Result(found, logging.CRITICAL)
		}
	}
}
