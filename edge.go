package graph

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

// CreateLink is an auxiliary function for creating a directional edge from one node to another.
func CreateLink(from, to *Node) *Edge {
	return from.DependOn(to)
}
