# Fenwick Tree

長さ $n$ の配列に対し、

* 要素の 1 点変更
* 区間の要素の総和

を $O(\log n)$ で求めることが出来るデータ構造です。

## コンストラクタ

```go
func NewFenwickTree[T Numeric](n int) *FenwickTree[T]
```

* 長さ $n$ の配列 $a_0, a_1, \cdots, a_{n-1}$ を作ります。初期値はすべて $0$ です。
* 型引数 `T` には、`int`, `uint`, `int64`, `uint64`, `float64` などの数値型 (`Numeric` 制約) を指定します。

**制約**
* $0 \leq n \leq 10^8$

**計算量**
* $O(n)$

## Add

```go
func (f *FenwickTree[T]) Add(p int, x T)
```

`a[p] += x` を行います。

**制約**
* $0 \leq p < n$

**計算量**
* $O(\log n)$

## Sum

```go
func (f *FenwickTree[T]) Sum(l, r int) T
```

`a[l] + a[l + 1] + ... + a[r - 1]` を返します。
`T` が符号付き・符号なし整数型の場合、答えがオーバーフローしたならば Go 言語の仕様に従ってラップアラウンドした値（$\bmod 2^{\mathrm{bit}}$ で等しい値）を返します。

**制約**
* $0 \leq l \leq r \leq n$

**計算量**
* $O(\log n)$

## 使用例

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/fenwicktree"
)

func main() {
	// 長さ 5 の Fenwick Tree を作成
	fw := fenwicktree.NewFenwickTree[int](5)

	// 値の加算
	fw.Add(0, 10)
	fw.Add(2, 20)
	fw.Add(4, 30)

	// 区間和の取得 (半開区間 [l, r))
	fmt.Println(fw.Sum(0, 3)) // a[0] + a[1] + a[2] = 10 + 0 + 20 = 30
	fmt.Println(fw.Sum(1, 5)) // a[1] + a[2] + a[3] + a[4] = 0 + 20 + 0 + 30 = 50
}
```
