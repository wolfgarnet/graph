package graph

import (
	"fmt"
)

// Graph represents a graph
type Graph struct {
	Nodes []*Node
	NodeStringer func(interface{}) string

	OnEdgeCreated func(*Edge)
	// TODO a sub map of regions/graphs
}

// NewGraph returns a new graph
func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) Stringify() {
	for _, n := range g.Nodes {
		for _, e := range n.Edges {
			if e.Node1 != n {
				continue
			}

			fmt.Printf("[%v] %v -(%v)> [%v] %v\n", n, g.NodeStringer(n.Data), e.Data, e.Node2, g.NodeStringer(e.Node2.Data))
		}
	}
}

func (g *Graph) PrintNodes() {
	for _, node := range g.Nodes {
		if node.Data != nil {
			fmt.Printf("[%v]: %v(%T)\n", node, g.NodeStringer(node.Data), node.Data)
		} else {
			fmt.Printf("[%v]: NIL(%T)\n", node, node.Data)
		}
	}
}

// NewNode will add a new node to the graph.
// If unique is set, the ast node cannot be inserted again.
func (g *Graph) NewNode(data interface{}, unique bool) (node *Node) {
	if unique {
		node = g.Find(data)
		if node != nil {
			return
		}
	}
	node = newNode(data)
	node.ID = uint32(len(g.Nodes))
	node.graph = g
	g.Nodes = append(g.Nodes, node)

	return
}

func (g *Graph) addNode(node *Node) {
	if node == nil {
		return
	}

	g.Nodes = append(g.Nodes, node)
}

// TopologicalSort will return a slice of the graph nodes,
// sorted topological.
func (g *Graph) TopologicalSort() (sorted []*Node, err error) {
	// Reset node marks
	for _, n := range g.Nodes {
		n.mark = unmarked
	}

	for {
		unmarked := g.findUnmarkedNode()
		if unmarked == nil {
			break
		}

		err = unmarked.topologicalSortVisit(&sorted)
		if err != nil {
			return
		}
	}

	return
}

func (n *Node) topologicalSortVisit(sorted *[]*Node) error {
	if n.mark == temporarilyMarked {
		return fmt.Errorf("Not a DAG")
	}

	if n.mark == unmarked {
		n.mark = temporarilyMarked

		for _, edge := range n.Edges {
			// Exclude edges where n is the dependency
			// Or, find the edges where n points to node2
			if n == edge.Node1 {
				continue
			}

			edge.Node1.topologicalSortVisit(sorted)
		}

		n.mark = permanentlyMarked
		*sorted = append([]*Node{n}, *sorted...)
	}

	return nil
}

func (g *Graph) findUnmarkedNode() *Node {
	for _, n := range g.Nodes {
		if n.mark == unmarked {
			return n
		}
	}

	return nil
}

// Find will find a graph node in the graph, given the ast node.
// If the ast node is not in the graph, nil is returned.
func (g *Graph) Find(data interface{}) *Node {
	for _, n := range g.Nodes {
		if n.Data == data {
			return n
		}
	}
	return nil
}

// FindById will find a node by its given id
func (g *Graph) FindById(i uint32) *Node {
	for _, n := range g.Nodes {
		if n.ID == i {
			return n
		}
	}

	return nil
}

// Size returns the number of nodes in the graph.
func (g *Graph) Size() int {
	return len(g.Nodes)
}

// HasCyclicDependencies returns true if the graph has cyclic dependencies.
func (g *Graph) HasCyclicDependencies() bool {
	for _, n := range g.Nodes {
		deps := []*Node{}
		if n.hasCyclicDependency(deps) {
			return true
		}
	}

	return false
}

// Edge represents an edge between two nodes in a graph.
type Edge struct {
	Node1       *Node
	Node2       *Node
	Data        interface{}
	CrossRegion bool
}

// Remove will remove an edge.
// This will remove this edge from both inbound and outbound nodes.
func (e *Edge) Remove() {
	e.Node1.RemoveEdge(e)
	e.Node2.RemoveEdge(e)
}

// topologicalMark is a mark used for sorting the graph topological
type topologicalMark uint8

const (
	unmarked topologicalMark = iota
	temporarilyMarked
	permanentlyMarked
)

// Node represents a node in a graph
type Node struct {
	// ID is the id of the node
	ID       uint32
	// Data is the data of the node
	Data     interface{}
	// Edges are the edges of the node
	Edges    []*Edge
	// Region defines which region the node belongs to
	Region   interface{}
	// Some data that can be piggy backed on the node
	Metadata interface{}

	// An internal reference to the graph the node is attached to
	graph    *Graph

	// For topological sort
	mark     topologicalMark
}

func newNode(data interface{}) *Node {
	return &Node{
		Data: data,
	}
}

func (n *Node) PutIntoRegion(region interface{}) *Node {
	n.Region = region
	// TODO insert into a map of regions in the graph
	return n
}

// DependOn inserts the other node as a dependency for this node
func (n *Node) DependOn(other *Node) *Edge {
	if n == other {
		return nil
	}
	if n.DependsOn(other) {
		return n.DependsOnAdjacent(other)
	}

	edge := &Edge{
		// The dependent
		Node1: n,
		// The dependency
		Node2: other,
	}

	if n.graph.OnEdgeCreated != nil {
		n.graph.OnEdgeCreated(edge)
	}

	// Insert the edge into node 1 and node 2
	n.Edges = append(n.Edges, edge)
	other.Edges = append(other.Edges, edge)

	return edge
}

// RemoveDependency will remove the given dependency for this node.
func (n *Node) RemoveDependency(other *Node) {
	i := 0
	for _, edge := range n.Edges {
		if edge.Node1 != n || edge.Node2 != other {
			n.Edges[i] = edge
			i++
			continue
		}

		other.RemoveEdge(edge)
	}

	n.Edges = n.Edges[:i]
}

// RemoveEdge will remove a given edge from the node.
// Note, that this will not remove the edge from the other node.
// Use edge.Remove() instead.
func (n *Node) RemoveEdge(e *Edge) {
	i := 0
	for _, edge := range n.Edges {
		if edge != e {
			n.Edges[i] = edge
			i++
		}
	}
	n.Edges = n.Edges[:i]
}

// DependsOn will return true if this node depends on the given node.
func (n *Node) DependsOn(other *Node) bool {
	for _, edge := range n.Edges {
		if edge.Node1 != n {
			continue
		}

		if edge.Node2 == other {
			return true
		}

		if edge.Node2.DependsOn(other) {
			return true
		}
	}

	return false
}

// DependsOnAdjacent will return true if this node is dependent and adjacent to the other node.
func (n *Node) DependsOnAdjacent(other *Node) *Edge {
	for _, edge := range n.Edges {
		if edge.Node1 != n {
			continue
		}

		if edge.Node2 == other {
			return edge
		}
	}

	return nil
}

// DistanceTo will return the distance from one node to another in terms of edges in between.
// Zero means no dependency.
func (n *Node) DistanceTo(other *Node) int {
	shortest := 9999999
	hasShortest := false
	for _, edge := range n.Edges {
		if edge.Node1 != n {
			continue
		}

		if edge.Node2 == other {
			hasShortest = true
			shortest = 1
		}

		d := edge.Node2.DistanceTo(other) + 1
		if d < shortest {
			hasShortest = true
			shortest = d
		}
	}

	if hasShortest {
		return shortest
	} else {
		return -1
	}
}

// DependencyLength will return the number of dependencies for the node.
func (n *Node) DependencyLength() (c int) {
	for _, e := range n.Edges {
		if e.Node1 != n {
			continue
		}

		c += e.Node2.DependencyLength()
		c++
	}

	return
}

// GetDependencies will return a slice of the nodes dependencies.
// unique - only unique nodes will be returned
// all - Not only the adjacent dependency nodes will be returned, but also dependencies dependencies.
func (n *Node) GetDependencies(unique bool, all bool) (deps []*Node) {
	for _, edge := range n.Edges {
		if edge.Node1 != n {
			continue
		}

		deps = append(deps, edge.Node2)
		if all {
			deps = append(deps, edge.Node2.GetDependencies(unique, all)...)
		}
	}

	// Prune duplicates
	if unique {
		unique := make([]*Node, 0)
		for i, n1 := range deps {
			isIn := false
			for j := i + 1; j < len(deps); j++ {
				n2 := deps[j]
				if n1 == n2 {
					isIn = true
					break
				}
			}
			if !isIn {
				unique = append(unique, n1)
			}
		}
		deps = unique
	}

	return
}

func (n *Node) hasCyclicDependency(deps []*Node) bool {
	deps = append(deps, n)
	for _, e := range n.Edges {
		if e.Node1 != n {
			continue
		}

		for _, d := range deps {
			if d == e.Node2 {
				return true
			}
		}

		if e.Node2.hasCyclicDependency(deps) {
			return true
		}
	}

	return false
}

func (n Node) String() string {
	return fmt.Sprintf("Node-%v", n.ID)
}

func (n Node) Stringify() string {
	return n.graph.NodeStringer(n.Data)
}
