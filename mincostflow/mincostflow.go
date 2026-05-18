package mincostflow

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type MCFEdge[Cap, Cost Integer] struct {
	From int
	To   int
	Cap  Cap
	Flow Cap
	Cost Cost
}

type _edge[Cap, Cost Integer] struct {
	to   int
	rev  int
	cap  Cap
	cost Cost
}

// Graph represents a graph for minimum cost flow.
type MCFGraph[Cap, Cost Integer] struct {
	n     int
	edges []MCFEdge[Cap, Cost]
}

// New creates a new minimum cost flow graph with n vertices.
func NewMCFGraph[Cap, Cost Integer](n int) *MCFGraph[Cap, Cost] {
	return &MCFGraph[Cap, Cost]{
		n: n,
	}
}

// AddEdge adds a directed edge from `from` to `to` with capacity `cap` and cost `cost`.
func (g *MCFGraph[Cap, Cost]) AddEdge(from, to int, cap Cap, cost Cost) int {
	if from < 0 || from >= g.n {
		panic("from is out of bounds")
	}
	if to < 0 || to >= g.n {
		panic("to is out of bounds")
	}
	if cap < 0 {
		panic("capacity cannot be negative")
	}
	if cost < 0 {
		panic("cost cannot be negative")
	}
	m := len(g.edges)
	g.edges = append(g.edges, MCFEdge[Cap, Cost]{
		From: from,
		To:   to,
		Cap:  cap,
		Flow: 0,
		Cost: cost,
	})
	return m
}

// GetEdge returns the i-th edge.
func (g *MCFGraph[Cap, Cost]) GetEdge(i int) MCFEdge[Cap, Cost] {
	if i < 0 || i >= len(g.edges) {
		panic("index out of bounds")
	}
	return g.edges[i]
}

// Edges returns all edges.
func (g *MCFGraph[Cap, Cost]) Edges() []MCFEdge[Cap, Cost] {
	res := make([]MCFEdge[Cap, Cost], len(g.edges))
	copy(res, g.edges)
	return res
}

func (g *MCFGraph[Cap, Cost]) Flow(s, t int) (Cap, Cost) {
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

func (g *MCFGraph[Cap, Cost]) FlowLimit(s, t int, flowLimit Cap) (Cap, Cost) {
	slope := g.SlopeLimit(s, t, flowLimit)
	last := slope[len(slope)-1]
	return last.Cap, last.Cost
}

func (g *MCFGraph[Cap, Cost]) Slope(s, t int) []struct {
	Cap  Cap
	Cost Cost
} {
	var maxCap Cap = 1
	for {
		if maxCap*2+1 > maxCap {
			maxCap = maxCap*2 + 1
		} else {
			break
		}
	}
	return g.SlopeLimit(s, t, maxCap)
}

func (g *MCFGraph[Cap, Cost]) SlopeLimit(s, t int, flowLimit Cap) []struct {
	Cap  Cap
	Cost Cost
} {
	if s < 0 || s >= g.n {
		panic("s is out of bounds")
	}
	if t < 0 || t >= g.n {
		panic("t is out of bounds")
	}
	if s == t {
		panic("s and t must be different")
	}

	m := len(g.edges)
	edgeIdx := make([]int, m)

	degree := make([]int, g.n)
	redgeIdx := make([]int, m)
	type elistItem struct {
		from int
		e    _edge[Cap, Cost]
	}
	elist := make([]elistItem, 0, 2*m)

	for i := 0; i < m; i++ {
		e := g.edges[i]
		edgeIdx[i] = degree[e.From]
		degree[e.From]++
		redgeIdx[i] = degree[e.To]
		degree[e.To]++
		elist = append(elist, elistItem{e.From, _edge[Cap, Cost]{e.To, -1, e.Cap - e.Flow, e.Cost}})
		elist = append(elist, elistItem{e.To, _edge[Cap, Cost]{e.From, -1, e.Flow, -e.Cost}})
	}

	start := make([]int, g.n+1)
	for i := 0; i < g.n; i++ {
		start[i+1] = start[i] + degree[i]
	}
	csrElist := make([]_edge[Cap, Cost], 2*m)
	counter := make([]int, g.n+1)
	copy(counter, start)

	for _, item := range elist {
		csrElist[counter[item.from]] = item.e
		counter[item.from]++
	}

	for i := 0; i < m; i++ {
		e := g.edges[i]
		edgeIdx[i] += start[e.From]
		redgeIdx[i] += start[e.To]
		csrElist[edgeIdx[i]].rev = redgeIdx[i]
		csrElist[redgeIdx[i]].rev = edgeIdx[i]
	}

	result := slopeImpl(g.n, start, csrElist, s, t, flowLimit)

	for i := 0; i < m; i++ {
		e := csrElist[edgeIdx[i]]
		g.edges[i].Flow = g.edges[i].Cap - e.cap
	}

	return result
}

func slopeImpl[Cap, Cost Integer](n int, start []int, elist []_edge[Cap, Cost], s, t int, flowLimit Cap) []struct {
	Cap  Cap
	Cost Cost
} {
	var maxCost Cost = 1
	for {
		if maxCost*2+1 > maxCost {
			maxCost = maxCost*2 + 1
		} else {
			break
		}
	}

	type dualDist struct {
		dual Cost
		dist Cost
	}
	dd := make([]dualDist, n)
	prevE := make([]int, n)
	vis := make([]bool, n)

	type Q struct {
		key Cost
		to  int
	}
	queMin := make([]int, 0)
	que := make([]Q, 0)

	pushHeap := func(q []Q) {
		child := len(q) - 1
		for child > 0 {
			parent := (child - 1) / 2
			if q[parent].key <= q[child].key {
				break
			}
			q[parent], q[child] = q[child], q[parent]
			child = parent
		}
	}

	popHeap := func(q []Q) {
		if len(q) == 0 {
			return
		}
		q[0], q[len(q)-1] = q[len(q)-1], q[0]
		parent := 0
		nHeap := len(q) - 1
		for {
			left := parent*2 + 1
			if left >= nHeap {
				break
			}
			right := left + 1
			minChild := left
			if right < nHeap && q[right].key < q[left].key {
				minChild = right
			}
			if q[parent].key <= q[minChild].key {
				break
			}
			q[parent], q[minChild] = q[minChild], q[parent]
			parent = minChild
		}
	}

	dualRef := func() bool {
		for i := 0; i < n; i++ {
			dd[i].dist = maxCost
			vis[i] = false
		}
		queMin = queMin[:0]
		que = que[:0]

		heapR := 0

		dd[s].dist = 0
		queMin = append(queMin, s)
		for len(queMin) > 0 || len(que) > 0 {
			var v int
			if len(queMin) > 0 {
				v = queMin[len(queMin)-1]
				queMin = queMin[:len(queMin)-1]
			} else {
				for heapR < len(que) {
					heapR++
					pushHeap(que[:heapR])
				}
				v = que[0].to
				popHeap(que[:heapR])
				que = que[:len(que)-1]
				heapR--
			}

			if vis[v] {
				continue
			}
			vis[v] = true
			if v == t {
				break
			}

			dualV := dd[v].dual
			distV := dd[v].dist
			for i := start[v]; i < start[v+1]; i++ {
				e := elist[i]
				if e.cap == 0 {
					continue
				}
				cost := e.cost - dd[e.to].dual + dualV
				if dd[e.to].dist-distV > cost {
					distTo := distV + cost
					dd[e.to].dist = distTo
					prevE[e.to] = e.rev
					if distTo == distV {
						queMin = append(queMin, e.to)
					} else {
						que = append(que, Q{distTo, e.to})
					}
				}
			}
		}

		if !vis[t] {
			return false
		}

		for v := 0; v < n; v++ {
			if !vis[v] {
				continue
			}
			dd[v].dual -= dd[t].dist - dd[v].dist
		}
		return true
	}

	flow := Cap(0)
	cost := Cost(0)
	var prevCostPerFlow Cost
	firstIteration := true
	result := []struct {
		Cap  Cap
		Cost Cost
	}{{0, 0}}

	for flow < flowLimit {
		if !dualRef() {
			break
		}
		c := flowLimit - flow
		for v := t; v != s; v = elist[prevE[v]].to {
			revCap := elist[elist[prevE[v]].rev].cap
			if c > revCap {
				c = revCap
			}
		}
		for v := t; v != s; v = elist[prevE[v]].to {
			elist[prevE[v]].cap += c
			elist[elist[prevE[v]].rev].cap -= c
		}
		d := -dd[s].dual
		flow += c
		cost += Cost(c) * d
		if !firstIteration && prevCostPerFlow == d {
			result = result[:len(result)-1]
		}
		result = append(result, struct {
			Cap  Cap
			Cost Cost
		}{flow, cost})
		prevCostPerFlow = d
		firstIteration = false
	}

	return result
}
