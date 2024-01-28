package analysis

import "analyzer/logging"

func checkForConcurrentRecv(routine int, id int, pos string, vc map[int]VectorClock) {
	if _, ok := lastRecvRoutine[routine]; !ok {
		lastRecvRoutine[routine] = make(map[int]VectorClockTID)
	}

	lastRecvRoutine[routine][id] = VectorClockTID{vc[routine].Copy(), pos}

	for r, elem := range lastRecvRoutine {
		if r == routine {
			continue
		}

		if elem[id].vc.clock == nil {
			continue
		}

		happensBefore := GetHappensBefore(elem[id].vc, vc[routine])
		if happensBefore == Concurrent {
			found := "Found concurrent Recv on same channel:\n"
			found += "\trecv: " + pos + "\n"
			found += "\trecv : " + lastRecvRoutine[r][id].tID
			logging.Result(found, logging.CRITICAL)
		}
	}
}
