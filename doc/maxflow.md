# MaxFlow

最大フロー問題 (MaxFlow) を解くライブラリです。

## コンストラクタ

```go
func NewMFGraph[Cap Integer](n int) *MFGraph[Cap]
```

$n$ 頂点 $0$ 辺のグラフを作る。`Cap` は容量の型（整数型）。

**制約**
* $0 \leq n \leq 10^8$
* `Cap` は組み込みの整数型

**計算量**
* $O(n)$

## AddEdge

```go
func (g *MFGraph[Cap]) AddEdge(from, to int, cap Cap) int
```

`from` から `to` へ最大容量 `cap`、流量 $0$ の辺を追加し、何番目に追加された辺かを返す。

**制約**
* $0 \leq \mathrm{from}, \mathrm{to} < n$
* $0 \leq \mathrm{cap}$

**計算量**
* ならし $O(1)$

## Flow / FlowLimit

```go
func (g *MFGraph[Cap]) Flow(s, t int) Cap
func (g *MFGraph[Cap]) FlowLimit(s, t int, flowLimit Cap) Cap
```

* `Flow`: 頂点 $s$ から $t$ へ流せる限り流し、流せた量を返す。
* `FlowLimit`: 頂点 $s$ から $t$ へ流量 `flowLimit` に達するまで流せる限り流し、流せた量を返す。
* 複数回呼ぶことも可能です。

**制約**
* $s \neq t$
* $0 \leq s, t < n$
* 返り値が `Cap` に収まる

**計算量**
$m$ を追加された辺数として
* $O(n^2 m)$
* $O((n + m) \sqrt{m})$ (辺の容量がすべて $1$ の時)

## MinCut

```go
func (g *MFGraph[Cap]) MinCut(s int) []bool
```

長さ $n$ のスライスを返す。$i$ 番目の要素には、頂点 $s$ から $i$ へ残余グラフで到達可能なとき、またその時のみ `true` を返す。`Flow(s, t)` を呼んだ後に呼ぶと、返り値は $s, t$ 間の mincut に対応します。

**計算量**
* $O(n + m)$

## GetEdge / Edges

```go
type MFEdge[Cap Integer] struct {
	From int
	To   int
	Cap  Cap
	Flow Cap
}

func (g *MFGraph[Cap]) GetEdge(i int) MFEdge[Cap]
func (g *MFGraph[Cap]) Edges() []MFMFEdge[Cap]
```

* 今の内部の辺の状態を返す
* 辺の順番は `AddEdge` で追加された順番と同一

**計算量**
* `GetEdge`: $O(1)$
* `Edges`: $O(m)$

## ChangeEdge

```go
func (g *MFGraph[Cap]) ChangeEdge(i int, newCap, newFlow Cap)
```

$i$ 番目に追加された辺の容量、流量を `newCap`, `newFlow` に変更する。他の辺の容量、流量は変更しない。

**制約**
* $0 \leq \mathrm{newFlow} \leq \mathrm{newCap}$

**計算量**
* $O(1)$

## 使用例

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/maxflow"
)

func main() {
	// 4頂点のグラフ (0: 始点, 3: 終点)
	g := maxflow.NewMFGraph[int](4)

	g.AddEdge(0, 1, 5)
	g.AddEdge(0, 2, 4)
	g.AddEdge(1, 2, 2)
	g.AddEdge(1, 3, 3)
	g.AddEdge(2, 3, 5)

	// 0 から 3 への最大流を計算
	fmt.Println(g.Flow(0, 3)) // 7

	// 各辺の流量情報を取得
	for i, e := range g.Edges() {
		fmt.Printf("Edge %d: %d -> %d (flow: %d / cap: %d)\n", i, e.From, e.To, e.Flow, e.Cap)
	}
}
```
