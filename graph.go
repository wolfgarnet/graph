package graph

import (
	"fmt"
	"sort"
)

// Graph represents a graph
type Graph struct {
	// Nodes contains all the nodes for the graph
	Nodes map[interface{}]*Node

	// NodeStringer is a function for stringifying nodes
	NodeStringer func(interface{}) string

	// OnNodeCreated is called when a node is created
	OnNodeCreated func(*Node)

	// OnEdgeCreated is called when an edge is created
	OnEdgeCreated func(*Edge)

	// OnSameNodeEdge is called when an edge is created that has same source and destination
	OnSameNodeEdge func(*Node)

	// OnDuplicateEdge is called when trying to create an edge that already exists
	OnDuplicateEdge func(*Edge)

	// Regions is used to group nodes into regions
	Regions map[interface{}][]*Node
}

// NewGraph returns a new graph
func NewGraph() *Graph {
	return &Graph{
		Nodes:   make(map[interface{}]*Node),
		Regions: make(map[interface{}][]*Node),
	}
}

// Stringify stringifies the graph
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

func (g *Graph) NumberOfEdges() int {
	c := 0
	for _, n := range g.Nodes {
		for _, e := range n.Edges {
			if e.Source != n {
				continue
			}

			c++
		}
	}

	return c
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

func (g *Graph) Print() {
	edges := make(map[*Edge]struct{ id, num int })
	id := 0
	edgeCounter := 0

	ids := make([]int, len(g.Nodes))
	i := 0
	for _, n := range g.Nodes {
		ids[i] = int(n.ID)
		i++
	}

	sort.Ints(ids)

	//for _, node := range g.Nodes {
	for _, i := range ids {
		node := g.FindById(uint32(i))
		if node.Data != nil {
			if g.NodeStringer != nil {
				fmt.Printf("[%v]: %v(%T)\n", node, g.NodeStringer(node.Data), node.Data)
			} else {
				fmt.Printf("[%v]: %v(%T)\n", node, node.Data, node.Data)
			}
		} else {
			fmt.Printf("[%v]: NIL(%T)\n", node, node.Data)
		}

		for i, edge := range node.Edges {
			edgeCounter++
			value, hasEdge := edges[edge]
			if !hasEdge {
				value = struct{ id, num int }{id, 1}
				edges[edge] = value
				id++
			} else {
				value.num++
				edges[edge] = value
			}

			inbound := edge.Destination == node
			if inbound {
				fmt.Printf("  [%v] %v <- %v (%p) [%v/%v]\n", i, edge.Destination.ID, edge.Source.ID, edge, value.id, value.num)
			} else {
				fmt.Printf("  [%v] %v -> %v (%p) [%v/%v]\n", i, edge.Source.ID, edge.Destination.ID, edge, value.id, value.num)
			}
		}
	}

	fmt.Printf("Nodes: %v, Unique edges: %v, total edges: %v\n", len(g.Nodes), len(edges), edgeCounter)
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
		if g.OnNodeCreated != nil {
			g.OnNodeCreated(node)
		}
	}

	return
}

func (g *Graph) addNode(node *Node) {
	if node == nil {
		return
	}

	g.Nodes[node.Data] = node
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

func (g *Graph) SizeEdges() int {
	edges := make(map[*Edge]bool)
	for _, n := range g.Nodes {
		for _, e := range n.Edges {
			edges[e] = true
		}
	}

	return len(edges)
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
