package graph

type Walk struct {
	FollowEdge func(node *Node, edge *Edge) bool
	CallBack   func(node *Node, edge *Edge)
}

func NewDepthFirstWalker(cb func(*Node, *Edge)) *Walk {
	return &Walk{
		FollowEdge:func(node *Node, edge *Edge) bool {
			if node != edge.Source {
				return false
			}

			return true
		},
		CallBack:cb,
	}
}

func NewDepthFirstWalkerWithinSameRegion(cb func(*Node, *Edge)) *Walk {
	return &Walk{
		FollowEdge:func(node *Node, edge *Edge) bool {
			if node != edge.Source {
				return false
			}
			if node.Region != edge.Destination.Region {
				return false
			}

			return true
		},
		CallBack:cb,
	}
}

func (n *Node) Walk(w *Walk) {
	w.Walk(n)
}

func (w *Walk) Walk(n *Node) {
	for _, edge := range n.Edges {
		if !w.FollowEdge(n, edge) {
			continue
		}

		w.CallBack(n, edge)
	}
}
