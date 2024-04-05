package rewriter

import (
	"analyzer/bugs"
	"analyzer/trace"
	"errors"
)

func rewriteCyclicDeadlock(bug bugs.Bug) error {
	firstTime := -1
	lastTime := -1

	if len(bug.TraceElement2) == 0 {
		return errors.New("No trace elements in bug")
	}

	for _, elem := range bug.TraceElement2 {
		// get the first and last mutex operation in the cycle
		time := (*elem).GetTSort()
		if firstTime == -1 || time < firstTime {
			firstTime = time
		}
		if lastTime == -1 || time > lastTime {
			lastTime = time
		}
	}

	// remove tail after lastTime
	trace.ShortenTrace(lastTime, true)

	// TODO: the locks must be ordered based on the lock tree
	routinesInCycle := make(map[int]struct{})

	for i := 0; i < len(bug.TraceElement2); i++ {
		routinesInCycle[(*bug.TraceElement2[i]).GetRoutine()] = struct{}{}

		j := (i + 1) % len(bug.TraceElement2)

		elem1 := bug.TraceElement2[i]
		elem2 := bug.TraceElement2[j]

		if (*elem1).GetRoutine() == (*elem2).GetRoutine() {
			continue
		}

		// shift the routine of elem1 so that elem 2 is before elem1
		trace.ShiftRoutine((*elem1).GetRoutine(), (*elem1).GetTSort(), (*elem2).GetTSort()-(*elem1).GetTSort()+1)
	}

	currentTrace := trace.GetTraces()
	lastTime = -1

	for routine := range routinesInCycle {
		found := false
		for i := len((*currentTrace)[routine]) - 1; i >= 0; i-- {
			elem := (*currentTrace)[routine][i]
			switch elem := elem.(type) {
			case *trace.TraceElementMutex:
				if (*elem).IsLock() {
					trace.ShortenRoutineIndex(routine, i, true)
					if lastTime == -1 || (*elem).GetTSort() > lastTime {
						lastTime = (*elem).GetTSort()
					}
					found = true
				}
			}
			if found {
				break
			}
		}
	}

	// add start and end signals
	trace.AddTraceElementReplay(firstTime, true)
	trace.AddTraceElementReplay(lastTime+1, false)

	// for i, elem := range bug.TraceElement2 {
	// 	(*elem).SetTSort(firstTime + i)
	// }

	// // find the rewritten order and store in newOrder
	// // run until partialTrace is empty
	// for len(partialTrace) != 0 {
	// 	println("Partial trace length: ", len(partialTrace))
	// 	// check for each routine, if the next operation is a lock operation
	// 	// if so, add it to nextElems and remove it from partialTrace
	// 	// also save the routine where the first element has the earliest timestamp
	// 	// regardless of operation type
	// 	nextElems := make([]*trace.TraceElement, 0)
	// 	nextPossibleTime := -1
	// 	nextPossibleRoutine := -1
	// 	println("\n\n\n")
	// 	for routine, traceRout := range partialTrace {
	// 		if len(traceRout) == 0 {
	// 			continue
	// 		}

	// 		switch elem := (*traceRout[0]).(type) {
	// 		case *trace.TraceElementMutex:
	// 			op := (*elem).GetOperation()
	// 			time := (*elem).GetTSort()

	// 			if op == trace.LockOp || op == trace.RLockOp ||
	// 				op == trace.TryLockOp || op == trace.TryRLockOp {
	// 				nextElems = append(nextElems, traceRout[0])
	// 				partialTrace[routine] = partialTrace[routine][1:]
	// 			}

	// 			if nextPossibleTime == -1 || time < nextPossibleTime {
	// 				nextPossibleTime = time
	// 				nextPossibleRoutine = routine
	// 			}
	// 		}
	// 	}

	// 	if len(nextElems) != 0 {
	// 		// if there are elements in nextElems sort it by timestamp and add them in
	// 		// sorted order into newOrder
	// 		sort.Slice(nextElems, func(i, j int) bool {
	// 			return (*nextElems[i]).GetTSort() < (*nextElems[j]).GetTSort()
	// 		})

	// 		newOrder = append(newOrder, nextElems...)
	// 		nextElems = make([]*trace.TraceElement, 0)
	// 	} else {
	// 		if nextPossibleRoutine == -1 {
	// 			println("No possible routine found")
	// 			break
	// 		}
	// 		// if there are no elements possible lock operations, add the first
	// 		// element of the routine with the earliest timestamp to newOrder
	// 		// and remove it from partialTrace
	// 		newOrder = append(newOrder, partialTrace[nextPossibleRoutine][0])
	// 		partialTrace[nextPossibleRoutine] = partialTrace[nextPossibleRoutine][1:]
	// 	}
	// }

	// // add the start signal to the trace
	// trace.AddTraceElementReplay(firstTime, true)

	// // set the times on all elements in newOrder to the correct time
	// for i, elem := range newOrder {
	// 	(*elem).SetTSort(firstTime + i + 1)
	// 	println((*elem).ToString())
	// 	println("\n\n\n")
	// }

	// // add the end signal
	// trace.AddTraceElementReplay(firstTime+len(newOrder)+1, false)

	return nil
}
