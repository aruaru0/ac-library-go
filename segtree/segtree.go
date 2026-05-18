package segtree

import "math/bits"

func countrZero(n uint32) int {
	return bits.TrailingZeros32(n)
}

func bitCeil(n uint32) uint32 {
	x := uint32(1)
	for x < n {
		x *= 2
	}
	return x
}

// SegTree represents a segment tree.
type SegTree[S any] struct {
	n    int
	size int
	log  int
	d    []S
	op   func(S, S) S
	e    func() S
}

// New creates a new segment tree of size n.
func NewSegTree[S any](n int, op func(S, S) S, e func() S) *SegTree[S] {
	v := make([]S, n)
	for i := 0; i < n; i++ {
		v[i] = e()
	}
	return NewSegTreeFromSlice(v, op, e)
}

// NewFromSlice creates a new segment tree from a slice.
func NewSegTreeFromSlice[S any](v []S, op func(S, S) S, e func() S) *SegTree[S] {
	n := len(v)
	size := int(bitCeil(uint32(n)))
	log := countrZero(uint32(size))
	d := make([]S, 2*size)
	for i := range d {
		d[i] = e()
	}
	for i := 0; i < n; i++ {
		d[size+i] = v[i]
	}
	st := &SegTree[S]{
		n:    n,
		size: size,
		log:  log,
		d:    d,
		op:   op,
		e:    e,
	}
	for i := size - 1; i >= 1; i-- {
		st.update(i)
	}
	return st
}

func (st *SegTree[S]) update(k int) {
	st.d[k] = st.op(st.d[2*k], st.d[2*k+1])
}

// Set sets the value at index p to x.
func (st *SegTree[S]) Set(p int, x S) {
	if p < 0 || p >= st.n {
		panic("p is out of bounds")
	}
	p += st.size
	st.d[p] = x
	for i := 1; i <= st.log; i++ {
		st.update(p >> i)
	}
}

// Get returns the value at index p.
func (st *SegTree[S]) Get(p int) S {
	if p < 0 || p >= st.n {
		panic("p is out of bounds")
	}
	return st.d[p+st.size]
}

// Prod returns the product of the elements in the range [l, r).
func (st *SegTree[S]) Prod(l, r int) S {
	if l < 0 || r < l || r > st.n {
		panic("invalid range")
	}
	sml := st.e()
	smr := st.e()
	l += st.size
	r += st.size

	for l < r {
		if l&1 == 1 {
			sml = st.op(sml, st.d[l])
			l++
		}
		if r&1 == 1 {
			r--
			smr = st.op(st.d[r], smr)
		}
		l >>= 1
		r >>= 1
	}
	return st.op(sml, smr)
}

// AllProd returns the product of all elements.
func (st *SegTree[S]) AllProd() S {
	return st.d[1]
}

// MaxRight returns the maximum r <= n such that f(op(a[l], ..., a[r-1])) = true.
func (st *SegTree[S]) MaxRight(l int, f func(S) bool) int {
	if l < 0 || l > st.n {
		panic("l is out of bounds")
	}
	if !f(st.e()) {
		panic("f(e) must be true")
	}
	if l == st.n {
		return st.n
	}
	l += st.size
	sm := st.e()
	for {
		for l%2 == 0 {
			l >>= 1
		}
		if !f(st.op(sm, st.d[l])) {
			for l < st.size {
				l = 2 * l
				if f(st.op(sm, st.d[l])) {
					sm = st.op(sm, st.d[l])
					l++
				}
			}
			return l - st.size
		}
		sm = st.op(sm, st.d[l])
		l++
		if (l & -l) == l {
			break
		}
	}
	return st.n
}

// MinLeft returns the minimum l >= 0 such that f(op(a[l], ..., a[r-1])) = true.
func (st *SegTree[S]) MinLeft(r int, f func(S) bool) int {
	if r < 0 || r > st.n {
		panic("r is out of bounds")
	}
	if !f(st.e()) {
		panic("f(e) must be true")
	}
	if r == 0 {
		return 0
	}
	r += st.size
	sm := st.e()
	for {
		r--
		for r > 1 && r%2 != 0 {
			r >>= 1
		}
		if !f(st.op(st.d[r], sm)) {
			for r < st.size {
				r = 2*r + 1
				if f(st.op(st.d[r], sm)) {
					sm = st.op(st.d[r], sm)
					r--
				}
			}
			return r + 1 - st.size
		}
		sm = st.op(st.d[r], sm)
		if (r & -r) == r {
			break
		}
	}
	return 0
}
