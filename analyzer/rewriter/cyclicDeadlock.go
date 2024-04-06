package rewriter

import (
	"analyzer/bugs"
	"analyzer/trace"
	"analyzer/utils"
	"errors"
	"fmt"
	"sort"
)

func rewriteCyclicDeadlock(bug bugs.Bug) error {
	firstTime := -1
	lastTime := -1

	if len(bug.TraceElement2) == 0 {
		return errors.New("No trace elements in bug")
	}

	for _, elem := range bug.TraceElement2 {
		// get the first and last mutex operation in the cycle
		time := (*elem).GetTPre()
		if firstTime == -1 || time < firstTime {
			firstTime = time
		}
		if lastTime == -1 || time > lastTime {
			lastTime = time
		}
	}

	// remove tail after lastTime
	trace.ShortenTrace(lastTime, true)

	fmt.Print("\n")
	PrintTrace([]string{"M"})
	fmt.Print("\n")

	routinesInCycle := make(map[int]struct{})

	maxIterations := 100 // prevent infinite loop
	for iter := 0; iter < maxIterations; iter++ {
		found := false
		// for all edges in the cycle shift the routine so that the next element is before the current element
		for i := 0; i < len(bug.TraceElement2); i++ {
			routinesInCycle[(*bug.TraceElement2[i]).GetRoutine()] = struct{}{}

			j := (i + 1) % len(bug.TraceElement2)

			elem1 := bug.TraceElement2[i]
			elem2 := bug.TraceElement2[j]

			if (*elem1).GetRoutine() == (*elem2).GetRoutine() {
				continue
			}

			// shift the routine of elem1 so that elem 2 is before elem1
			res := trace.ShiftRoutine((*elem1).GetRoutine(), (*elem1).GetTPre(), (*elem2).GetTPre()-(*elem1).GetTPre()+1)

			if res {
				found = true
			}
		}

		if !found {
			break
		}
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

	fmt.Println(firstTime, lastTime+1)

	fmt.Print("\n")
	PrintTrace([]string{"M"})
	fmt.Print("\n")

	return nil
}

/*
* Print the trace sorted by tPre
* Args:
*   types: types of the elements to print. If empty, all elements will be printed
* TODO: remove
 */
func PrintTrace(types []string) {
	elements := make([]struct {
		string
		int
	}, 0)
	for _, tra := range *trace.GetTraces() {
		for _, elem := range tra {
			elemStr := elem.ToString()
			if len(types) == 0 || utils.Contains(types, elemStr[0:1]) {
				elements = append(elements, struct {
					string
					int
				}{elemStr, elem.GetTPre()})
			}
		}
	}

	// sort elements by timestamp
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].int < elements[j].int
	})

	for _, elem := range elements {
		fmt.Println(elem.string)
	}
}
