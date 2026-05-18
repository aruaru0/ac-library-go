package maxflow

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func assertEdgeEq[Cap Integer](t *testing.T, expect, actual MFEdge[Cap]) {
	if expect.From != actual.From || expect.To != actual.To || expect.Cap != actual.Cap || expect.Flow != actual.Flow {
		t.Errorf("Edge mismatch. Expect: %+v, Actual: %+v", expect, actual)
	}
}

func TestMaxflowZero(t *testing.T) {
	_ = NewMFGraph[int](0)
}

func TestMaxflowSimple(t *testing.T) {
	g := NewMFGraph[int](4)
	if g.AddEdge(0, 1, 1) != 0 {
		t.Errorf("AddEdge failed")
	}
	if g.AddEdge(0, 2, 1) != 1 {
		t.Errorf("AddEdge failed")
	}
	if g.AddEdge(1, 3, 1) != 2 {
		t.Errorf("AddEdge failed")
	}
	if g.AddEdge(2, 3, 1) != 3 {
		t.Errorf("AddEdge failed")
	}
	if g.AddEdge(1, 2, 1) != 4 {
		t.Errorf("AddEdge failed")
	}
	if g.Flow(0, 3) != 2 {
		t.Errorf("Flow failed")
	}

	assertEdgeEq(t, MFEdge[int]{0, 1, 1, 1}, g.GetEdge(0))
	assertEdgeEq(t, MFEdge[int]{0, 2, 1, 1}, g.GetEdge(1))
	assertEdgeEq(t, MFEdge[int]{1, 3, 1, 1}, g.GetEdge(2))
	assertEdgeEq(t, MFEdge[int]{2, 3, 1, 1}, g.GetEdge(3))
	assertEdgeEq(t, MFEdge[int]{1, 2, 1, 0}, g.GetEdge(4))

	if !reflect.DeepEqual([]bool{true, false, false, false}, g.MinCut(0)) {
		t.Errorf("MinCut failed")
	}
}

func TestMaxflowNotSimple(t *testing.T) {
	g := NewMFGraph[int](2)
	g.AddEdge(0, 1, 1)
	g.AddEdge(0, 1, 2)
	g.AddEdge(0, 1, 3)
	g.AddEdge(0, 1, 4)
	g.AddEdge(0, 1, 5)
	g.AddEdge(0, 0, 6)
	g.AddEdge(1, 1, 7)

	if g.Flow(0, 1) != 15 {
		t.Errorf("Flow failed")
	}

	assertEdgeEq(t, MFEdge[int]{0, 1, 1, 1}, g.GetEdge(0))
	assertEdgeEq(t, MFEdge[int]{0, 1, 2, 2}, g.GetEdge(1))
	assertEdgeEq(t, MFEdge[int]{0, 1, 3, 3}, g.GetEdge(2))
	assertEdgeEq(t, MFEdge[int]{0, 1, 4, 4}, g.GetEdge(3))
	assertEdgeEq(t, MFEdge[int]{0, 1, 5, 5}, g.GetEdge(4))

	if !reflect.DeepEqual([]bool{true, false}, g.MinCut(0)) {
		t.Errorf("MinCut failed")
	}
}

func TestMaxflowCut(t *testing.T) {
	g := NewMFGraph[int](3)
	g.AddEdge(0, 1, 2)
	g.AddEdge(1, 2, 1)
	if g.Flow(0, 2) != 1 {
		t.Errorf("Flow failed")
	}
	assertEdgeEq(t, MFEdge[int]{0, 1, 2, 1}, g.GetEdge(0))
	assertEdgeEq(t, MFEdge[int]{1, 2, 1, 1}, g.GetEdge(1))

	if !reflect.DeepEqual([]bool{true, true, false}, g.MinCut(0)) {
		t.Errorf("MinCut failed")
	}
}

func TestMaxflowTwice(t *testing.T) {
	g := NewMFGraph[int](3)
	g.AddEdge(0, 1, 1)
	g.AddEdge(0, 2, 1)
	g.AddEdge(1, 2, 1)

	if g.Flow(0, 2) != 2 {
		t.Errorf("Flow failed")
	}

	assertEdgeEq(t, MFEdge[int]{0, 1, 1, 1}, g.GetEdge(0))
	assertEdgeEq(t, MFEdge[int]{0, 2, 1, 1}, g.GetEdge(1))
	assertEdgeEq(t, MFEdge[int]{1, 2, 1, 1}, g.GetEdge(2))

	g.ChangeEdge(0, 100, 10)
	assertEdgeEq(t, MFEdge[int]{0, 1, 100, 10}, g.GetEdge(0))

	if g.Flow(0, 2) != 0 {
		t.Errorf("Flow failed")
	}
	if g.Flow(0, 1) != 90 {
		t.Errorf("Flow failed")
	}

	assertEdgeEq(t, MFEdge[int]{0, 1, 100, 100}, g.GetEdge(0))
	assertEdgeEq(t, MFEdge[int]{0, 2, 1, 1}, g.GetEdge(1))
	assertEdgeEq(t, MFEdge[int]{1, 2, 1, 1}, g.GetEdge(2))

	if g.Flow(2, 0) != 2 {
		t.Errorf("Flow failed")
	}

	assertEdgeEq(t, MFEdge[int]{0, 1, 100, 99}, g.GetEdge(0))
	assertEdgeEq(t, MFEdge[int]{0, 2, 1, 0}, g.GetEdge(1))
	assertEdgeEq(t, MFEdge[int]{1, 2, 1, 0}, g.GetEdge(2))
}

func TestMaxflowBound(t *testing.T) {
	INF := math.MaxInt32
	g := NewMFGraph[int](3)
	g.AddEdge(0, 1, INF)
	g.AddEdge(1, 0, INF)
	g.AddEdge(0, 2, INF)

	if g.Flow(0, 2) != INF {
		t.Errorf("Flow failed")
	}

	assertEdgeEq(t, MFEdge[int]{0, 1, INF, 0}, g.GetEdge(0))
	assertEdgeEq(t, MFEdge[int]{1, 0, INF, 0}, g.GetEdge(1))
	assertEdgeEq(t, MFEdge[int]{0, 2, INF, INF}, g.GetEdge(2))
}

func TestMaxflowBoundUint(t *testing.T) {
	INF := ^uint32(0)
	g := NewMFGraph[uint32](3)
	g.AddEdge(0, 1, INF)
	g.AddEdge(1, 0, INF)
	g.AddEdge(0, 2, INF)

	if g.Flow(0, 2) != INF {
		t.Errorf("Flow failed")
	}

	assertEdgeEq(t, MFEdge[uint32]{0, 1, INF, 0}, g.GetEdge(0))
	assertEdgeEq(t, MFEdge[uint32]{1, 0, INF, 0}, g.GetEdge(1))
	assertEdgeEq(t, MFEdge[uint32]{0, 2, INF, INF}, g.GetEdge(2))
}

func TestMaxflowSelfLoop(t *testing.T) {
	g := NewMFGraph[int](3)
	g.AddEdge(0, 0, 100)
	assertEdgeEq(t, MFEdge[int]{0, 0, 100, 0}, g.GetEdge(0))
}

func TestMaxflowInvalid(t *testing.T) {
	g := NewMFGraph[int](2)

	assertPanic := func(f func(), name string) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic", name)
			}
		}()
		f()
	}

	assertPanic(func() { g.Flow(0, 0) }, "Flow(0, 0)")
	assertPanic(func() { g.FlowLimit(0, 0, 0) }, "FlowLimit(0, 0, 0)")
}

func TestMaxflowStress(t *testing.T) {
	for phase := 0; phase < 10000; phase++ {
		n := rand.Intn(19) + 2 // 2 to 20
		m := rand.Intn(100) + 1 // 1 to 100
		s, t_ := rand.Intn(n), rand.Intn(n)
		for s == t_ {
			t_ = rand.Intn(n)
		}
		if rand.Intn(2) == 1 {
			s, t_ = t_, s
		}

		g := NewMFGraph[int](n)
		for i := 0; i < m; i++ {
			u := rand.Intn(n)
			v := rand.Intn(n)
			c := rand.Intn(10001)
			g.AddEdge(u, v, c)
		}
		flow := g.Flow(s, t_)
		dual := 0
		cut := g.MinCut(s)
		vFlow := make([]int, n)
		for _, e := range g.Edges() {
			vFlow[e.From] -= e.Flow
			vFlow[e.To] += e.Flow
			if cut[e.From] && !cut[e.To] {
				dual += e.Cap
			}
		}
		if flow != dual {
			t.Errorf("Stress failed: flow != dual")
		}
		if vFlow[s] != -flow {
			t.Errorf("Stress failed: vFlow[s]")
		}
		if vFlow[t_] != flow {
			t.Errorf("Stress failed: vFlow[t]")
		}
		for i := 0; i < n; i++ {
			if i == s || i == t_ {
				continue
			}
			if vFlow[i] != 0 {
				t.Errorf("Stress failed: vFlow[i] != 0")
			}
		}
	}
}
