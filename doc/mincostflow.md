# MinCostFlow

最小費用流問題 (Minimum-cost flow problem) を扱うライブラリです。

## コンストラクタ

```go
func NewMCFGraph[Cap, Cost Integer](n int) *MCFGraph[Cap, Cost]
```

$n$ 頂点 $0$ 辺のグラフを作る。`Cap` は容量の型、`Cost` はコストの型（共に整数型）。

**制約**
* $0 \leq n \leq 10^8$
* `Cap, Cost` は組み込みの整数型

**計算量**
* $O(n)$

## AddEdge

```go
func (g *MCFGraph[Cap, Cost]) AddEdge(from, to int, cap Cap, cost Cost) int
```

`from` から `to` へ最大容量 `cap`, コスト `cost` の辺を追加する。何番目に追加された辺かを返す。

**制約**
* $0 \leq \mathrm{from}, \mathrm{to} < n$
* $0 \leq \mathrm{cap}, \mathrm{cost}$

**計算量**
* ならし $O(1)$

## Flow / FlowLimit

```go
func (g *MCFGraph[Cap, Cost]) Flow(s, t int) (Cap, Cost)
func (g *MCFGraph[Cap, Cost]) FlowLimit(s, t int, flowLimit Cap) (Cap, Cost)
```

$s$ から $t$ へ流せるだけ流し、その流量とコストを返す。

* `Flow`: $s$ から $t$ へ流せるだけ流す
* `FlowLimit`: $s$ から $t$ へ流量 `flowLimit` まで流せるだけ流す

**計算量**
* `Slope` と同じ

## Slope / SlopeLimit

```go
func (g *MCFGraph[Cap, Cost]) Slope(s, t int) []struct{ Cap Cap; Cost Cost }
func (g *MCFGraph[Cap, Cost]) SlopeLimit(s, t int, flowLimit Cap) []struct{ Cap Cap; Cost Cost }
```

返り値に流量とコストの関係の折れ線が入る。全ての $x$ について、流量 $x$ の時の最小コストを $g(x)$ とすると、$(x, g(x))$ は返り値を折れ線として見たものに含まれる。

* 返り値の最初の要素は $(0, 0)$
* 返り値の `Cap` は狭義単調増加、`Cost` は広義単調増加
* 3点が同一線上にあることはない

**制約**
* $s \neq t$
* $0 \leq s, t < n$
* `Slope` や `Flow` を合わせて複数回呼んだときの挙動は未定義
* $s$ から $t$ へ流したフローの流量が `Cap` に収まる。
* 流したフローのコストの総和が `Cost` に収まる。

**計算量**
$F$ を流量、$m$ を追加した辺の本数として
* $O(F (n + m) \log (n + m))$

## GetEdge / Edges

```go
type MCFEdge[Cap, Cost Integer] struct {
	From int
	To   int
	Cap  Cap
	Flow Cap
	Cost Cost
}

func (g *MCFGraph[Cap, Cost]) GetEdge(i int) MCFEdge[Cap, Cost]
func (g *MCFGraph[Cap, Cost]) Edges() []MCFMCFEdge[Cap, Cost]
```

* 今の内部の辺の状態を返す
* 辺の順番は `AddEdge` で追加された順番と同一

**計算量**
* `GetEdge`: $O(1)$
* `Edges`: $O(m)$

## 使用例

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/mincostflow"
)

func main() {
	// 4頂点のグラフ (0: 始点, 3: 終点)
	g := mincostflow.NewMCFGraph[int, int](4)

	g.AddEdge(0, 1, 2, 1) // from, to, cap, cost
	g.AddEdge(0, 2, 1, 2)
	g.AddEdge(1, 2, 1, 1)
	g.AddEdge(1, 3, 1, 3)
	g.AddEdge(2, 3, 2, 1)

	// 0 から 3 へ流量 2 を流す最小費用
	flow, cost := g.FlowLimit(0, 3, 2)
	fmt.Printf("Flow: %d, Cost: %d\n", flow, cost) // Flow: 2, Cost: 6
}
```
