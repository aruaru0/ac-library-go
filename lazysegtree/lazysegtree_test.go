package lazysegtree

import (
	"testing"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func opSS(a, b int) int { return maxInt(a, b) }
func opTS(a, b int) int { return a + b }
func opTT(a, b int) int { return a + b }
func eS() int           { return -1000000000 }
func eT() int           { return 0 }

func TestLazySegTreeZero(t *testing.T) {
	s0 := NewLazySegTree(0, opSS, eS, opTS, opTT, eT)
	if s0.AllProd() != -1000000000 {
		t.Errorf("Zero length AllProd failed")
	}

	s10 := NewLazySegTree(10, opSS, eS, opTS, opTT, eT)
	if s10.AllProd() != -1000000000 {
		t.Errorf("10 length AllProd failed")
	}
}

func TestLazySegTreePanic(t *testing.T) {
	s := NewLazySegTree(10, opSS, eS, opTS, opTT, eT)

	assertPanic := func(f func(), name string) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic", name)
			}
		}()
		f()
	}

	assertPanic(func() { s.Get(-1) }, "Get(-1)")
	assertPanic(func() { s.Get(10) }, "Get(10)")

	assertPanic(func() { s.Prod(-1, -1) }, "Prod(-1, -1)")
	assertPanic(func() { s.Prod(3, 2) }, "Prod(3, 2)")
	assertPanic(func() { s.Prod(0, 11) }, "Prod(0, 11)")
	assertPanic(func() { s.Prod(-1, 11) }, "Prod(-1, 11)")
}

func TestLazySegTreeNaiveProd(t *testing.T) {
	for n := 0; n <= 50; n++ {
		seg := NewLazySegTree(n, opSS, eS, opTS, opTT, eT)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			p[i] = (i*i + 100) % 31
			seg.Set(i, p[i])
		}
		for l := 0; l <= n; l++ {
			for r := l; r <= n; r++ {
				e := -1000000000
				for i := l; i < r; i++ {
					e = maxInt(e, p[i])
				}
				if seg.Prod(l, r) != e {
					t.Errorf("NaiveProd mismatch at n=%d, l=%d, r=%d", n, l, r)
				}
			}
		}
	}
}

func TestLazySegTreeUsage(t *testing.T) {
	v := make([]int, 10)
	seg := NewLazySegTreeFromSlice(v, opSS, eS, opTS, opTT, eT)
	if seg.AllProd() != 0 {
		t.Errorf("Usage: AllProd expected 0")
	}
	seg.ApplyRange(0, 3, 5)
	if seg.AllProd() != 5 {
		t.Errorf("Usage: AllProd expected 5")
	}
	seg.Apply(2, -10)
	if seg.Prod(2, 3) != -5 {
		t.Errorf("Usage: Prod(2,3) expected -5")
	}
	if seg.Prod(2, 4) != 0 {
		t.Errorf("Usage: Prod(2,4) expected 0")
	}
}
