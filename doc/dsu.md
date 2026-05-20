# DSU (Disjoint Set Union)

無向グラフに対して、

* 辺の追加
* 2 頂点が連結かの判定

をならし $O(\alpha(n))$ 時間で処理することが出来ます。

また、内部的に各連結成分ごとに代表となる頂点を 1 つ持っています。辺の追加により連結成分がマージされる時、新たな代表元は元の連結成分の代表元のうちどちらかになります。

## コンストラクタ

```go
func NewDSU(n int) *DSU
```

* $n$ 頂点 $0$ 辺の無向グラフを作ります。

**制約**
* $0 \leq n \leq 10^8$

**計算量**
* $O(n)$

## Merge

```go
func (d *DSU) Merge(a, b int) int
```

辺 $(a, b)$ を足します。
$a, b$ が連結だった場合はその代表元、非連結だった場合は新たな代表元を返します。

**制約**
* $0 \leq a < n$
* $0 \leq b < n$

**計算量**
* ならし $O(\alpha(n))$

## Same

```go
func (d *DSU) Same(a, b int) bool
```

頂点 $a, b$ が連結かどうかを返します。

**制約**
* $0 \leq a < n$
* $0 \leq b < n$

**計算量**
* ならし $O(\alpha(n))$

## Leader

```go
func (d *DSU) Leader(a int) int
```

頂点 $a$ の属する連結成分の代表元を返します。

**制約**
* $0 \leq a < n$

**計算量**
* ならし $O(\alpha(n))$

## Size

```go
func (d *DSU) Size(a int) int
```

頂点 $a$ の属する連結成分のサイズを返します。

**制約**
* $0 \leq a < n$

**計算量**
* ならし $O(\alpha(n))$

## Groups

```go
func (d *DSU) Groups() [][]int
```

グラフを連結成分に分け、その情報を返します。
返り値は「『一つの連結成分の頂点番号のリスト』のリスト」です。
(内側外側限らず) スライス内でどの順番で頂点が格納されているかは未定義です。

**計算量**
* $O(n)$

## 使用例

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/dsu"
)

func main() {
	// 5頂点のDSUを作成 (0, 1, 2, 3, 4)
	d := dsu.NewDSU(5)

	d.Merge(0, 1)
	d.Merge(2, 3)

	fmt.Println(d.Same(0, 1)) // true
	fmt.Println(d.Same(0, 2)) // false

	// 1と2を結合することで、{0, 1} と {2, 3} が結合して {0, 1, 2, 3} になる
	d.Merge(1, 2)
	fmt.Println(d.Same(0, 2)) // true
	fmt.Println(d.Size(0))   // 4 (0の属するグループ의サイズ)

	// グループ一覧の取得
	fmt.Println(d.Groups()) // [[0 1 2 3] [4]] (出力順は未定義)
}
```
