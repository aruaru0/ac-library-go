# Modint

自動で mod を取る構造体です。Go 言語では演算子オーバーロードがないため、メソッドチェーンを用いて演算を行います。

多くの問題では `Modint998244353`、`Modint1000000007`、または `DynamicModint` のどれかを使えば十分です。

## 使用例 (Static Modint)

```go
package main

import (
	"fmt"
	"ac-library-go/modint"
)

func main() {
	// sum of array (mod 998244353)
	var n int
	fmt.Scan(&n)
	sum := modint.NewModint998244353(0)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Scan(&x)
		sum = sum.Add(modint.NewModint998244353(x))
	}
	fmt.Println(sum.Val())
}
```

## 使用例 (Dynamic Modint)

実行時に mod が決まる場合は `DynamicModint` を使用します。

```go
package main

import (
	"fmt"
	"ac-library-go/modint"
)

func main() {
	var n int
	var mod uint32
	fmt.Scan(&n, &mod)
	sum := modint.NewDynamic(0, mod)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Scan(&x)
		sum = sum.Add(modint.NewDynamic(x, mod))
	}
	fmt.Println(sum.Val())
}
```

## コンストラクタ

```go
func NewStatic[M Modulus](v int64) StaticModint[M]
func NewDynamic(v int64, mod uint32) DynamicModint
```

整数 `v` の mod を取って `Modint` に格納します。

## メソッド

* `Val() uint32` : 格納されている値を返します。
* `Add(rhs) Modint` : 加算 (`+`) の結果を返します。
* `Sub(rhs) Modint` : 減算 (`-`) の結果を返します。
* `Mul(rhs) Modint` : 乗算 (`*`) の結果を返します。
* `Div(rhs) Modint` : 除算 (`/`) の結果を返します。（逆元を掛けるのと同じです）
* `Pow(n int64) Modint` : $x^n$ を返します。
* `Inv() Modint` : 逆元を返します。

**計算量**
* `Add`, `Sub`, `Mul`: $O(1)$
* `Div`, `Inv`, `Pow`: $O(\log \mathrm{mod})$
