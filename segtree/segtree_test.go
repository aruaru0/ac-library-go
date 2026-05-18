package segtree

import (
	"testing"
)

func opString(a, b string) string {
	if a == "$" {
		return b
	}
	if b == "$" {
		return a
	}
	return a + b
}

func eString() string {
	return "$"
}

type segtreeNaive struct {
	n int
	d []string
}

func newSegtreeNaive(n int) *segtreeNaive {
	d := make([]string, n)
	for i := range d {
		d[i] = eString()
	}
	return &segtreeNaive{n: n, d: d}
}

func (s *segtreeNaive) set(p int, x string) {
	s.d[p] = x
}

func (s *segtreeNaive) get(p int) string {
	return s.d[p]
}

func (s *segtreeNaive) prod(l, r int) string {
	sum := eString()
	for i := l; i < r; i++ {
		sum = opString(sum, s.d[i])
	}
	return sum
}

func (s *segtreeNaive) allProd() string {
	return s.prod(0, s.n)
}

func (s *segtreeNaive) maxRight(l int, f func(string) bool) int {
	sum := eString()
	if !f(sum) {
		panic("f(e) must be true")
	}
	for i := l; i < s.n; i++ {
		sum = opString(sum, s.d[i])
		if !f(sum) {
			return i
		}
	}
	return s.n
}

func (s *segtreeNaive) minLeft(r int, f func(string) bool) int {
	sum := eString()
	if !f(sum) {
		panic("f(e) must be true")
	}
	for i := r - 1; i >= 0; i-- {
		sum = opString(s.d[i], sum)
		if !f(sum) {
			return i + 1
		}
	}
	return 0
}

func TestSegTreeZero(t *testing.T) {
	s := NewSegTree(0, opString, eString)
	if s.AllProd() != "$" {
		t.Errorf("AllProd failed on empty")
	}
}

func TestSegTreeOne(t *testing.T) {
	s := NewSegTree(1, opString, eString)
	if s.AllProd() != "$" {
		t.Errorf("AllProd failed")
	}
	if s.Get(0) != "$" {
		t.Errorf("Get failed")
	}
	if s.Prod(0, 1) != "$" {
		t.Errorf("Prod failed")
	}
	s.Set(0, "dummy")
	if s.Get(0) != "dummy" {
		t.Errorf("Get after Set failed")
	}
	if s.Prod(0, 0) != "$" {
		t.Errorf("Prod 0,0 failed")
	}
	if s.Prod(0, 1) != "dummy" {
		t.Errorf("Prod 0,1 failed")
	}
	if s.Prod(1, 1) != "$" {
		t.Errorf("Prod 1,1 failed")
	}
}

func TestSegTreeCompareNaive(t *testing.T) {
	for n := 0; n < 30; n++ {
		seg0 := newSegtreeNaive(n)
		seg1 := NewSegTree(n, opString, eString)
		for i := 0; i < n; i++ {
			char := string(rune('a' + i))
			seg0.set(i, char)
			seg1.Set(i, char)
		}

		for l := 0; l <= n; l++ {
			for r := l; r <= n; r++ {
				if seg0.prod(l, r) != seg1.Prod(l, r) {
					t.Errorf("Prod mismatch at %d, %d", l, r)
				}
			}
		}

		for l := 0; l <= n; l++ {
			for r := l; r <= n; r++ {
				y := seg1.Prod(l, r)
				f := func(x string) bool { return len(x) <= len(y) }
				if seg0.maxRight(l, f) != seg1.MaxRight(l, f) {
					t.Errorf("MaxRight mismatch at %d", l)
				}
			}
		}

		for r := 0; r <= n; r++ {
			for l := 0; l <= r; l++ {
				y := seg1.Prod(l, r)
				f := func(x string) bool { return len(x) <= len(y) }
				if seg0.minLeft(r, f) != seg1.MinLeft(r, f) {
					t.Errorf("MinLeft mismatch at %d", r)
				}
			}
		}
	}
}

func TestSegTreePanic(t *testing.T) {
	s := NewSegTree(10, opString, eString)

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

	assertPanic(func() { s.MaxRight(11, func(string) bool { return true }) }, "MaxRight(11)")
	assertPanic(func() { s.MinLeft(-1, func(string) bool { return true }) }, "MinLeft(-1)")
	assertPanic(func() { s.MaxRight(0, func(string) bool { return false }) }, "MaxRight(0) with f(e)=false")
}
