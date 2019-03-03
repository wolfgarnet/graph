package graph

import "fmt"

// Node represents a node in a graph
type Node struct {
	// ID is the id of the node
	ID uint32

	// Data is the data of the node
	Data interface{}

	// Edges are the edges of the node
	Edges []*Edge

	// Region defines which region the node belongs to
	Region interface{}

	// Some data that can be piggy backed on the node
	Metadata interface{}

	// An internal reference to the graph the node is attached to
	graph *Graph

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
		if n.graph.OnSameNodeEdge != nil {
			n.graph.OnSameNodeEdge(n)
		}
		return nil
	}
	d := n.DependsOnAdjacent(other)
	if d != nil {
		if n.graph.OnDuplicateEdge != nil {
			n.graph.OnDuplicateEdge(d)
		}
		return d
	}

	// The source points to destination
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

// DependOn2 does a dependency check before adding the edge
func (n *Node) DependOn2(other *Node) *Edge {
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
//
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
// Points from this node
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

// GetDependent will return a slice of nodes that depends on this node.
// Points to this node
// unique - only unique nodes will be returned
// all - Not only the adjacent dependency nodes will be returned, but also dependencies dependencies.
func (n *Node) GetDependents(unique bool, all bool) (deps []*Node) {
	for _, edge := range n.Edges {
		if edge.Destination != n {
			continue
		}

		deps = append(deps, edge.Source)
		if all {
			deps = append(deps, edge.Source.GetDependents(unique, all)...)
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
