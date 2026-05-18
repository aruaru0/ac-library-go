package scc

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

// SCCGraph represents a directed graph for finding strongly connected components.
type SCCGraph struct {
	n     int
	edges []pair
}

// New creates a new SCCGraph with n vertices.
func NewSCCGraph(n int) *SCCGraph {
	if n < 0 {
		panic("n must be non-negative")
	}
	return &SCCGraph{n: n}
}

// AddEdge adds a directed edge from `from` to `to`.
func (g *SCCGraph) AddEdge(from, to int) {
	if from < 0 || from >= g.n {
		panic("from is out of bounds")
	}
	if to < 0 || to >= g.n {
		panic("to is out of bounds")
	}
	g.edges = append(g.edges, pair{from, edge{to}})
}

func (g *SCCGraph) sccIDs() (int, []int) {
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

// SCC returns the strongly connected components of the graph.
// The returned groups are sorted in topological order.
func (g *SCCGraph) SCC() [][]int {
	groupNum, ids := g.sccIDs()
	counts := make([]int, groupNum)
	for _, x := range ids {
		counts[x]++
	}
	groups := make([][]int, groupNum)
	for i := 0; i < groupNum; i++ {
		groups[i] = make([]int, 0, counts[i])
	}
	for i := 0; i < g.n; i++ {
		groups[ids[i]] = append(groups[ids[i]], i)
	}
	return groups
}
