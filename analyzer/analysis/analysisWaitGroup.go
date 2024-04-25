package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
	"analyzer/utils"
	"errors"
	"strconv"
)

func checkForDoneBeforeAddChange(routine int, id int, delta int, pos string, vc clock.VectorClock) {
	if delta > 0 {
		checkForDoneBeforeAddAdd(routine, id, pos, vc, delta)
	} else if delta < 0 {
		checkForDoneBeforeAddDone(routine, id, pos, vc)
	} else {
		// checkForImpossibleWait(routine, id, pos, vc)
	}
}

func checkForDoneBeforeAddAdd(routine int, id int, pos string, vc clock.VectorClock, delta int) {
	// if necessary, create maps and lists
	if _, ok := wgAdd[id]; !ok {
		wgAdd[id] = make(map[int][]VectorClockTID)
	}
	if _, ok := wgAdd[id][routine]; !ok {
		wgAdd[id][routine] = make([]VectorClockTID, 0)
	}

	// add the vector clock and position to the list
	for i := 0; i < delta; i++ {
		if delta > 1 {
			pos = pos + "+" + strconv.Itoa(i) // add a unique identifier to the position
		}
		wgAdd[id][routine] = append(wgAdd[id][routine], VectorClockTID{vc.Copy(), pos})
	}
}

func checkForDoneBeforeAddDone(routine int, id int, pos string, vc clock.VectorClock) {
	// if necessary, create maps and lists
	if _, ok := wgDone[id]; !ok {
		wgDone[id] = make(map[int][]VectorClockTID)

	}
	if _, ok := wgDone[id][routine]; !ok {
		wgDone[id][routine] = make([]VectorClockTID, 0)
	}

	// add the vector clock and position to the list
	wgDone[id][routine] = append(wgDone[id][routine], VectorClockTID{vc.Copy(), pos})
}

/*
 * Build a st graph for a wait group.
 * The graph has the following structure:
 * - a start node s
 * - a end node t
 * - edges from s to all done operations
 * - edges from all add operations to t
 * - edges from done to add if the add happens before the done
 * Args:
 *   adds (map[int][]VectorClockTID): The add operations
 *   dones (map[int][]VectorClockTID): The done operations
 * Returns:
 *   []Edge: The graph
 */
func buildResidualGraph(adds map[int][]VectorClockTID, dones map[int][]VectorClockTID) map[string][]string {
	graph := make(map[string][]string, 0)
	graph["s"] = []string{}
	graph["t"] = []string{}

	// add edges from s to all done operations
	for _, done := range dones {
		for _, vc := range done {
			graph[vc.tID] = []string{}
			graph["s"] = append(graph["s"], vc.tID)
		}
	}

	// add edges from all add operations to t
	for _, add := range adds {
		for _, vc := range add {
			graph[vc.tID] = []string{"t"}
		}
	}

	// add edge from done to add if the add happens before the done
	for _, done := range dones {
		for _, vcDone := range done {
			for _, add := range adds {
				for _, vcAdd := range add {
					if clock.GetHappensBefore(vcAdd.vc, vcDone.vc) == clock.Before {
						graph[vcDone.tID] = append(graph[vcDone.tID], vcAdd.tID)

					}
				}
			}
		}
	}

	return graph
}

/*
 * Calculate the maximum flow of a graph using the ford fulkerson algorithm
 * Args:
 *   graph ([]Edge): The graph
 * Returns:
 *   int: The maximum flow
 */
func calculateMaxFlow(graph map[string][]string) (int, map[string][]string) {
	maxFlow := 0
	for {
		path, flow := findPath(graph)
		if flow == 0 {
			break
		}

		maxFlow += flow
		for i := 0; i < len(path)-1; i++ {
			graph[path[i]] = append(graph[path[i]], path[i+1])
			graph[path[i+1]] = remove(graph[path[i+1]], path[i])
		}
	}

	return maxFlow, graph
}

/*
 * Find a path in a graph using a breadth-first search
 * Args:
 *   graph ([]Edge): The graph
 * Returns:
 *   []string: The path
 *   int: The flow
 */
func findPath(graph map[string][]string) ([]string, int) {
	visited := make(map[string]bool, 0)
	queue := []string{"s"}
	visited["s"] = true
	parents := make(map[string]string, 0)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node == "t" {
			path := []string{}
			for node != "s" {
				path = append(path, node)
				node = parents[node]
			}
			path = append(path, "s")

			return path, 1
		}

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
				parents[neighbor] = node
			}
		}
	}

	return []string{}, 0
}

/*
 * Remove an element from a list
 * Args:
 *   list ([]string): The list
 *   element (string): The element to remove
 * Returns:
 *   []string: The list without the element
 */
func remove(list []string, element string) []string {
	for i, e := range list {
		if e == element {
			list = append(list[:i], list[i+1:]...)
			return list
		}
	}
	return list
}

func numberDone(id int) int {
	res := 0
	for _, dones := range wgDone[id] {
		res += len(dones)
	}
	return res
}

/*
- Check if a wait group counter could become negative
- For each done operation, build a bipartite st graph.
- Use the Ford-Fulkerson algorithm to find the maximum flow.
- If the maximum flow is smaller than the number of done operations, a negative wait group counter is possible.
*/
func CheckForDoneBeforeAdd() {
	for id := range wgAdd { // for all waitgroups
		graph := buildResidualGraph(wgAdd[id], wgDone[id])

		maxFlow, graph := calculateMaxFlow(graph)
		nrDone := numberDone(id)

		addsVcTIDs := []VectorClockTID{}
		donesVcTIDs := []VectorClockTID{}

		if maxFlow < nrDone {
			// sort the adds and dones, that do not have a partner is such a way,
			// that the i-th add in the result message is concurrent with the
			// i-th done in the result message

			for _, adds := range wgAdd[id] {
				for _, add := range adds {
					if !utils.Contains(graph["t"], add.tID) {
						addsVcTIDs = append(addsVcTIDs, add)
					}
				}
			}
			for _, dones := range graph["s"] {
				doneVcTID, err := getDoneVcTIDFromTID(id, dones)
				if err != nil {
					logging.Debug(err.Error(), logging.ERROR)
				} else {
					donesVcTIDs = append(donesVcTIDs, doneVcTID)
				}
			}

			addsVcTIDSorted := make([]VectorClockTID, 0)
			donesVcTIDSorted := make([]VectorClockTID, 0)

			for i := 0; i < len(addsVcTIDs); i++ {
				for j := 0; j < len(donesVcTIDs); j++ {
					if clock.GetHappensBefore(addsVcTIDs[i].vc, addsVcTIDs[j].vc) == clock.Concurrent {
						addsVcTIDSorted = append(addsVcTIDSorted, addsVcTIDs[i])
						donesVcTIDSorted = append(donesVcTIDSorted, donesVcTIDs[j])
						// remove the element from the list
						addsVcTIDs = append(addsVcTIDs[:i], addsVcTIDs[i+1:]...)
						donesVcTIDs = append(donesVcTIDs[:j], donesVcTIDs[j+1:]...)
						// fix the index
						i--
						j = 0
					}
				}
			}

			message := "Possible negative waitgroup counter:\n"
			message += "\tadd: "
			for _, add := range addsVcTIDSorted {
				if add.tID == "\n" {
					message += add.tID + " (unknown); "
				}
				message += add.tID + "; "
			}

			message += "\n\tdone: "
			for _, done := range donesVcTIDSorted {
				message += done.tID + "; "
			}

			logging.Result(message, logging.CRITICAL)
		}
	}
}

func getDoneVcTIDFromTID(id int, tID string) (VectorClockTID, error) {
	for _, dones := range wgDone[id] {
		for _, done := range dones {
			if done.tID == tID {
				return done, nil
			}
		}
	}
	return VectorClockTID{}, errors.New("Could not find done operation with tID " + tID)
}
