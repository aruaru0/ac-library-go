# SCC

有向グラフを強連結成分分解 (Strongly Connected Components) します。

## コンストラクタ

```go
func NewSCCGraph(n int) *SCCGraph
```

$n$ 頂点 $0$ 辺の有向グラフを作ります。

**制約**
* $0 \leq n \leq 10^8$

**計算量**
* $O(n)$

## AddEdge

```go
func (g *SCCGraph) AddEdge(from, to int)
```

頂点 `from` から頂点 `to` へ有向辺を足します。

**制約**
* $0 \leq \mathrm{from} < n$
* $0 \leq \mathrm{to} < n$

**計算量**
* ならし $O(1)$

## SCC

```go
func (g *SCCGraph) SCC() [][]int
```

以下の条件を満たすような、「頂点のリスト」のリストを返します。

* 全ての頂点がちょうど 1 つずつ、どれかのリストに含まれます。
* 内側のリストと強連結成分が一対一に対応します。リスト内での頂点の順序は未定義です。
* リストはトポロジカルソートされています。異なる強連結成分の頂点 $u, v$ について、$u$ から $v$ に到達できる時、$u$ の属するリストは $v$ の属するリストよりも前です。

**計算量**
追加した辺の本数を $m$ として
* $O(n + m)$

## 使用例

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/scc"
)

func main() {
	// 4頂点のグラフ
	g := scc.NewSCCGraph(4)

	// サイクル 0 -> 1 -> 2 -> 0 を作成
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	// 2 から 3 へ伸ばす (3は単体で連結成分になる)
	g.AddEdge(2, 3)

	sccList := g.SCC()
	fmt.Println(sccList) // [[0 1 2] [3]] (トポロジカル順にグループ化されて出力)
}
```
