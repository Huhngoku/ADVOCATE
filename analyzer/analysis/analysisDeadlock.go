package analysis

import (
	"strconv"
)

/*
 * Struct to represent a node in a lock graph
 */
type lockGraphNode struct {
	id       int              // id of the mutex represented by the node
	routine  int              // id of the routine that holds the lock
	rw       bool             // true if the mutex is a read-write lock
	rLock    bool             // true if the lock was a read lock
	children []*lockGraphNode // children of the node
	outside  []*lockGraphNode // nodes with the same lock ID that are in the tree of another routine
	parent   *lockGraphNode   // parent of the node
	visited  bool             // true if the node was already visited by the deadlock detection algorithm  // TODO: can we find the actual path with that or only that there is a circle
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
 */
func (node *lockGraphNode) addChild(childID int, childRw bool, childRLock bool) *lockGraphNode {
	child := &lockGraphNode{id: childID, parent: node, rw: childRw, rLock: childRLock, routine: node.routine}
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
 */
func AnalysisDeadlockMutexLock(id int, routine int, rw bool, rLock bool) {
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
	node := currentNode[routine][len(currentNode[routine])-1].addChild(id, rw, rLock)
	currentNode[routine] = append(currentNode[routine], node)
	nodesPerID[id][routine] = append(nodesPerID[id][routine], node)
}

/*
 * Remove the lock from the currently hold locks
 * Args:
 *   id (int): The id of the lock
 *   routine (int): The id of the routine
 */
func AnalysisDeadlockMutexUnLock(id int, routine int) {
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
	printTrees()

	findOutsideConnections()
	cycles := findCycles()

	if len(cycles) == 0 {
		return
	}

	// TODO: remove duplicate cycles

	for _, cycle := range cycles {
		res := isCycleDeadlock(cycle)
		if res {
			// TODO: log the cycle
			println("Deadlock detected")
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

	for routine, outsideNode := range nodesPerID[node.id] {
		if routine == node.routine {
			continue
		}

		for _, node := range outsideNode {
			node.outside = append(node.outside, node)
		}
	}

	for _, child := range node.children {
		traverseTreeAndAddOutsideConnections(child)
	}
}

/*
 * Find all unique cycles in the lock graph formed by connecting all lock trees
 * using the outside connections.
 * Return all the cycles as a list of nodes
 * Returns:
 *   ([][]*lockGraphNode): A list of cycles, where each cycle is a list of nodes
 */
func findCycles() [][]*lockGraphNode {
	// TODO: implement using a DFS search
	panic("Not implemented yet")
}

/*
 * Check if a cycle can create a deadlock
 * Args:
 *   cycle ([]*lockGraphNode): The cycle to check
 * Returns:
 *   (bool): True if the cycle can create a deadlock
 */
func isCycleDeadlock(cycle []*lockGraphNode) bool {
	// TODO: implement
	panic("Not implemented yet")
}
