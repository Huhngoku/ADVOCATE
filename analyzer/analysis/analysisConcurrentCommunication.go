package analysis

import "analyzer/logging"

func checkForConcurrentRecv(routine int, id int, pos string, vc map[int]VectorClock) {
	if _, ok := lastRecvRoutine[routine]; !ok {
		lastRecvRoutine[routine] = make(map[int]VectorClock)
		lastRecvRoutinePos[routine] = make(map[int]string)
	}

	lastRecvRoutine[routine][id] = vc[routine].Copy()
	lastRecvRoutinePos[routine][id] = pos

	for r, elem := range lastRecvRoutine {
		if r == routine {
			continue
		}

		if elem[id].clock == nil {
			continue
		}

		happensBefore := GetHappensBefore(elem[id], vc[routine])
		if happensBefore == Concurrent {
			found := "Found concurrent Recv on same channel:\n"
			found += "\trecv: " + pos + "\n"
			found += "\trecv : " + lastRecvRoutinePos[r][id]
			logging.Result(found, logging.CRITICAL)
		}
	}
}
