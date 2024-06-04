package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
)

func checkForConcurrentRecv(routine int, id int, tID string, vc map[int]clock.VectorClock, tPost int) {
	for r, elem := range lastRecvRoutine {
		println(r, elem[id].Vc.ToString())
		if r == routine {
			continue
		}

		if elem[id].Vc.GetClock() == nil {
			continue
		}

		happensBefore := clock.GetHappensBefore(elem[id].Vc, vc[routine])
		if happensBefore == clock.Concurrent {

			file1, line1, tPre1, err := infoFromTID(tID)
			if err != nil {
				logging.Debug(err.Error(), logging.ERROR)
				return
			}

			file2, line2, tPre2, err := infoFromTID(lastRecvRoutine[r][id].TID)

			arg1 := logging.TraceElementResult{
				RoutineID: routine,
				ObjID:     id,
				TPre:      tPre1,
				ObjType:   "CR",
				File:      file1,
				Line:      line1,
			}

			arg2 := logging.TraceElementResult{
				RoutineID: r,
				ObjID:     id,
				TPre:      tPre2,
				ObjType:   "CR",
				File:      file2,
				Line:      line2,
			}

			logging.Result(logging.WARNING, logging.AConcurrentRecv,
				"recv", []logging.ResultElem{arg1}, "recv", []logging.ResultElem{arg2})
		}
	}

	if tPost != 0 {
		if _, ok := lastRecvRoutine[routine]; !ok {
			lastRecvRoutine[routine] = make(map[int]VectorClockTID)
		}

		lastRecvRoutine[routine][id] = VectorClockTID{vc[routine].Copy(), tID, routine}
	}
}
