package twosat

type edge struct {
	to int
}

type csr struct {
	start []int
	elist []edge
}

type pair struct {
	from int
	e    edge
}

func newCSR(n int, edges []pair) *csr {
	start := make([]int, n+1)
	elist := make([]edge, len(edges))
	for _, e := range edges {
		start[e.from+1]++
	}
	for i := 1; i <= n; i++ {
		start[i] += start[i-1]
	}
	counter := make([]int, n+1)
	copy(counter, start)
	for _, e := range edges {
		elist[counter[e.from]] = e.e
		counter[e.from]++
	}
	return &csr{start: start, elist: elist}
}

type sccGraph struct {
	n     int
	edges []pair
}

func newSCCGraph(n int) *sccGraph {
	return &sccGraph{n: n}
}

func (g *sccGraph) addEdge(from, to int) {
	g.edges = append(g.edges, pair{from, edge{to}})
}

func (g *sccGraph) sccIDs() (int, []int) {
	c := newCSR(g.n, g.edges)
	nowOrd := 0
	groupNum := 0
	visited := make([]int, 0, g.n)
	low := make([]int, g.n)
	ord := make([]int, g.n)
	for i := range ord {
		ord[i] = -1
	}
	ids := make([]int, g.n)

	var dfs func(int)
	dfs = func(v int) {
		low[v] = nowOrd
		ord[v] = nowOrd
		nowOrd++
		visited = append(visited, v)
		for i := c.start[v]; i < c.start[v+1]; i++ {
			to := c.elist[i].to
			if ord[to] == -1 {
				dfs(to)
				if low[to] < low[v] {
					low[v] = low[to]
				}
			} else {
				if ord[to] < low[v] {
					low[v] = ord[to]
				}
			}
		}
		if low[v] == ord[v] {
			for {
				u := visited[len(visited)-1]
				visited = visited[:len(visited)-1]
				ord[u] = g.n
				ids[u] = groupNum
				if u == v {
					break
				}
			}
			groupNum++
		}
	}
	for i := 0; i < g.n; i++ {
		if ord[i] == -1 {
			dfs(i)
		}
	}
	for i := range ids {
		ids[i] = groupNum - 1 - ids[i]
	}
	return groupNum, ids
}

// TwoSAT represents a 2-SAT solver.
type TwoSAT struct {
	n      int
	answer []bool
	scc    *sccGraph
}

// New creates a new TwoSAT instance with n variables.
func NewTwoSAT(n int) *TwoSAT {
	if n < 0 {
		panic("n must be non-negative")
	}
	return &TwoSAT{
		n:      n,
		answer: make([]bool, n),
		scc:    newSCCGraph(2 * n),
	}
}

// AddClause adds a clause (x_i == f) || (x_j == g).
func (ts *TwoSAT) AddClause(i int, f bool, j int, g bool) {
	if i < 0 || i >= ts.n {
		panic("i is out of bounds")
	}
	if j < 0 || j >= ts.n {
		panic("j is out of bounds")
	}
	fromI := 2*i
	if !f {
		fromI++
	}
	toJ := 2*j
	if g {
		toJ++
	}
	ts.scc.addEdge(fromI, toJ)

	fromJ := 2*j
	if !g {
		fromJ++
	}
	toI := 2*i
	if f {
		toI++
	}
	ts.scc.addEdge(fromJ, toI)
}

// Satisfiable returns true if the 2-SAT formula is satisfiable, false otherwise.
func (ts *TwoSAT) Satisfiable() bool {
	_, id := ts.scc.sccIDs()
	for i := 0; i < ts.n; i++ {
		if id[2*i] == id[2*i+1] {
			return false
		}
		ts.answer[i] = id[2*i] < id[2*i+1]
	}
	return true
}

// Answer returns the assignment of variables that satisfies the formula.
// The result is only valid if Satisfiable() previously returned true.
func (ts *TwoSAT) Answer() []bool {
	return ts.answer
}
