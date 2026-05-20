# TwoSAT

2-SATを解きます。
変数 $x_0, x_1, \cdots, x_{N - 1}$ に関して、

* $(x_i = f) \lor (x_j = g)$

というクローズ（節）を足し、これをすべて満たす変数の割当があるかを判定・構築します。

## コンストラクタ

```go
func NewTwoSAT(n int) *TwoSAT
```

$n$ 変数の2-SATを作ります。

**制約**
* $0 \leq n \leq 10^8$

**計算量**
* $O(n)$

## AddClause

```go
func (ts *TwoSAT) AddClause(i int, f bool, j int, g bool)
```

$(x_i = f) \lor (x_j = g)$ というクローズを足します。

**制約**
* $0 \leq i < n$
* $0 \leq j < n$

**計算量**
* ならし $O(1)$

## Satisfiable

```go
func (ts *TwoSAT) Satisfiable() bool
```

条件を満たす割当が存在するかどうかを判定します。割当が存在するならば `true`、そうでないなら `false` を返します。

**制約**
* 複数回呼ぶことも可能。

**計算量**
足した制約の個数を $m$ として
* $O(n + m)$

## Answer

```go
func (ts *TwoSAT) Answer() []bool
```

最後に呼んだ `Satisfiable` の、クローズを満たす割当を返します。
`Satisfiable` を呼ぶ前や、`Satisfiable` で割当が存在しなかったときにこの関数を呼ぶと、中身が未定義の長さ $n$ のスライスを返します。

**計算量**
* $O(n)$

## 使用例

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/twosat"
)

func main() {
	// 3変数の 2-SAT (x_0, x_1, x_2)
	ts := twosat.NewTwoSAT(3)

	// (x_0 または x_1)
	ts.AddClause(0, true, 1, true)
	// (not x_0 または not x_2)
	ts.AddClause(0, false, 2, false)
	// (x_1 または x_2)
	ts.AddClause(1, true, 2, true)

	if ts.Satisfiable() {
		fmt.Println("Satisfiable")
		fmt.Println(ts.Answer()) // [true true false] などの解 (満たす割り当ての1つ)
	} else {
		fmt.Println("Unsatisfiable")
	}
}
```
