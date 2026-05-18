package dsu

import (
	"testing"
)

func TestDSUZero(t *testing.T) {
	uf := NewDSU(0)
	groups := uf.Groups()
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %v", groups)
	}
}

func TestDSUSimple(t *testing.T) {
	uf := NewDSU(2)
	if uf.Same(0, 1) {
		t.Errorf("expected 0 and 1 to be different")
	}

	x := uf.Merge(0, 1)

	if x != uf.Leader(0) {
		t.Errorf("expected %d, got %d", x, uf.Leader(0))
	}
	if x != uf.Leader(1) {
		t.Errorf("expected %d, got %d", x, uf.Leader(1))
	}
	if !uf.Same(0, 1) {
		t.Errorf("expected 0 and 1 to be same")
	}
	if uf.Size(0) != 2 {
		t.Errorf("expected size 2, got %d", uf.Size(0))
	}
}

func TestDSULine(t *testing.T) {
	n := 500000
	uf := NewDSU(n)
	for i := 0; i < n-1; i++ {
		uf.Merge(i, i+1)
	}
	if uf.Size(0) != n {
		t.Errorf("expected size %d, got %d", n, uf.Size(0))
	}
	groups := uf.Groups()
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestDSULineReverse(t *testing.T) {
	n := 500000
	uf := NewDSU(n)
	for i := n - 2; i >= 0; i-- {
		uf.Merge(i, i+1)
	}
	if uf.Size(0) != n {
		t.Errorf("expected size %d, got %d", n, uf.Size(0))
	}
	groups := uf.Groups()
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestDSUGroups(t *testing.T) {
	uf := NewDSU(5)
	uf.Merge(0, 1)
	uf.Merge(2, 3)
	uf.Merge(0, 4)

	groups := uf.Groups()
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	// Just checking the sizes to avoid strict order checking
	size3Count := 0
	size2Count := 0
	for _, g := range groups {
		if len(g) == 3 {
			size3Count++
		} else if len(g) == 2 {
			size2Count++
		}
	}
	if size3Count != 1 || size2Count != 1 {
		t.Errorf("expected one group of size 3 and one of size 2, got groups: %v", groups)
	}
}
