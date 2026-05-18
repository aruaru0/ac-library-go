# String

文字列アルゴリズム詰め合わせです。

## SuffixArray

```go
func SuffixArrayString(s string) []int
func SuffixArray[T Integer](s []T) []int
func SuffixArrayUpperBound(s []int, upper int) []int
```

長さ $n$ の文字列 `s` の Suffix Array として、長さ $n$ の配列を返す。
Suffix Array `sa` は $(0, 1, \dots, n - 1)$ の順列であって、各 $i = 0,1, \cdots ,n-2$ について `s[sa[i]..n) < s[sa[i+1]..n)` を満たすもの。

**制約**
* $0 \leq n \leq 10^8$
* `SuffixArrayUpperBound`: $0 \leq \mathrm{upper} \leq 10^8$ かつ $s$ のすべての要素 $x$ について $0 \leq x \leq \mathrm{upper}$

**計算量**
* `SuffixArrayString`, `SuffixArrayUpperBound`: $O(n + \mathrm{upper})$
* `SuffixArray`: $O(n \log n)$

## LcpArray

```go
func LcpArrayString(s string, sa []int) []int
func LcpArray[T Integer](s []T, sa []int) []int
```

長さ $n$ の文字列 `s` の LCP Array として、長さ $n-1$ の配列を返す。$i$ 番目の要素は `s[sa[i]..n)` と `s[sa[i+1]..n)` の LCP (Longest Common Prefix) の長さ。

**制約**
* `sa` は `s` の Suffix Array
* $1 \leq n \leq 10^8$

**計算量**
* $O(n)$

## ZAlgorithm

```go
func ZAlgorithmString(s string) []int
func ZAlgorithm[T Integer](s []T) []int
```

入力の長さを $n$ として、長さ $n$ の配列を返す。
$i$ 番目の要素は `s[0..n)` と `s[i..n)` の LCP (Longest Common Prefix) の長さ。

**制約**
* $0 \leq n \leq 10^8$

**計算量**
* $O(n)$
