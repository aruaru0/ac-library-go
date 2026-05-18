package scc

import (
	"testing"
)

func TestSCCEmpty(t *testing.T) {
	graph0 := NewSCCGraph(0)
	if len(graph0.SCC()) != 0 {
		t.Errorf("Empty graph should have 0 SCCs")
	}
}

func TestSCCSimple(t *testing.T) {
	graph := NewSCCGraph(2)
	graph.AddEdge(0, 1)
	graph.AddEdge(1, 0)
	scc := graph.SCC()
	if len(scc) != 1 {
		t.Errorf("Expected 1 SCC, got %d", len(scc))
	}
	if len(scc[0]) != 2 {
		t.Errorf("Expected SCC to contain 2 vertices, got %d", len(scc[0]))
	}
}

func TestSCCSelfLoop(t *testing.T) {
	graph := NewSCCGraph(2)
	graph.AddEdge(0, 0)
	graph.AddEdge(0, 0)
	graph.AddEdge(1, 1)
	scc := graph.SCC()
	if len(scc) != 2 {
		t.Errorf("Expected 2 SCCs, got %d", len(scc))
	}
}

func TestSCCTopologicalSort(t *testing.T) {
	graph := NewSCCGraph(4)
	graph.AddEdge(0, 1)
	graph.AddEdge(1, 2)
	graph.AddEdge(2, 0)
	graph.AddEdge(3, 2)

	scc := graph.SCC()
	if len(scc) != 2 {
		t.Errorf("Expected 2 SCCs, got %d", len(scc))
	}
	
	// SCCs should be returned in topological order: {3}, {0, 1, 2}
	if len(scc[0]) != 1 || scc[0][0] != 3 {
		t.Errorf("Expected first SCC to be {3}, got %v", scc[0])
	}
	if len(scc[1]) != 3 {
		t.Errorf("Expected second SCC to have 3 vertices, got %v", scc[1])
	}
}

func TestSCCComplex(t *testing.T) {
	graph := NewSCCGraph(6)
	graph.AddEdge(1, 4)
	graph.AddEdge(5, 2)
	graph.AddEdge(3, 0)
	graph.AddEdge(5, 5)
	graph.AddEdge(4, 1)
	graph.AddEdge(0, 3)
	graph.AddEdge(4, 2)
	
	scc := graph.SCC()
	if len(scc) != 4 {
		t.Errorf("Expected 4 SCCs, got %d", len(scc))
	}
	// Expected topological order could be: [5], [1, 4], [2], [0, 3] or similar
	// we just check if it contains all elements correctly
	var flattened []int
	for _, g := range scc {
		flattened = append(flattened, g...)
	}
	if len(flattened) != 6 {
		t.Errorf("Expected total 6 elements, got %d", len(flattened))
	}
}
