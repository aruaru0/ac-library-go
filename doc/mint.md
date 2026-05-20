# Mint (Modular Arithmetic)

法（mod）をパラメータ化したモジュラ演算のヘルパークラスです。
通常の `int` を引数に取り、演算結果を `int` で返すため、演算子オーバーロードがない Go においても簡潔に数式を記述することができます。
また、組み合わせ（$nCr$）の高速な計算用のメモ化テーブルや、行列累乗の機能も内包しています。

## コンストラクタ

```go
func NewModInt(mod int) *ModInt
```

* コンテクストの `ModInt` インスタンスを作ります。

**制約**
* $1 \leq mod \leq 2 \times 10^9$

**計算量**
* $O(1)$

## Add

```go
func (m *ModInt) Add(a, b int) int
```

* $(a + b) \pmod{mod}$ を返します。

## Sub

```go
func (m *ModInt) Sub(a, b int) int
```

* $(a - b) \pmod{mod}$ を返します。

## Mul

```go
func (m *ModInt) Mul(a, b int) int
```

* $(a \times b) \pmod{mod}$ を返します。内部で一時的に `int64` にキャストされるため、オーバーフローを防ぐことができます。

## Div

```go
func (m *ModInt) Div(a, b int) int
```

* $(a / b) \pmod{mod}$ （すなわち $a \times b^{-1} \pmod{mod}$）を返します。

## Pow

```go
func (m *ModInt) Pow(p, n int) int
```

* $(p^n) \pmod{mod}$ を返します。

**計算量**
* $O(\log n)$

## ModInv

```go
func (m *ModInt) ModInv(a int) int
```

* 拡張ユークリッドの互除法を用いて、$a$ の $\pmod{mod}$ における逆元 $a^{-1}$ を計算して返します。

**計算量**
* $O(\log mod)$

## InitComb

```go
func (m *ModInt) InitComb(limit int)
```

* 組み合わせ計算（$nCr$）に必要な階乗およびその逆元のテーブルを $0$ から `limit` まで $O(limit)$ で一括計算・初期化します。

**計算量**
* $O(limit)$

## Fact

```go
func (m *ModInt) Fact(n int) int
```

* $n! \pmod{mod}$ を返します。テーブルサイズが不足している場合は、自動的に拡張されます。

## IFact

```go
func (m *ModInt) IFact(n int) int
```

* $(n!)^{-1} \pmod{mod}$ を返します。テーブルサイズが不足している場合は、自動的に拡張されます。

## NCr

```go
func (m *ModInt) NCr(n, r int) int
```

* $_n\text{C}_r \pmod{mod}$ を返します。

**計算量**
* $O(1)$ （事前に `InitComb` を行っている場合）

## PowModMatrix

```go
func (m *ModInt) PowModMatrix(A [][]int, p int) [][]int
```

* 行列 $A$ の $p$ 乗 $A^p \pmod{mod}$ を計算して返します。

**計算量**
* $O(N^3 \log p)$ （$N$ は行列のサイズ）

## MulMod

```go
func (m *ModInt) MulMod(A, B [][]int) [][]int
```

* 行列の積 $(A \times B) \pmod{mod}$ を返します。

## 使用例

```go
package main

import (
	"fmt"
	"github.com/aruaru0/ac-library-go/mint"
)

func main() {
	m := mint.NewModInt(998244353)

	// 基本演算
	a := 1000000
	b := 2000000
	c := m.Add(a, b)
	d := m.Mul(a, b)
	fmt.Println(c) // 3000000
	fmt.Println(d) // 2000000000000 % 998244353 = 4310574

	// 組み合わせ計算 (nCr)
	// 事前に10万までの階乗テーブルを初期化
	m.InitComb(100000)
	fmt.Println(m.NCr(100000, 50000)) // 100000 C 50000 % 998244353
}
```
