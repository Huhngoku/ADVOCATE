package analysis

import (
	"analyzer/clock"
	"strconv"
)

/*
 * Struct to represent a node in a lock graph
 */
type lockGraphNode struct {
	id       int               // id of the mutex represented by the node
	routine  int               // id of the routine that holds the lock
	rw       bool              // true if the mutex is a read-write lock
	rLock    bool              // true if the lock was a read lock
	children []*lockGraphNode  // children of the node
	outside  []*lockGraphNode  // nodes with the same lock ID that are in the tree of another routine
	lockSet  []int             // ids of the nodes that are hold by the routine, when the node was created
	vc       clock.VectorClock // vector clock of the node, is equal to the vector clock of the lock event
	parent   *lockGraphNode    // parent of the node
	tID      string            // trace id of the lock
	visited  map[int]struct{}  // map to store the routine, for which the node was already visited when starting the DFS from the routines lock tree root
}

/*
 * Create a new lock graph
 * Returns:
 *   (*lockGraphNode): The new node
 */
func newLockGraph(routine int) *lockGraphNode {
	return &lockGraphNode{id: -1, routine: routine}
}

/*
 * Add a child to the node
 * Args:
 *   childID (int): The id of the child
 *   childRw (bool): True if the child is a read-write lock
 *   childRLock (bool): True if the child is a read lock
 *   vc (VectorClock): The vector clock of the childs lock operation
 *   lockSet ([]int): The lockSet of the child
 */
func (node *lockGraphNode) addChild(childID int, tID string, childRw bool, childRLock bool, vc clock.VectorClock, lockSet []int) *lockGraphNode {
	child := &lockGraphNode{id: childID, parent: node, rw: childRw,
		rLock: childRLock, routine: node.routine, vc: vc, lockSet: lockSet, tID: tID}
	node.children = append(node.children, child)
	return node.children[len(node.children)-1]
}

func (node *lockGraphNode) print() {
	result := node.toString()
	println(result)
}

func (node *lockGraphNode) toString() string {
	if node == nil {
		return ""
	}
	result := ""

	for _, child := range node.children {
		result += child.toStringTraverse(1)
	}

	return result
}

func (node *lockGraphNode) toStringTraverse(depth int) string {
	if node == nil {
		return ""
	}

	result := ""
	for i := 0; i < depth-1; i++ {
		result += "  "
	}
	result += strconv.Itoa(node.id) + "\n"

	for _, child := range node.children {
		result += child.toStringTraverse(depth + 1)
	}
	return result
}

func printTrees() {
	for routine, node := range lockGraphs {
		println("Routine " + strconv.Itoa(routine))
		node.print()
	}

}

// currend node for each routine
var currentNode = make(map[int][]*lockGraphNode) // routine -> []*lockGraphNode
// lock graph for each routine
var lockGraphs = make(map[int]*lockGraphNode) // routine -> lockGraphNode
// all nodes for each id
var nodesPerID = make(map[int]map[int][]*lockGraphNode) // id -> routine -> []*lockGraphNode

/*
 * Add the lock to the currently hold locks
 * Add the node to the lock tree
 * Args:
 *   id (int): The id of the lock
 *   routine (int): The id of the routine
 *   rw (bool): True if the lock is a read-write lock
 *   rLock (bool): True if the lock is a read lock
 *   vc (VectorClock): The vector clock of the lock event
 *   tPre (int): The timestamp at the end of the event
 */
func AnalysisCyclickDeadlockMutexLock(id int, tID string, routine int, rw bool, rLock bool, vc clock.VectorClock, tPost int) {
	if tPost == 0 {
		return
	}

	// create new lock tree if it does not exist yet
	if _, ok := lockGraphs[routine]; !ok {
		lockGraphs[routine] = newLockGraph(routine)
		currentNode[routine] = []*lockGraphNode{lockGraphs[routine]}
	}

	// create empty map for nodesPerID if it does not exist yet
	if _, ok := nodesPerID[id]; !ok {
		nodesPerID[id] = make(map[int][]*lockGraphNode)
	}
	if _, ok := nodesPerID[id][routine]; !ok {
		nodesPerID[id][routine] = []*lockGraphNode{}
	}

	// add the lock element to the lock tree
	// update the current lock
	node := currentNode[routine][len(currentNode[routine])-1].addChild(id, tID, rw, rLock, vc.Copy(), getCurrentLockSet(routine))
	currentNode[routine] = append(currentNode[routine], node)
	nodesPerID[id][routine] = append(nodesPerID[id][routine], node)
}

/*
 * Remove the lock from the currently hold locks
 * Args:
 *   id (int): The id of the lock
 *   routine (int): The id of the routine
 *   tPost (int): The timestamp at the end of the event
 */
func AnalysisCyclicDeadlockMutexUnLock(id int, routine int, tPost int) {
	if tPost == 0 {
		return
	}

	for i := len(currentNode[routine]) - 1; i >= 0; i-- {
		if currentNode[routine][i].id == id {
			currentNode[routine] = currentNode[routine][:i]
			return
		}
	}
}

/*
 * Check if the lock graph created by connecting all lock trees is cyclic
 * If there are cycles, log the results
 */
func CheckForCyclicDeadlock() {
	findOutsideConnections()
	found, cycles := findCycles() // find all cycles in the lock graph

	if !found { // no cycles
		return
	}

	// remove duplicate cycles
	cycles = removeCyclicPermutations(cycles)

	for _, cycle := range cycles {
		// check if the cycle can create a deadlock
		res := isCycleDeadlock(cycle)
		if res {
			// TODO: add a result. Cyclic deadlock detection is currently disabled
			// found := "Possible cyclic deadlock:\n"
			// found += "\thead: " + cycle[0].tID + "\n"
			// found += "\ttail: "
			// for i := 0; i < len(cycle); i++ {
			// 	found += cycle[i].tID + ";"
			// }
			// logging.Result(found, logging.CRITICAL)
		}
	}
}

/*
 * Find all connections between lock trees for different routines
 * A connection exists iff both nodes have the same id but different routines
 */
func findOutsideConnections() {
	for _, tree := range lockGraphs { // for each lock tree
		traverseTreeAndAddOutsideConnections(tree)
	}
}

/*
 * Traverse all nodes of the tree recursively.
 * For each node, add all nodes with the same id but different routine to the outside connections
 * Args:
 *   node (*lockGraphNode): The node to start the traversal
 */
func traverseTreeAndAddOutsideConnections(node *lockGraphNode) {
	if node == nil {
		return
	}

	for routine, outsideNodes := range nodesPerID[node.id] {
		if routine == node.routine {
			continue
		}

		for _, outsideNode := range outsideNodes {
			node.outside = append(node.outside, outsideNode)
		}
	}

	for _, child := range node.children {
		traverseTreeAndAddOutsideConnections(child)
	}
}

/*
 * Find all cycles in the lock graph formed by connecting all lock trees
 * using the outside connections.
 * Return all the cycles as a list of nodes
 * Returns:
 *  (bool): True if there are cycles
 *  ([][]*lockGraphNode): A list of cycles, where each cycle is a list of nodes
 */
func findCycles() (bool, [][]*lockGraphNode) {
	cycles := [][]*lockGraphNode{}
	for routine, tree := range lockGraphs { // for each lock tree
		findCyclesDFS(tree, &([]*lockGraphNode{}), &cycles, routine, nil)
	}

	if len(cycles) == 0 {
		return false, nil
	}
	return true, cycles
}

func findCyclesDFS(node *lockGraphNode, currentPath *([]*lockGraphNode),
	cycles *([][]*lockGraphNode), routine int, last *lockGraphNode) {
	if node == nil {
		return
	}

	// make node.visited if it does not exist yet
	if node.visited == nil {
		node.visited = make(map[int]struct{})
	}

	if _, ok := node.visited[routine]; ok { // node was already visited
		cycle, index := isInCurrentPath(node, currentPath)
		if cycle {
			copySlice := make([]*lockGraphNode, len(*currentPath)-index)
			copy(copySlice, (*currentPath)[index:])
			*cycles = append(*cycles, copySlice)
		}
		return
	}

	if node.id != -1 { // not for root
		node.visited[routine] = struct{}{}
		*currentPath = append(*currentPath, node)
	}

	// recursion step for each child
	for _, child := range node.children {
		findCyclesDFS(child, currentPath, cycles, routine, nil)
	}

	// recursion step for each outside connection
	for _, outside := range node.outside {
		if outside == last {
			continue
		}
		findCyclesDFS(outside, currentPath, cycles, routine, node)
	}

	// remove node from current path
	if node.id != -1 {
		*currentPath = (*currentPath)[:len(*currentPath)-1]
	}

}

func isInCurrentPath(node *lockGraphNode, currentPath *([]*lockGraphNode)) (bool, int) {
	for i, pathNode := range *currentPath {
		if pathNode == node {
			return true, i
		}
	}
	return false, -1
}

/*
 * Remove cyclic permutations
 * Args:
 *   cycles ([][]*lockGraphNode): The cycles to remove permutations from
 * Returns:
 *   ([][]*lockGraphNode): The cycles without cyclic permutations
 */
func removeCyclicPermutations(cycles [][]*lockGraphNode) [][]*lockGraphNode {
	// remove cyclic permutations (same cycle but different starting point)
	for i := 0; i < len(cycles); i++ {
		for j := i + 1; j < len(cycles); j++ {
			if len(cycles[i]) == len(cycles[j]) {
				if isCyclicPermutation(cycles[i], cycles[j]) {
					cycles = append(cycles[:j], cycles[j+1:]...)
					j--
				}
			}
		}
	}
	return cycles
}

/*
 * Check if two cycles are cyclic permutations of each other. The function
 * assumes that the cycles have the same length.
 * Args:
 *   cycle1 ([]*lockGraphNode): The first cycle
 *   cycle2 ([]*lockGraphNode): The second cycle
 * Returns:
 *   (bool): True if the cycles are cyclic permutations of each other
 */
func isCyclicPermutation(cycle1 []*lockGraphNode, cycle2 []*lockGraphNode) bool {
	for i := 0; i < len(cycle1); i++ {
		if cycle1[0] == cycle2[i] {
			for j := 0; j < len(cycle1); j++ {
				if cycle1[j] != cycle2[(i+j)%len(cycle1)] {
					return false
				}
			}
			return true
		}
	}
	return false
}

/*
 * Check if a cycle can create a deadlock
 * It can not be a deadlock, if at least on of the following is false:
 * - the cycle consists of more than one different lock (R1)
 * - the lock operations in the cycle for different routines are concurrent (R2)
 * - two operations on the same lock connected by an edge are not both read operations (R3)
 * - the cycle is valid considering gate locks (R4)
 * Args:
 *   cycle ([]*lockGraphNode): The cycle to check
 * Returns:
 *   (bool): True if the cycle can create a deadlock
 */
func isCycleDeadlock(cycle []*lockGraphNode) bool {
	// does the cycle consists of more than one different lock? (R1)
	if !isCycleMoreThanOneMutex(cycle) {
		return false
	}

	// are the lock operation int the cycle for different routines concurrent? (R2)
	if !isCycleConcurrent(cycle) {
		return false
	}

	// check, that the cycle is valid considering read-write locks (R3)
	if !isCycleValidRead(cycle) {
		return false
	}

	if !isCycleValidGate(cycle) {
		return false
	}

	//

	return true
}

/*
 * Check if the cycle consists of more than one different lock
 * Args:
 *   cycle ([]*lockGraphNode): The cycle to check
 * Returns:
 *   (bool): True if the cycle consists of more than one different lock
 */
func isCycleMoreThanOneMutex(cycle []*lockGraphNode) bool {
	moreThanOneMutexIndex := -1
	moreThanOneMutexBool := false

	for _, node := range cycle {
		if moreThanOneMutexIndex == -1 {
			moreThanOneMutexIndex = node.id
		} else if moreThanOneMutexIndex != node.id {
			moreThanOneMutexBool = true
		}
	}

	return moreThanOneMutexBool
}

/*
 * Check if all lock operations in the cycle are concurrent
 * Args:
 *   cycle ([]*lockGraphNode): The cycle to check
 * Returns:
 *   (bool): True if all lock operations in the cycle are concurrent
 */
func isCycleConcurrent(cycle []*lockGraphNode) bool {
	for i := 0; i < len(cycle); i++ {
		for j := i + 1; j < len(cycle); j++ {
			if cycle[i].routine == cycle[j].routine {
				continue
			}

			happensBefore := clock.GetHappensBefore(cycle[i].vc, cycle[j].vc)
			if happensBefore != clock.Concurrent {
				return false
			}
		}
	}
	return true
}

/*
 * Check, that the cycle is valid considering read-write locks
 * Two operations on the same lock connected by an edge are not both read operations
 * Args:
 *   cycle ([]*lockGraphNode): The cycle to check
 * Returns:
 *   (bool): True if the cycle is valid considering read-write locks
 */
func isCycleValidRead(cycle []*lockGraphNode) bool {
	for i := 0; i < len(cycle); i++ {
		for j := i + 1; j < len(cycle); j++ {
			for _, ls1 := range cycle[i].lockSet {
				for _, ls2 := range cycle[j].lockSet {
					if ls1 != ls2 {
						continue
					}

					if cycle[i].rw && cycle[j].rw {
						return false
					}
				}
			}
		}
	}
	return true
}

/*
 * Check, that the cycle is valid considering gate locks
 * Args:
 *   cycle ([]*lockGraphNode): The cycle to check
 * Returns:
 *   (bool): True if the cycle is valid considering gate locks
 */
func isCycleValidGate(cycle []*lockGraphNode) bool {
	for i := 0; i < len(cycle); i++ {
		for j := i + 1; j < len(cycle); j++ {
			if cycle[i].id != cycle[j].id {
				continue
			}

			if cycle[i].rw && cycle[j].rw {
				return false
			}
		}
	}
	return true
}

/*
 * Get the current lock set of a routine
 * Args:
 *   routine (int): The id of the routine
 * Returns:
 *   ([]int): The current lock set of the routine
 */
func getCurrentLockSet(routine int) []int {
	ls := make([]int, len(currentNode[routine]))
	for id, _ := range lockSet[routine] {
		ls = append(ls, id)
	}
	return ls
}
