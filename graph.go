package graph

import (
	"fmt"
)

// Graph represents a graph
type Graph struct {
	Nodes map[interface{}]*Node
	NodeStringer func(interface{}) string

	OnEdgeCreated func(*Edge)
	Regions map[interface{}][]*Node
}

// NewGraph returns a new graph
func NewGraph() *Graph {
	return &Graph{
		Nodes:make(map[interface{}]*Node),
		Regions:make(map[interface{}][]*Node),
	}
}

func (g *Graph) Stringify() {
	for _, n := range g.Nodes {
		for _, e := range n.Edges {
			if e.Source != n {
				continue
			}

			fmt.Printf("[%v] %v -(%v)> [%v] %v\n", n, g.NodeStringer(n.Data), e.Data, e.Destination, g.NodeStringer(e.Destination.Data))
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
// If the node exists, it will return that node.
func (g *Graph) NewNode(data interface{}) (node *Node) {
	node = g.Find(data)
	if node == nil {
		node = newNode(data)
		node.ID = uint32(len(g.Nodes))
		node.graph = g
		g.Nodes[data] = node
	}

	return
}

func (g *Graph) addNode(node *Node) {
	if node == nil {
		return
	}

	g.Nodes[node.Data] = node
}

// TopologicalSort will return a slice of the graph nodes,
// sorted topological.
func (g *Graph) TopologicalSort(region interface{}) (sorted []*Node, err error) {
	// Reset node marks and set regionality
	if region == nil {
		for _, n := range g.Nodes {
			n.mark = unmarked
			n.localSort = false
		}
	} else {
		for _, n := range g.Regions[region] {
			n.mark = unmarked
			n.localSort = true
		}
	}

	for {
		unmarked := g.findUnmarkedNode(region)
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
			// Or, find the edges where n points to the destination
			if n == edge.Source {
				continue
			}
			// If doing a regional search, exclude those edges where
			// regions do not match.
			if n.localSort && n.Region != edge.Source.Region {
				continue
			}

			edge.Source.topologicalSortVisit(sorted)
		}

		n.mark = permanentlyMarked
		*sorted = append([]*Node{n}, *sorted...)
	}

	return nil
}

func (g *Graph) findUnmarkedNode(region interface{}) *Node {
	if region == nil {
		for _, n := range g.Nodes {
			if n.mark == unmarked {
				return n
			}
		}
	} else {
		for _, n := range g.Regions[region] {
			if n.mark == unmarked {
				return n
			}
		}
	}

	return nil
}

// Find will find a graph node in the graph, given the ast node.
// If the ast node is not in the graph, nil is returned.
func (g *Graph) Find(data interface{}) *Node {
	node, hasNode := g.Nodes[data]
	if hasNode {
		return node
	} else {
		return nil
	}
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

	Source      *Node
	Destination *Node
	Data        interface{}
	CrossRegion bool
}

// Remove will remove an edge.
// This will remove this edge from both inbound and outbound nodes.
func (e *Edge) Remove() {
	e.Source.RemoveEdge(e)
	e.Destination.RemoveEdge(e)
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
	mark      topologicalMark
	localSort bool
}

func newNode(data interface{}) *Node {
	return &Node{
		Data: data,
	}
}

func (n *Node) PutIntoRegion(region interface{}) *Node {
	n.Region = region
	n.graph.Regions[region] = append(n.graph.Regions[region], n)
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
		Source: n,
		// The dependency
		Destination: other,
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
		if edge.Source != n || edge.Destination != other {
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
		if edge.Source != n {
			continue
		}

		if edge.Destination == other {
			return true
		}

		if edge.Destination.DependsOn(other) {
			return true
		}
	}

	return false
}

// DependsOnAdjacent will return true if this node is dependent and adjacent to the other node.
func (n *Node) DependsOnAdjacent(other *Node) *Edge {
	for _, edge := range n.Edges {
		if edge.Source != n {
			continue
		}

		if edge.Destination == other {
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
		if edge.Source != n {
			continue
		}

		if edge.Destination == other {
			hasShortest = true
			shortest = 1
		}

		d := edge.Destination.DistanceTo(other) + 1
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
		if e.Source != n {
			continue
		}

		c += e.Destination.DependencyLength()
		c++
	}

	return
}

// IsDependency will return true if other nodes are
// dependent on this node.
func (n *Node) IsDependency() bool {
	for _, e := range n.Edges {
		if e.Destination != n {
			continue
		}

		return true
	}

	return false
}

// GetDependencies will return a slice of the nodes dependencies.
// unique - only unique nodes will be returned
// all - Not only the adjacent dependency nodes will be returned, but also dependencies dependencies.
func (n *Node) GetDependencies(unique bool, all bool) (deps []*Node) {
	for _, edge := range n.Edges {
		if edge.Source != n {
			continue
		}

		deps = append(deps, edge.Destination)
		if all {
			deps = append(deps, edge.Destination.GetDependencies(unique, all)...)
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
		if e.Source != n {
			continue
		}

		for _, d := range deps {
			if d == e.Destination {
				return true
			}
		}

		if e.Destination.hasCyclicDependency(deps) {
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
