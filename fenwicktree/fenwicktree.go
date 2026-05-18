package fenwicktree

// Numeric represents the supported types for Fenwick Tree.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// FenwickTree is a data structure that can efficiently update elements
// and calculate prefix sums in an array of numbers.
type FenwickTree[T Numeric] struct {
	n    int
	data []T
}

// New creates a new FenwickTree of length n.
func NewFenwickTree[T Numeric](n int) *FenwickTree[T] {
	if n < 0 {
		panic("n must be non-negative")
	}
	return &FenwickTree[T]{
		n:    n,
		data: make([]T, n),
	}
}

// Add adds x to the p-th element.
func (f *FenwickTree[T]) Add(p int, x T) {
	if p < 0 || p >= f.n {
		panic("p is out of bounds")
	}
	p++
	for p <= f.n {
		f.data[p-1] += x
		p += p & -p
	}
}

// Sum returns the sum of elements in the range [l, r).
func (f *FenwickTree[T]) Sum(l, r int) T {
	if l < 0 || l > r || r > f.n {
		panic("invalid range")
	}
	return f.sum(r) - f.sum(l)
}

func (f *FenwickTree[T]) sum(r int) T {
	var s T
	for r > 0 {
		s += f.data[r-1]
		r -= r & -r
	}
	return s
}
