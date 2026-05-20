# Modint

自動で mod を取る構造体です。Go 言語では演算子オーバーロードがないため、メソッドチェーンを用いて演算を行います。

多くの問題では `Modint998244353`、`Modint1000000007`、または `DynamicModint` のどれかを使えば十分です。

## 使用例 (Static Modint)

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/modint"
)

func main() {
	a := modint.NewModint998244353(1000000000)
	b := modint.NewModint998244353(2000000000)

	// (10^9 + 2*10^9) % 998244353
	fmt.Println(a.Add(b)) // 3511294

	// (10^9 * 2*10^9) % 998244353
	fmt.Println(a.Mul(b)) // 1599388

	// 累乗 a^10 % 998244353
	fmt.Println(a.Pow(10))

	// 逆元 a^-1 % 998244353
	fmt.Println(a.Inv())
}
```

## 使用例 (Dynamic Modint)

実行時に mod が決まる場合は `DynamicModint` を使用します。

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/modint"
)

func main() {
	var mod uint32 = 1000000007
	a := modint.NewDynamic(1000000000, mod)
	b := modint.NewDynamic(2000000000, mod)

	// (10^9 + 2*10^9) % 1000000007
	fmt.Println(a.Add(b)) // 999999993
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
