package dsu

// DSU (Disjoint Set Union) implements (union by size) + (path compression).
// It maintains a collection of disjoint sets, allowing to merge two sets
// and to determine whether two elements belong to the same set.
type DSU struct {
	n             int
	parentOrSize  []int
}

// New creates a new DSU with n vertices and 0 edges.
func NewDSU(n int) *DSU {
	if n < 0 {
		panic("n must be non-negative")
	}
	parentOrSize := make([]int, n)
	for i := 0; i < n; i++ {
		parentOrSize[i] = -1
	}
	return &DSU{
		n:            n,
		parentOrSize: parentOrSize,
	}
}

// Merge adds an edge (a, b).
// It returns the leader of the merged component.
func (d *DSU) Merge(a, b int) int {
	if a < 0 || a >= d.n {
		panic("a is out of bounds")
	}
	if b < 0 || b >= d.n {
		panic("b is out of bounds")
	}
	x, y := d.Leader(a), d.Leader(b)
	if x == y {
		return x
	}
	if -d.parentOrSize[x] < -d.parentOrSize[y] {
		x, y = y, x
	}
	d.parentOrSize[x] += d.parentOrSize[y]
	d.parentOrSize[y] = x
	return x
}

// Same returns whether the vertices a and b are in the same connected component.
func (d *DSU) Same(a, b int) bool {
	if a < 0 || a >= d.n {
		panic("a is out of bounds")
	}
	if b < 0 || b >= d.n {
		panic("b is out of bounds")
	}
	return d.Leader(a) == d.Leader(b)
}

// Leader returns the leader of the connected component that a belongs to.
func (d *DSU) Leader(a int) int {
	if a < 0 || a >= d.n {
		panic("a is out of bounds")
	}
	if d.parentOrSize[a] < 0 {
		return a
	}
	d.parentOrSize[a] = d.Leader(d.parentOrSize[a])
	return d.parentOrSize[a]
}

// Size returns the size of the connected component that a belongs to.
func (d *DSU) Size(a int) int {
	if a < 0 || a >= d.n {
		panic("a is out of bounds")
	}
	return -d.parentOrSize[d.Leader(a)]
}

// Groups divides the graph into connected components and returns them.
func (d *DSU) Groups() [][]int {
	leaderBuf := make([]int, d.n)
	groupSize := make([]int, d.n)
	for i := 0; i < d.n; i++ {
		leaderBuf[i] = d.Leader(i)
		groupSize[leaderBuf[i]]++
	}

	result := make([][]int, d.n)
	for i := 0; i < d.n; i++ {
		result[i] = make([]int, 0, groupSize[i])
	}

	for i := 0; i < d.n; i++ {
		result[leaderBuf[i]] = append(result[leaderBuf[i]], i)
	}

	// Remove empty groups
	var groups [][]int
	for i := 0; i < d.n; i++ {
		if len(result[i]) > 0 {
			groups = append(groups, result[i])
		}
	}
	return groups
}
