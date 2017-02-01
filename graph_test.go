package graph

import (
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {
	g := NewGraph()

	g.NewNode(1)
	g.NewNode(2)

	if len(g.Nodes) != 2 {
		t.Errorf("The graph does not have two root nodes")
	}
}

func TestGraph_size(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)

	n1.DependOn(n2)

	if g.Size() != 2 {
		t.Errorf("The graph does not have two nodes")
	}
}

func TestGraph_size2(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)

	n1.DependOn(n2)
	n1.DependOn(n3)

	n2.DependOn(n3)

	if g.Size() != 3 {
		t.Errorf("The graph does not have three nodes, but %v", g.Size())
	}
}

func TestGraph_size3(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)

	n1.DependOn(n2)
	n2.DependOn(n3)
	n3.DependOn(n4)

	if g.Size() != 4 {
		t.Errorf("The graph does not have four nodes, but %v", g.Size())
	}
}

func TestGraph_dependencies(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)

	n1.DependOn(n2)
	n2.DependOn(n3)
	n3.DependOn(n4)

	n1numdeps := len(n1.GetDependencies(true, true))
	if n1numdeps != 3 {
		t.Errorf("The graph does not have three dependencies, but %v", n1numdeps)
	}

	n2numdeps := len(n2.GetDependencies(true, true))
	if n2numdeps != 2 {
		t.Errorf("The graph does not have two dependencies, but %v", n2numdeps)
	}

	n3numdeps := len(n3.GetDependencies(true, true))
	if n3numdeps != 1 {
		t.Errorf("The graph does not have one dependency, but %v", n3numdeps)
	}

	n4numdeps := len(n4.GetDependencies(true, true))
	if n4numdeps != 0 {
		t.Errorf("The graph does not have zero dependencies, but %v", n4numdeps)
	}
}

func TestGraph_uniqueness(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)
	n5 := g.NewNode(5)
	n6 := g.NewNode(6)
	n7 := g.NewNode(7)

	n1.DependOn(n2)
	n1.DependOn(n3)
	n2.DependOn(n4)
	n3.DependOn(n4)
	n4.DependOn(n5)
	n4.DependOn(n6)
	n5.DependOn(n7)
	n6.DependOn(n7)

	n1numdepsunique := len(n1.GetDependencies(true, true))
	if n1numdepsunique != 6 {
		t.Errorf("The node 1 does not have six unique dependencies, but %v", n1numdepsunique)
	}

	n1numdeps := len(n1.GetDependencies(false, true))
	if n1numdeps != 12 {
		t.Errorf("The node 1 does not have 12 total dependencies, but %v", n1numdeps)
	}

}

func TestGraph_uniqueness2(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)
	n5 := g.NewNode(5)
	n6 := g.NewNode(6)
	n7 := g.NewNode(7)

	n1.DependOn(n2)
	n1.DependOn(n3)
	n2.DependOn(n4)
	n3.DependOn(n4)
	n4.DependOn(n5)
	n4.DependOn(n6)
	n5.DependOn(n7)
	n6.DependOn(n7)

	n1numdepsunique := len(n1.GetDependencies(true, false))
	if n1numdepsunique != 2 {
		t.Errorf("The node 1 does not have two unique dependencies, but %v", n1numdepsunique)
	}

	n1numdeps := len(n1.GetDependencies(false, false))
	if n1numdeps != 2 {
		t.Errorf("The node 1 does not have two total dependencies, but %v", n1numdeps)
	}

}

func TestGraph_dependencies2(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)

	n1.DependOn(n2)
	n1.DependOn(n3)
	n2.DependOn(n4)
	n3.DependOn(n4)

	n1numdeps := len(n1.GetDependencies(true, true))
	if n1numdeps != 3 {
		t.Errorf("The graph does not have three dependencies, but %v", n1numdeps)
	}

	n2numdeps := len(n2.GetDependencies(true, true))
	if n2numdeps != 1 {
		t.Errorf("The graph does not have two dependencies, but %v", n2numdeps)
	}

	n3numdeps := len(n3.GetDependencies(true, true))
	if n3numdeps != 1 {
		t.Errorf("The graph does not have one dependency, but %v", n3numdeps)
	}

	n4numdeps := len(n4.GetDependencies(true, true))
	if n4numdeps != 0 {
		t.Errorf("The graph does not have zero dependencies, but %v", n4numdeps)
	}
}

func TestGraph_dependencies3(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)

	n1.DependOn(n2)
	n1.DependOn(n3)
	n2.DependOn(n4)
	n3.DependOn(n4)

	n1numdeps := len(n1.GetDependencies(true, false))
	if n1numdeps != 2 {
		t.Errorf("The graph does not have two dependencies, but %v", n1numdeps)
	}

	n2numdeps := len(n2.GetDependencies(true, false))
	if n2numdeps != 1 {
		t.Errorf("The graph does not have two dependencies, but %v", n2numdeps)
	}

	n3numdeps := len(n3.GetDependencies(true, false))
	if n3numdeps != 1 {
		t.Errorf("The graph does not have one dependency, but %v", n3numdeps)
	}

	n4numdeps := len(n4.GetDependencies(true, false))
	if n4numdeps != 0 {
		t.Errorf("The graph does not have zero dependencies, but %v", n4numdeps)
	}
}

func TestGraph_cyclic(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)

	n1.DependOn(n2)
	n2.DependOn(n3)
	n3.DependOn(n1)

	if !g.HasCyclicDependencies() {
		t.Errorf("The graph should have a cyclic dependency")
	}
}

func TestGraph_cyclic2(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)

	n1.DependOn(n2)
	n2.DependOn(n3)
	n3.DependOn(n4)
	n4.DependOn(n2)

	if !g.HasCyclicDependencies() {
		t.Errorf("The graph should have a cyclic dependency")
	}
}

func TestGraph_cyclic3(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)

	n1.DependOn(n2)
	n1.DependOn(n3)
	n2.DependOn(n4)
	n3.DependOn(n4)

	if g.HasCyclicDependencies() {
		t.Errorf("The graph should NOT have a cyclic dependency")
	}
}

func TestGraph_dependsOn(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)
	n5 := g.NewNode(5)
	n6 := g.NewNode(6)

	n1.DependOn(n2)
	n1.DependOn(n6)
	n2.DependOn(n3)
	n2.DependOn(n5)
	n2.DependOn(n4)
	n4.DependOn(n5)

	if !n1.DependsOn(n6) {
		t.Errorf("%v should depend on %v\n", n1, n6)
	}

	if !n1.DependsOn(n5) {
		t.Errorf("%v should depend on %v\n", n1, n5)
	}

	if !n4.DependsOn(n5) {
		t.Errorf("%v should depend on %v\n", n4, n5)
	}

	if n6.DependsOn(n1) {
		t.Errorf("%v should NOT depend on %v\n", n6, n1)
	}

	if n5.DependsOn(n4) {
		t.Errorf("%v should NOT depend on %v\n", n5, n4)
	}
}

func TestGraph_dependent(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)
	n5 := g.NewNode(5)
	n6 := g.NewNode(6)

	n1.DependOn(n2)
	n1.DependOn(n6)
	n2.DependOn(n3)
	n2.DependOn(n5)
	n2.DependOn(n4)
	n4.DependOn(n5)

	deps := n5.GetDependents(true, false)
	if len(deps) != 2 {
		t.Errorf("Number of dependents should be 2")
	}

	if !n2.DependsOn(n5) {
		t.Errorf("Node 2 should depend on node 5")
	}

	if !n4.DependsOn(n5) {
		t.Errorf("Node 4 should depend on node 5")
	}
}

func TestGraph_topologicalSort(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)

	n1.DependOn(n2)
	n2.DependOn(n3)

	sorted, err := g.TopologicalSort(nil)
	if err != nil {
		t.Errorf("Unable to sort topological: %v", err)
	}

	err = validateTopologicalSort(sorted)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
}

func TestGraph_topologicalSort2(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)

	n1.DependOn(n2)
	n1.DependOn(n3)
	n2.DependOn(n4)
	n3.DependOn(n4)

	sorted, err := g.TopologicalSort(nil)
	if err != nil {
		t.Errorf("Unable to sort topological: %v", err)
	}

	err = validateTopologicalSort(sorted)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
}

func TestGraph_topologicalSort3(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)
	n5 := g.NewNode(5)
	n6 := g.NewNode(6)
	n7 := g.NewNode(7)

	n1.DependOn(n2)
	n1.DependOn(n3)
	n2.DependOn(n5)
	n3.DependOn(n4)
	n3.DependOn(n6)
	n4.DependOn(n7)
	n6.DependOn(n7)

	sorted, err := g.TopologicalSort(nil)
	if err != nil {
		t.Errorf("Unable to sort topological: %v", err)
	}

	err = validateTopologicalSort(sorted)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
}

func TestGraph_topologicalSort4(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)
	n5 := g.NewNode(5)
	n6 := g.NewNode(6)
	n7 := g.NewNode(7)

	n1.DependOn(n2)
	n1.DependOn(n3)
	n2.DependOn(n5)
	n3.DependOn(n4)
	n3.DependOn(n6)
	n4.DependOn(n7)
	n5.DependOn(n4)
	n6.DependOn(n7)

	sorted, err := g.TopologicalSort(nil)
	if err != nil {
		t.Errorf("Unable to sort topological: %v", err)
	}

	err = validateTopologicalSort(sorted)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
}

func TestGraph_topologicalSort5(t *testing.T) {
	g := NewGraph()

	n0 := g.NewNode(0)
	n1 := g.NewNode(1)
	n2 := g.NewNode(2)
	n3 := g.NewNode(3)
	n4 := g.NewNode(4)
	n5 := g.NewNode(5)
	n6 := g.NewNode(6)
	n7 := g.NewNode(7)
	n8 := g.NewNode(8)
	n9 := g.NewNode(9)

	n9.DependOn(n8)
	n7.DependOn(n6)
	n6.DependOn(n5)
	n5.DependOn(n4)
	n4.DependOn(n1)
	n3.DependOn(n0)
	n5.DependOn(n3)
	n5.DependOn(n2)
	n8.DependOn(n6)

	nodes := []*Node{n1, n3, n5, n2, n8, n7, n0, n4, n9, n6}

	sorted, err := SortTopological(nodes)
	if err != nil {
		t.Errorf("Unable to sort topological: %v", err)
	}

	err = validateTopologicalSort(sorted)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
}

func TestGraph_topologicalSortRegions(t *testing.T) {
	g := NewGraph()

	n1 := g.NewNode(1).PutIntoRegion(1)
	n2 := g.NewNode(2).PutIntoRegion(1)
	n3 := g.NewNode(3).PutIntoRegion(1)
	n4 := g.NewNode(4).PutIntoRegion(1)
	n5 := g.NewNode(5).PutIntoRegion(2)
	n6 := g.NewNode(6).PutIntoRegion(2)
	n7 := g.NewNode(7).PutIntoRegion(2)
	n8 := g.NewNode(8).PutIntoRegion(2)

	n1.DependOn(n2)
	n2.DependOn(n4)
	n2.DependOn(n3)
	n4.DependOn(n3)
	n5.DependOn(n1)
	n6.DependOn(n5)
	n7.DependOn(n6)
	n8.DependOn(n7)

	sorted, err := g.TopologicalSort(1)
	if err != nil {
		t.Errorf("Unable to sort topological: %v", err)
	}

	if len(sorted) != 4 {
		t.Errorf("The number of topologically sorted nodes are not 4, but %v", len(sorted))
	}

	err = validateTopologicalSort(sorted)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}
}

// Determine if the order of the topological sort is correct.
// For each node at position i, each of its dependencies must be
// at a position lower than it self.
// Or, nodes at a given index cannot have dependencies at a higher index.
func validateTopologicalSort(sorted []*Node) error {
	for i, n := range sorted {
		for _, e := range n.Edges {
			if n == e.Destination {
				continue
			}

			dependencyIdx := findNode(sorted, e.Destination)
			if dependencyIdx > i {
				return fmt.Errorf("Topological sort order failed: %v(%v) < %v(%v)\n", n, i, e.Destination, dependencyIdx)
			}
		}
	}

	return nil
}

func findNode(l []*Node, n *Node) int {
	for i, node := range l {
		if n == node {
			return i
		}
	}

	return -1
}
