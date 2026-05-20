# SegTree

モノイド $(S, \cdot: S \times S \to S, e \in S)$ に対し使用できるデータ構造（セグメント木）です。

長さ $n$ の $S$ の配列に対し、
* 要素の $1$ 点変更
* 区間の要素の総積の取得

を $O(\log n)$ で行うことが出来ます。

## コンストラクタ

```go
func NewSegTree[S any](n int, op func(S, S) S, e func() S) *SegTree[S]
func NewSegTreeFromSlice[S any](v []S, op func(S, S) S, e func() S) *SegTree[S]
```

* `New` : 長さ $n$ の数列 $a$ を作ります。初期値は全部 `e()` です。
* `NewFromSlice` : 長さ $n = \mathrm{len}(v)$ の数列 $a$ を作ります。$v$ の内容が初期値となります。

**制約**
* $0 \leq n \leq 10^8$

**計算量**
* $O(n)$

## Set

```go
func (st *SegTree[S]) Set(p int, x S)
```

$a[p]$ に $x$ を代入します。

**制約**
* $0 \leq p < n$

**計算量**
* $O(\log n)$

## Get

```go
func (st *SegTree[S]) Get(p int) S
```

$a[p]$ を返します。

**制約**
* $0 \leq p < n$

**計算量**
* $O(1)$

## Prod

```go
func (st *SegTree[S]) Prod(l, r int) S
```

`op(a[l], ..., a[r - 1])` を返します。$l = r$ のときは `e()` を返します。

**制約**
* $0 \leq l \leq r \leq n$

**計算量**
* $O(\log n)$

## AllProd

```go
func (st *SegTree[S]) AllProd() S
```

`op(a[0], ..., a[n - 1])` を計算します。$n = 0$ のときは `e()` を返します。

**計算量**
* $O(1)$

## MaxRight

```go
func (st *SegTree[S]) MaxRight(l int, f func(S) bool) int
```

以下の条件を両方満たす $r$ を返します。

* $r = l$ もしくは `f(op(a[l], a[l + 1], ..., a[r - 1])) = true`
* $r = n$ もしくは `f(op(a[l], a[l + 1], ..., a[r])) = false`

`f` が単調だとすれば、条件を満たす最大の $r$ となります。

**制約**
* `f(e()) = true`
* $0 \leq l \leq n$

**計算量**
* $O(\log n)$

## MinLeft

```go
func (st *SegTree[S]) MinLeft(r int, f func(S) bool) int
```

以下の条件を両方満たす $l$ を返します。

* $l = r$ もしくは `f(op(a[l], a[l + 1], ..., a[r - 1])) = true`
* $l = 0$ もしくは `f(op(a[l - 1], a[l], ..., a[r - 1])) = false`

`f` が単調だとすれば、条件を満たす最小の $l$ となります。

**制約**
* `f(e()) = true`
* $0 \leq r \leq n$

**計算量**
* $O(\log n)$

## 使用例

### サンプル 1: 基本的な型 (`int`) を用いた区間最小値 (RMQ)

```go
package main

import (
	"fmt"
	"math"
	"github.com/aruaru0/ac-library-go/segtree"
)

func main() {
	op := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	e := func() int {
		return math.MaxInt32
	}

	v := []int{5, 8, 2, 10, 4}
	// 初期配列からセグメント木を作成
	st := segtree.NewSegTreeFromSlice(v, op, e)

	fmt.Println(st.Prod(0, 3)) // [5, 8, 2] の最小値 -> 2

	// 2番目の要素(値2)を 9 に変更
	st.Set(2, 9)
	fmt.Println(st.Prod(0, 3)) // [5, 8, 9] の最小値 -> 5
	fmt.Println(st.AllProd())  // 全体の最小値 -> 4
}
```

### サンプル 2: 構造体を用いた「区間最大値とその個数」の取得

```go
package main

import (
	"fmt"
	"math"
	"github.com/aruaru0/ac-library-go/segtree"
)

// セグメント木の各ノードで管理する構造体
type Node struct {
	Max   int // 区間の最大値
	Count int // 最大値の出現回数
}

func main() {
	// 2つのノードをマージする関数 op
	op := func(a, b Node) Node {
		if a.Max > b.Max {
			return a
		}
		if a.Max < b.Max {
			return b
		}
		// 最大値が等しい場合は、個数を合算する
		return Node{Max: a.Max, Count: a.Count + b.Count}
	}

	// 単位元 e (初期値・空区間の値)
	e := func() Node {
		return Node{Max: math.MinInt32, Count: 0}
	}

	// 初期データ: [3, 1, 3, 2]
	v := []Node{
		{Max: 3, Count: 1},
		{Max: 1, Count: 1},
		{Max: 3, Count: 1},
		{Max: 2, Count: 1},
	}

	st := segtree.NewSegTreeFromSlice(v, op, e)

	// 全体 [0, 4) の最大値とその個数を取得
	res1 := st.AllProd()
	fmt.Printf("Max: %d, Count: %d\n", res1.Max, res1.Count) // Max: 3, Count: 2

	// インデックス 1 (値 1) を 値 4 に更新
	st.Set(1, Node{Max: 4, Count: 1})

	// 再度全体の最大値を取得
	res2 := st.AllProd()
	fmt.Printf("Max: %d, Count: %d\n", res2.Max, res2.Count) // Max: 4, Count: 1
}
```
