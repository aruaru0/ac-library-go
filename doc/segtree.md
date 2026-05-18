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
