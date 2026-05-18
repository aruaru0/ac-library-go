package maxflow

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type MFEdge[Cap Integer] struct {
	From int
	To   int
	Cap  Cap
	Flow Cap
}

type _edge[Cap Integer] struct {
	to  int
	rev int
	cap Cap
}

// Graph represents a graph for max flow.
type MFGraph[Cap Integer] struct {
	n   int
	pos []struct {
		first  int
		second int
	}
	g [][]_edge[Cap]
}

// New creates a new max flow graph with n vertices.
func NewMFGraph[Cap Integer](n int) *MFGraph[Cap] {
	return &MFGraph[Cap]{
		n: n,
		g: make([][]_edge[Cap], n),
	}
}

// AddEdge adds a directed edge from `from` to `to` with capacity `cap`.
func (g *MFGraph[Cap]) AddEdge(from, to int, cap Cap) int {
	if from < 0 || from >= g.n {
		panic("from is out of bounds")
	}
	if to < 0 || to >= g.n {
		panic("to is out of bounds")
	}
	if cap < 0 {
		panic("capacity cannot be negative")
	}
	m := len(g.pos)
	g.pos = append(g.pos, struct {
		first  int
		second int
	}{from, len(g.g[from])})
	fromId := len(g.g[from])
	toId := len(g.g[to])
	if from == to {
		toId++
	}
	g.g[from] = append(g.g[from], _edge[Cap]{to, toId, cap})
	g.g[to] = append(g.g[to], _edge[Cap]{from, fromId, 0})
	return m
}

// GetEdge returns the i-th edge.
func (g *MFGraph[Cap]) GetEdge(i int) MFEdge[Cap] {
	if i < 0 || i >= len(g.pos) {
		panic("index out of bounds")
	}
	_e := g.g[g.pos[i].first][g.pos[i].second]
	_re := g.g[_e.to][_e.rev]
	return MFEdge[Cap]{
		From: g.pos[i].first,
		To:   _e.to,
		Cap:  _e.cap + _re.cap,
		Flow: _re.cap,
	}
}

// Edges returns all edges.
func (g *MFGraph[Cap]) Edges() []MFEdge[Cap] {
	m := len(g.pos)
	result := make([]MFEdge[Cap], m)
	for i := 0; i < m; i++ {
		result[i] = g.GetEdge(i)
	}
	return result
}

// ChangeEdge changes the capacity and flow of the i-th edge.
func (g *MFGraph[Cap]) ChangeEdge(i int, newCap, newFlow Cap) {
	if i < 0 || i >= len(g.pos) {
		panic("index out of bounds")
	}
	if newFlow < 0 || newFlow > newCap {
		panic("invalid flow or capacity")
	}
	first := g.pos[i].first
	second := g.pos[i].second
	_e := &g.g[first][second]
	_re := &g.g[_e.to][_e.rev]
	_e.cap = newCap - newFlow
	_re.cap = newFlow
}

// Flow calculates the maximum flow from s to t.
func (g *MFGraph[Cap]) Flow(s, t int) Cap {
	var maxCap Cap = 1
	for {
		if maxCap*2+1 > maxCap {
			maxCap = maxCap*2 + 1
		} else {
			break
		}
	}
	return g.FlowLimit(s, t, maxCap)
}

func minCap[Cap Integer](a, b Cap) Cap {
	if a < b {
		return a
	}
	return b
}

// FlowLimit calculates the maximum flow from s to t up to a flow limit.
func (g *MFGraph[Cap]) FlowLimit(s, t int, flowLimit Cap) Cap {
	if s < 0 || s >= g.n {
		panic("s is out of bounds")
	}
	if t < 0 || t >= g.n {
		panic("t is out of bounds")
	}
	if s == t {
		panic("s and t must be different")
	}

	level := make([]int, g.n)
	iter := make([]int, g.n)
	que := make([]int, 0, g.n)

	bfs := func() {
		for i := range level {
			level[i] = -1
		}
		level[s] = 0
		que = que[:0]
		que = append(que, s)
		for len(que) > 0 {
			v := que[0]
			que = que[1:]
			for _, e := range g.g[v] {
				if e.cap == 0 || level[e.to] >= 0 {
					continue
				}
				level[e.to] = level[v] + 1
				if e.to == t {
					return
				}
				que = append(que, e.to)
			}
		}
	}

	var dfs func(v int, up Cap) Cap
	dfs = func(v int, up Cap) Cap {
		if v == s {
			return up
		}
		res := Cap(0)
		levelV := level[v]
		for ; iter[v] < len(g.g[v]); iter[v]++ {
			i := iter[v]
			e := &g.g[v][i]
			if levelV <= level[e.to] || g.g[e.to][e.rev].cap == 0 {
				continue
			}
			d := dfs(e.to, minCap(up-res, g.g[e.to][e.rev].cap))
			if d <= 0 {
				continue
			}
			g.g[v][i].cap += d
			g.g[e.to][e.rev].cap -= d
			res += d
			if res == up {
				return res
			}
		}
		level[v] = g.n
		return res
	}

	flow := Cap(0)
	for flow < flowLimit {
		bfs()
		if level[t] == -1 {
			break
		}
		for i := range iter {
			iter[i] = 0
		}
		f := dfs(t, flowLimit-flow)
		if f == 0 {
			break
		}
		flow += f
	}
	return flow
}

// MinCut returns the minimum cut from s.
func (g *MFGraph[Cap]) MinCut(s int) []bool {
	visited := make([]bool, g.n)
	que := make([]int, 0, g.n)
	que = append(que, s)
	visited[s] = true
	for len(que) > 0 {
		p := que[0]
		que = que[1:]
		for _, e := range g.g[p] {
			if e.cap > 0 && !visited[e.to] {
				visited[e.to] = true
				que = append(que, e.to)
			}
		}
	}
	return visited
}
