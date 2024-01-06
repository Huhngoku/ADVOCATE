package analysis

import (
	"strconv"
)

/*
 * Struct to represent a node in a lock graph
 */
type lockGraphNode struct {
	id       int              // id of the mutex represented by the node
	rw       bool             // true if the mutex is a read-write lock
	rLock    bool             // true if the lock was a read lock
	children []*lockGraphNode // children of the node
	parent   *lockGraphNode   // parent of the node
	visited  bool             // true if the node was already visited by the deadlock detection algorithm  // TODO: can we find the actual path with that or only that there is a circle
}

/*
 * Create a new lock graph
 * Returns:
 *   (*lockGraphNode): The new node
 */
func newLockGraph() *lockGraphNode {
	return &lockGraphNode{id: -1}
}

/*
 * Add a child to the node
 * Args:
 *   childID (int): The id of the child
 *   childRw (bool): True if the child is a read-write lock
 *   childRLock (bool): True if the child is a read lock
 */
func (node *lockGraphNode) addChild(childID int, childRw bool, childRLock bool) *lockGraphNode {
	child := &lockGraphNode{id: childID, parent: node, rw: childRw, rLock: childRLock}
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
	result := "root\n"

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
	for i := 0; i < depth; i++ {
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

var currentNode = make(map[int][]*lockGraphNode) // routine -> []*lockGraphNode
var lockGraphs = make(map[int]*lockGraphNode)    // routine -> lockGraphNode

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
		lockGraphs[routine] = newLockGraph()
		currentNode[routine] = []*lockGraphNode{lockGraphs[routine]}
	}

	// add the lock element to the lock tree
	// update the current lock
	node := currentNode[routine][len(currentNode[routine])-1].addChild(id, rw, rLock)
	currentNode[routine] = append(currentNode[routine], node)
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
	panic("Not implemented yet")
}
