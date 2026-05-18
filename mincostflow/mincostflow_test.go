package mincostflow

import (
	"math/rand"
	"reflect"
	"testing"

	"ac-library-go/maxflow"
)

func assertEdgeEq[Cap, Cost Integer](t *testing.T, expect, actual MCFEdge[Cap, Cost]) {
	if expect.From != actual.From || expect.To != actual.To || expect.Cap != actual.Cap || expect.Flow != actual.Flow || expect.Cost != actual.Cost {
		t.Errorf("Edge mismatch. Expect: %+v, Actual: %+v", expect, actual)
	}
}

func TestMincostflowZero(t *testing.T) {
	_ = NewMCFGraph[int, int](0)
}

func TestMincostflowSimple(t *testing.T) {
	g := NewMCFGraph[int, int](4)
	g.AddEdge(0, 1, 1, 1)
	g.AddEdge(0, 2, 1, 1)
	g.AddEdge(1, 3, 1, 1)
	g.AddEdge(2, 3, 1, 1)
	g.AddEdge(1, 2, 1, 1)

	expect := []struct{ Cap, Cost int }{{0, 0}, {2, 4}}
	actual := g.SlopeLimit(0, 3, 10)
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("SlopeLimit failed. Expect: %+v, Actual: %+v", expect, actual)
	}

	assertEdgeEq(t, MCFEdge[int, int]{0, 1, 1, 1, 1}, g.GetEdge(0))
	assertEdgeEq(t, MCFEdge[int, int]{0, 2, 1, 1, 1}, g.GetEdge(1))
	assertEdgeEq(t, MCFEdge[int, int]{1, 3, 1, 1, 1}, g.GetEdge(2))
	assertEdgeEq(t, MCFEdge[int, int]{2, 3, 1, 1, 1}, g.GetEdge(3))
	assertEdgeEq(t, MCFEdge[int, int]{1, 2, 1, 0, 1}, g.GetEdge(4))
}

func TestMincostflowUsage(t *testing.T) {
	g := NewMCFGraph[int, int](2)
	g.AddEdge(0, 1, 1, 2)
	flow, cost := g.Flow(0, 1)
	if flow != 1 || cost != 2 {
		t.Errorf("Flow failed")
	}

	g2 := NewMCFGraph[int, int](2)
	g2.AddEdge(0, 1, 1, 2)
	expect := []struct{ Cap, Cost int }{{0, 0}, {1, 2}}
	if !reflect.DeepEqual(expect, g2.Slope(0, 1)) {
		t.Errorf("Slope failed")
	}
}

func TestMincostflowOutOfRange(t *testing.T) {
	g := NewMCFGraph[int, int](10)

	assertPanic := func(f func(), name string) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic", name)
			}
		}()
		f()
	}

	assertPanic(func() { g.Slope(-1, 3) }, "Slope(-1, 3)")
	assertPanic(func() { g.Slope(3, 3) }, "Slope(3, 3)")
}

func TestMincostflowSelfLoop(t *testing.T) {
	g := NewMCFGraph[int, int](3)
	if g.AddEdge(0, 0, 100, 123) != 0 {
		t.Errorf("AddEdge failed")
	}
	assertEdgeEq(t, MCFEdge[int, int]{0, 0, 100, 0, 123}, g.GetEdge(0))
}

func TestMincostflowSameCostPaths(t *testing.T) {
	g := NewMCFGraph[int, int](3)
	g.AddEdge(0, 1, 1, 1)
	g.AddEdge(1, 2, 1, 0)
	g.AddEdge(0, 2, 2, 1)

	expected := []struct{ Cap, Cost int }{{0, 0}, {3, 3}}
	if !reflect.DeepEqual(expected, g.Slope(0, 2)) {
		t.Errorf("Slope failed")
	}
}

func TestMincostflowInvalid(t *testing.T) {
	g := NewMCFGraph[int, int](2)

	assertPanic := func(f func(), name string) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic", name)
			}
		}()
		f()
	}

	assertPanic(func() { g.AddEdge(0, 0, -1, 0) }, "AddEdge negative cap")
	assertPanic(func() { g.AddEdge(0, 0, 0, -1) }, "AddEdge negative cost")
}

func TestMincostflowStress(t *testing.T) {
	for phase := 0; phase < 1000; phase++ {
		n := rand.Intn(19) + 2
		m := rand.Intn(100) + 1
		s, t_ := rand.Intn(n), rand.Intn(n)
		for s == t_ {
			t_ = rand.Intn(n)
		}
		if rand.Intn(2) == 1 {
			s, t_ = t_, s
		}

		gMf := maxflow.NewMFGraph[int](n)
		g := NewMCFGraph[int, int](n)
		for i := 0; i < m; i++ {
			u := rand.Intn(n)
			v := rand.Intn(n)
			cap := rand.Intn(11)
			cost := rand.Intn(10001)
			g.AddEdge(u, v, cap, cost)
			gMf.AddEdge(u, v, cap)
		}

		flow, cost := g.Flow(s, t_)
		if gMf.Flow(s, t_) != flow {
			t.Errorf("Stress failed: max flow mismatch")
		}

		cost2 := 0
		vCap := make([]int, n)
		for _, e := range g.Edges() {
			vCap[e.From] -= e.Flow
			vCap[e.To] += e.Flow
			cost2 += e.Flow * e.Cost
		}
		if cost != cost2 {
			t.Errorf("Stress failed: cost mismatch")
		}

		for i := 0; i < n; i++ {
			if i == s {
				if vCap[i] != -flow {
					t.Errorf("Stress failed: flow mismatch at s")
				}
			} else if i == t_ {
				if vCap[i] != flow {
					t.Errorf("Stress failed: flow mismatch at t")
				}
			} else {
				if vCap[i] != 0 {
					t.Errorf("Stress failed: flow mismatch at intermediate")
				}
			}
		}

		// check: there is no negative cycle
		dist := make([]int, n)
		for {
			update := false
			for _, e := range g.Edges() {
				if e.Flow < e.Cap {
					ndist := dist[e.From] + e.Cost
					if ndist < dist[e.To] {
						update = true
						dist[e.To] = ndist
					}
				}
				if e.Flow > 0 {
					ndist := dist[e.To] - e.Cost
					if ndist < dist[e.From] {
						update = true
						dist[e.From] = ndist
					}
				}
			}
			if !update {
				break
			}
		}
	}
}
