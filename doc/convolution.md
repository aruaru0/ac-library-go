# Convolution

畳み込み (Convolution) を計算します。

## Convolution

```go
func Convolution[M modint.Modulus](a, b []modint.StaticModint[M]) []modint.StaticModint[M]
```

要素数 $N$ の配列 `a` と要素数 $M$ の配列 `b` を畳み込んだ結果（長さ $N+M-1$ の配列）を返します。
計算には Number Theoretic Transform (NTT) が使用されます。

**制約**
* $\mathrm{mod}$ は素数
* $2^c | (\mathrm{mod} - 1)$ かつ $N + M - 1 \le 2^c$ を満たす $c$ が存在する

**計算量**
* $O((N+M) \log(N+M))$

**使用例**
```go
import "ac-library-go/modint"

func main() {
    a := []modint.Modint998244353{ /* ... */ }
    b := []modint.Modint998244353{ /* ... */ }
    c := convolution.Convolution(a, b)
}
```

## ConvolutionLL

```go
func ConvolutionLL(a, b []int64) []int64
```

整数配列の畳み込みを計算します。
内部的に 3 つの素数を用いた NTT と Garner のアルゴリズムを用いて計算されます。

**制約**
* 長さ $N + M - 1 \le 2^{24}$
* 答えの各要素が $2^{63}-1$ 以下の整数に収まること

**計算量**
* $O((N+M) \log(N+M))$
