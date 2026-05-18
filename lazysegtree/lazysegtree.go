package lazysegtree

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

// LazySegTree represents a segment tree with lazy propagation.
type LazySegTree[S any, F any] struct {
	n           int
	size        int
	log         int
	d           []S
	lz          []F
	op          func(S, S) S
	e           func() S
	mapping     func(F, S) S
	composition func(F, F) F
	id          func() F
}

// New creates a new lazy segment tree of size n.
func NewLazySegTree[S any, F any](
	n int,
	op func(S, S) S,
	e func() S,
	mapping func(F, S) S,
	composition func(F, F) F,
	id func() F,
) *LazySegTree[S, F] {
	v := make([]S, n)
	for i := 0; i < n; i++ {
		v[i] = e()
	}
	return NewLazySegTreeFromSlice(v, op, e, mapping, composition, id)
}

// NewFromSlice creates a new lazy segment tree from a slice.
func NewLazySegTreeFromSlice[S any, F any](
	v []S,
	op func(S, S) S,
	e func() S,
	mapping func(F, S) S,
	composition func(F, F) F,
	id func() F,
) *LazySegTree[S, F] {
	n := len(v)
	size := int(bitCeil(uint32(n)))
	log := countrZero(uint32(size))
	d := make([]S, 2*size)
	for i := range d {
		d[i] = e()
	}
	lz := make([]F, size)
	for i := range lz {
		lz[i] = id()
	}
	for i := 0; i < n; i++ {
		d[size+i] = v[i]
	}
	st := &LazySegTree[S, F]{
		n:           n,
		size:        size,
		log:         log,
		d:           d,
		lz:          lz,
		op:          op,
		e:           e,
		mapping:     mapping,
		composition: composition,
		id:          id,
	}
	for i := size - 1; i >= 1; i-- {
		st.update(i)
	}
	return st
}

func (st *LazySegTree[S, F]) update(k int) {
	st.d[k] = st.op(st.d[2*k], st.d[2*k+1])
}

func (st *LazySegTree[S, F]) allApply(k int, f F) {
	st.d[k] = st.mapping(f, st.d[k])
	if k < st.size {
		st.lz[k] = st.composition(f, st.lz[k])
	}
}

func (st *LazySegTree[S, F]) push(k int) {
	st.allApply(2*k, st.lz[k])
	st.allApply(2*k+1, st.lz[k])
	st.lz[k] = st.id()
}

// Set sets the value at index p to x.
func (st *LazySegTree[S, F]) Set(p int, x S) {
	if p < 0 || p >= st.n {
		panic("p is out of bounds")
	}
	p += st.size
	for i := st.log; i >= 1; i-- {
		st.push(p >> i)
	}
	st.d[p] = x
	for i := 1; i <= st.log; i++ {
		st.update(p >> i)
	}
}

// Get returns the value at index p.
func (st *LazySegTree[S, F]) Get(p int) S {
	if p < 0 || p >= st.n {
		panic("p is out of bounds")
	}
	p += st.size
	for i := st.log; i >= 1; i-- {
		st.push(p >> i)
	}
	return st.d[p]
}

// Prod returns the product of the elements in the range [l, r).
func (st *LazySegTree[S, F]) Prod(l, r int) S {
	if l < 0 || r < l || r > st.n {
		panic("invalid range")
	}
	if l == r {
		return st.e()
	}
	l += st.size
	r += st.size

	for i := st.log; i >= 1; i-- {
		if ((l >> i) << i) != l {
			st.push(l >> i)
		}
		if ((r >> i) << i) != r {
			st.push((r - 1) >> i)
		}
	}

	sml := st.e()
	smr := st.e()
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
func (st *LazySegTree[S, F]) AllProd() S {
	return st.d[1]
}

// Apply applies f to the element at index p.
func (st *LazySegTree[S, F]) Apply(p int, f F) {
	if p < 0 || p >= st.n {
		panic("p is out of bounds")
	}
	p += st.size
	for i := st.log; i >= 1; i-- {
		st.push(p >> i)
	}
	st.d[p] = st.mapping(f, st.d[p])
	for i := 1; i <= st.log; i++ {
		st.update(p >> i)
	}
}

// ApplyRange applies f to the elements in the range [l, r).
func (st *LazySegTree[S, F]) ApplyRange(l, r int, f F) {
	if l < 0 || r < l || r > st.n {
		panic("invalid range")
	}
	if l == r {
		return
	}

	l += st.size
	r += st.size

	for i := st.log; i >= 1; i-- {
		if ((l >> i) << i) != l {
			st.push(l >> i)
		}
		if ((r >> i) << i) != r {
			st.push((r - 1) >> i)
		}
	}

	l2, r2 := l, r
	for l < r {
		if l&1 == 1 {
			st.allApply(l, f)
			l++
		}
		if r&1 == 1 {
			r--
			st.allApply(r, f)
		}
		l >>= 1
		r >>= 1
	}
	l, r = l2, r2

	for i := 1; i <= st.log; i++ {
		if ((l >> i) << i) != l {
			st.update(l >> i)
		}
		if ((r >> i) << i) != r {
			st.update((r - 1) >> i)
		}
	}
}

// MaxRight returns the maximum r <= n such that g(op(a[l], ..., a[r-1])) = true.
func (st *LazySegTree[S, F]) MaxRight(l int, g func(S) bool) int {
	if l < 0 || l > st.n {
		panic("l is out of bounds")
	}
	if !g(st.e()) {
		panic("g(e) must be true")
	}
	if l == st.n {
		return st.n
	}
	l += st.size
	for i := st.log; i >= 1; i-- {
		st.push(l >> i)
	}
	sm := st.e()
	for {
		for l%2 == 0 {
			l >>= 1
		}
		if !g(st.op(sm, st.d[l])) {
			for l < st.size {
				st.push(l)
				l = 2 * l
				if g(st.op(sm, st.d[l])) {
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

// MinLeft returns the minimum l >= 0 such that g(op(a[l], ..., a[r-1])) = true.
func (st *LazySegTree[S, F]) MinLeft(r int, g func(S) bool) int {
	if r < 0 || r > st.n {
		panic("r is out of bounds")
	}
	if !g(st.e()) {
		panic("g(e) must be true")
	}
	if r == 0 {
		return 0
	}
	r += st.size
	for i := st.log; i >= 1; i-- {
		st.push((r - 1) >> i)
	}
	sm := st.e()
	for {
		r--
		for r > 1 && r%2 != 0 {
			r >>= 1
		}
		if !g(st.op(st.d[r], sm)) {
			for r < st.size {
				st.push(r)
				r = 2*r + 1
				if g(st.op(st.d[r], sm)) {
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
