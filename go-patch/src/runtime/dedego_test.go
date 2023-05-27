// DEDEGO-FILE-START

package runtime_test

import (
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDedegoMutex(t *testing.T) {
	var m sync.Mutex
	var n sync.Mutex

	var trace2 string

	m.Lock()
	m.Unlock()
	m.TryLock()
	m.Unlock()

	trace1 := runtime.TraceToString()

	go func() {
		m.Lock()
		m.TryLock()
		m.Unlock()
		n.Lock()
		n.Unlock()
		trace2 = runtime.TraceToString()
	}()

	time.Sleep(100 * time.Millisecond)

	traceTotal := trace1 + "|" + trace2

	// check that form is correct
	exp := regexp.MustCompile("(?i)M,[0-9]+,1,2,-,L;M,[0-9]+,3,3,-,U;" +
		"M,[0-9]+,4,5,-,T,succ;M,[0-9]+,6,6,-,U\\|M,[0-9]+,1,2,-,L;" +
		"M,[0-9]+,3,4,-,T,fail;M,[0-9]+,5,5,-,U;M,[0-9]+,6,7,-,L;" +
		"M,[0-9]+,8,8,-,U")

	if !exp.MatchString(traceTotal) {
		t.Errorf("Trace in TestDedegoMutex is not correct: %s", traceTotal)
	}

	// check that the ids of the same mutex are the same
	traceTotal = trace1 + ";" + trace2
	traces := strings.Split(traceTotal, ";")
	elems := make([]string, 0)
	for _, elem := range traces {
		elems = append(elems, strings.Split(elem, ",")[1])
	}

	equal := map[int][]int{
		0: []int{1, 2, 3, 4, 5, 6},
		7: []int{8},
	}

	different := map[int][]int{
		0: []int{7},
	}

	for i, ids := range equal {
		for _, id := range ids {
			if elems[i] != elems[id] {
				t.Errorf("Mutex ids should be equal: (%s, %s) != (%s, %s)", strconv.Itoa(i), elems[i], strconv.Itoa(id), elems[id])
			}
		}
	}

	for i, ids := range different {
		for _, id := range ids {
			if elems[i] == elems[id] {
				t.Errorf("Mutex ids should not be equal: (%s, %s) == (%s, %s)", strconv.Itoa(i), elems[i], strconv.Itoa(id), elems[id])
			}
		}
	}
}

// func TestDedegoRWMutex(t *testing.T) {
// 	var m sync.RWMutex
// 	var n sync.RWMutex

// 	m.RLock()
// 	n.Lock()
// 	n.TryRLock()
// 	n.Unlock()
// 	m.RUnlock()

// 	trace1 := runtime.TraceToString()
// 	var trace2 string

// 	go func() {
// 		m.RLock()
// 		m.TryRLock()
// 		m.RUnlock()
// 		m.RUnlock()
// 		n.RLock()
// 		n.RUnlock()
// 		n.TryLock()
// 		n.Unlock()

// 		trace2 = runtime.TraceToString()
// 	}()

// 	time.Sleep(100 * time.Millisecond)

// 	traceTotal := trace1 + "|" + trace2

// 	exp := regexp.MustCompile("(?i)M,[0-9]+,1,2,RW,LR;M,[0-9]+,3,6,RW,L;" +
// 		"M,[0-9]+,4,5,-,L;M,[0-9]+,7,8,RW,TR,fail;M,[0-9]+,9,9,-,U;" +
// 		"M,[0-9]+,10,10,RW,U;M,[0-9]+,11,11,RW,UR\\|M,[0-9]+,1,2,RW,LR;" +
// 		"M,[0-9]+,3,4,RW,TR,succ;M,[0-9]+,5,5,RW,UR;M,[0-9]+,6,6,RW,UR;" +
// 		"M,[0-9]+,7,8,RW,LR;M,[0-9]+,9,9,RW,UR;M,[0-9]+,10,13,RW,T,succ;" +
// 		"M,[0-9]+,11,12,-,T,succ;M,[0-9]+,14,14,-,U;M,[0-9]+,15,15,RW,U")

// 	// check that form is correct
// 	if !exp.MatchString(traceTotal) {
// 		t.Errorf("Trace in TestDedegoRWMutex is not correct: %s", traceTotal)
// 	}

// 	traceTotal = trace1 + ";" + trace2
// 	traces := strings.Split(traceTotal, ";")
// 	elems := make([]string, 0)
// 	for _, elem := range traces {
// 		elems = append(elems, strings.Split(elem, ",")[1])
// 	}

// 	equal := map[int][]int{
// 		0: []int{4, 5, 6, 7, 8},
// 		1: []int{2, 3, 9, 10, 11, 12},
// 	}

// 	different := map[int][]int{
// 		0: []int{1},
// 	}

// 	for i, ids := range equal {
// 		for _, id := range ids {
// 			if elems[i] != elems[id] {
// 				t.Errorf("RWMutex ids should be equal: (%s, %s) != (%s, %s)", strconv.Itoa(i), elems[i], strconv.Itoa(id), elems[id])
// 			}
// 		}
// 	}

// 	for i, ids := range different {
// 		for _, id := range ids {
// 			if elems[i] == elems[id] {
// 				t.Errorf("RWMutex ids should not be equal: (%s, %s) == (%s, %s)", strconv.Itoa(i), elems[i], strconv.Itoa(id), elems[id])
// 			}
// 		}
// 	}
// }

/*
Test for the recording of wait group operations in the trace
*/
func TestDedegoWaitGroup(t *testing.T) {
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup

	traces := make([]string, 4)

	wg1.Add(3)
	wg2.Add(3)
	for i := 0; i < 3; i++ {
		go func(i int) {
			time.Sleep(500 * time.Millisecond)
			wg1.Done()
			wg2.Done()
			traces[i+1] = runtime.TraceToString()
		}(i)
	}

	wg1.Wait()
	wg2.Wait()

	traces[0] = runtime.TraceToString()

	// check that form is correct
	tracesTotal := strings.Join(traces, "|")
	exp := regexp.MustCompile("(?i)W,[0-9]+,1,0,W,3,3;W,[0-9]+,2,0,W,3,3;W,[0-9]+,3,4,W,0,0;W,[0-9]+,5,6,W,0,0\\|W,[0-9]+,1,0,W,-1,0;W,[0-9]+,2,0,W,-1,0\\|W,[0-9]+,1,0,W,-1,2;W,[0-9]+,2,0,W,-1,2\\|W,[0-9]+,1,0,W,-1,1;W,[0-9]+,2,0,W,-1,1")
	if !exp.MatchString(tracesTotal) {
		t.Errorf("Trace in TestDedegoWaitGroup is not correct: %s", tracesTotal)
	}

	// check that ids are correct
	tracesTotal = strings.Join(traces, ";")
	traces = strings.Split(tracesTotal, ";")
	elems := make([]string, 0)
	for _, elem := range traces {
		elems = append(elems, strings.Split(elem, ",")[1])
	}
	equal := map[int][]int{
		0: []int{2, 4, 6, 8},
		1: []int{3, 5, 7, 9},
	}

	different := map[int][]int{
		0: []int{1},
	}

	for i, ids := range equal {
		for _, id := range ids {
			if elems[i] != elems[id] {
				t.Errorf("WaitGroup ids should be equal: (%s, %s) != (%s, %s)", strconv.Itoa(i), elems[i], strconv.Itoa(id), elems[id])
			}
		}
	}

	for i, ids := range different {
		for _, id := range ids {
			if elems[i] == elems[id] {
				t.Errorf("WaitGroup ids should not be equal: (%s, %s) == (%s, %s)", strconv.Itoa(i), elems[i], strconv.Itoa(id), elems[id])
			}
		}
	}

}
