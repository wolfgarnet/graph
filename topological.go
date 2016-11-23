package graph

import "fmt"

type TopologicalSort struct {
	Nodes        []*Node
	edgeCriteria func(*Node, *Edge) bool
}

func NewTopologicalSort() *TopologicalSort {
	return &TopologicalSort{
		edgeCriteria: func(node *Node, edge *Edge) bool {
			// Exclude edges where n is the dependency
			// Or, find the edges where n points to the destination
			if node == edge.Source {
				return false
			}
			// If doing a regional search, exclude those edges where
			// regions do not match.
			if node.localSort && node.Region != edge.Source.Region {
				return false
			}
			return true
		},
	}
}

func NewCustomTopologicalSort(edgeCriteria func(*Node, *Edge) bool) *TopologicalSort {
	return &TopologicalSort{
		edgeCriteria: edgeCriteria,
	}
}

func (ts *TopologicalSort) Sort(nodes []*Node) (sorted []*Node, err error) {
	for _, n := range nodes {
		n.mark = unmarked
		n.localSort = true
	}

	for {
		unmarked := findUnmarkedNode(nodes)
		if unmarked == nil {
			break
		}

		err = ts.topologicalSortVisit(unmarked, &sorted)
		if err != nil {
			return
		}
	}

	return
}

func (ts *TopologicalSort) topologicalSortVisit(n *Node, sorted *[]*Node) error {
	if n.mark == temporarilyMarked {
		return fmt.Errorf("Not a DAG")
	}

	if n.mark == unmarked {
		n.mark = temporarilyMarked

		for _, edge := range n.Edges {
			// If the edge criteria is not met, the edge is skipped
			if !ts.edgeCriteria(n, edge) {
				continue
			}

			edge.Source.topologicalSortVisit(sorted)
		}

		n.mark = permanentlyMarked
		*sorted = append([]*Node{n}, *sorted...)
	}

	return nil
}

func (ts *TopologicalSort) findUnmarkedNode(nodes []*Node) *Node {
	for _, n := range nodes {
		if n.mark == unmarked {
			return n
		}
	}

	return nil
}
