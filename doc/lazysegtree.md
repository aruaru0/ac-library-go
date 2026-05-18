# LazySegTree

モノイド $(S, \cdot: S \times S \to S, e \in S)$ と、$S$ から $S$ への写像の集合 $F$ であって、以下の条件を満たすようなものについて使用できるデータ構造（遅延評価セグメント木）です。

* $F$ は恒等写像 $\mathrm{id}$ を含む。つまり、任意の $x \in S$ に対し $\mathrm{id}(x) = x$ をみたす。
* $F$ は写像の合成について閉じている。つまり、任意の $f, g \in F$ に対し $f \circ g \in F$ である。
* 任意の $f \in F, x, y \in S$ に対し $f(x \cdot y) = f(x) \cdot f(y)$ をみたす。

長さ $n$ の $S$ の配列に対し、
* 区間の要素に一括で $F$ の要素 $f$ を作用（$x = f(x)$）
* 区間の要素の総積の取得

を $O(\log n)$ で行うことが出来ます。

## コンストラクタ

```go
func NewLazySegTree[S any, F any](n int, op func(S, S) S, e func() S, mapping func(F, S) S, composition func(F, F) F, id func() F) *LazySegTree[S, F]
func NewLazySegTreeFromSlice[S any, F any](v []S, op func(S, S) S, e func() S, mapping func(F, S) S, composition func(F, F) F, id func() F) *LazySegTree[S, F]
```

* `op`: $\cdot: S \times S \to S$ を計算する関数
* `e`: $e$ を返す関数
* `mapping`: $f(x)$ を返す関数
* `composition`: $f \circ g$ を返す関数
* `id`: $\mathrm{id}$ を返す関数

**制約**
* $0 \leq n \leq 10^8$

**計算量**
* $O(n)$

## Set

```go
func (st *LazySegTree[S, F]) Set(p int, x S)
```

$a[p] = x$ とします。

**制約**
* $0 \leq p < n$

**計算量**
* $O(\log n)$

## Get

```go
func (st *LazySegTree[S, F]) Get(p int) S
```

$a[p]$ を返します。

**制約**
* $0 \leq p < n$

**計算量**
* $O(\log n)$

## Prod

```go
func (st *LazySegTree[S, F]) Prod(l, r int) S
```

`op(a[l], ..., a[r - 1])` を返します。$l = r$ のときは `e()` を返します。

**制約**
* $0 \leq l \leq r \leq n$

**計算量**
* $O(\log n)$

## AllProd

```go
func (st *LazySegTree[S, F]) AllProd() S
```

`op(a[0], ..., a[n - 1])` を計算します。$n = 0$ のときは `e()` を返します。

**計算量**
* $O(1)$

## Apply

```go
func (st *LazySegTree[S, F]) Apply(p int, f F)
func (st *LazySegTree[S, F]) ApplyRange(l, r int, f F)
```

* `Apply`: $a[p] = f(a[p])$ を実行します。
* `ApplyRange`: $i = l \dots r-1$ について $a[i] = f(a[i])$ を実行します。

**制約**
* $0 \leq p < n$
* $0 \leq l \leq r \leq n$

**計算量**
* $O(\log n)$

## MaxRight

```go
func (st *LazySegTree[S, F]) MaxRight(l int, g func(S) bool) int
```

以下の条件を両方満たす $r$ を返します。

* $r = l$ もしくは `g(op(a[l], a[l + 1], ..., a[r - 1])) = true`
* $r = n$ もしくは `g(op(a[l], a[l + 1], ..., a[r])) = false`

**制約**
* `g(e()) = true`
* $0 \leq l \leq n$

**計算量**
* $O(\log n)$

## MinLeft

```go
func (st *LazySegTree[S, F]) MinLeft(r int, g func(S) bool) int
```

以下の条件を両方満たす $l$ を返します。

* $l = r$ もしくは `g(op(a[l], a[l + 1], ..., a[r - 1])) = true`
* $l = 0$ もしくは `g(op(a[l - 1], a[l], ..., a[r - 1])) = false`

**制約**
* `g(e()) = true`
* $0 \leq r \leq n$

**計算量**
* $O(\log n)$
