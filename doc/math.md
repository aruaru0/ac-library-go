# Math

数論的アルゴリズム詰め合わせです。

## PowMod

```go
func PowMod(x, n int64, m int) int64
```

$x^n \bmod m$ を返します。

**制約**
* $0 \leq n$
* $1 \leq m$

**計算量**
* $O(\log n)$

## InvMod

```go
func InvMod(x, m int64) int64
```

$xy \equiv 1 \pmod m$ なる $y$ のうち、$0 \leq y < m$ を満たすものを返します。

**制約**
* $\gcd(x, m) = 1$
* $1 \leq m$

**計算量**
* $O(\log m)$

## CRT

```go
func CRT(r, m []int64) (int64, int64)
```

同じ長さのスライス `r`, `m` を渡します。このスライスの長さを $n$ とした時、

$$x \equiv r[i] \pmod{m[i]}, \forall i \in \lbrace 0,1,\cdots, n - 1 \rbrace$$

を解きます。答えは(存在するならば) $y, z \ (0 \leq y < z = \mathrm{lcm}(m[i]))$ を用いて $x \equiv y \pmod z$ の形で書けることが知られており、この $(y, z)$ を返します。答えがない場合は $(0, 0)$ を返します。$n=0$ の時は $(0, 1)$ を返します。

**制約**
* $\mathrm{len}(r) = \mathrm{len}(m)$
* $1 \leq m[i]$
* $\mathrm{lcm}(m[i])$ が `int64` に収まる。

**計算量**
* $O(n \log{\mathrm{lcm}(m[i])})$

## FloorSum

```go
func FloorSum(n, m, a, b int64) int64
```

$$\sum_{i = 0}^{n - 1} \left\lfloor \frac{a \times i + b}{m} \right\rfloor$$

を返します。答えがオーバーフローした場合は、Go 言語の仕様に従ってラップアラウンドします。

**制約**
* $0 \leq n < 2^{32}$
* $1 \leq m < 2^{32}$

**計算量**
* $O(\log m)$
