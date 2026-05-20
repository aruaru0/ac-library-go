package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aruaru0/ac-library-go/convolution"
	"github.com/aruaru0/ac-library-go/modint"
)

const MOD = 998244353

var fact []int64
var invFact []int64
var doubleFact []int64
var invPowerOf2 []int64

func initFact(maxN int) {
	fact = make([]int64, maxN+1)
	invFact = make([]int64, maxN+1)
	doubleFact = make([]int64, maxN+1)
	invPowerOf2 = make([]int64, maxN+1)

	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}

	invFact[maxN] = powMod(fact[maxN], MOD-2, MOD)
	for i := maxN - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % MOD
	}

	doubleFact[0] = 1
	for i := 1; i <= maxN; i++ {
		doubleFact[i] = doubleFact[i-1] * int64(2*i-1) % MOD
	}

	powerOf2 := make([]int64, maxN+1)
	powerOf2[0] = 1
	for i := 1; i <= maxN; i++ {
		powerOf2[i] = powerOf2[i-1] * 2 % MOD
	}
	invPowerOf2[maxN] = powMod(powerOf2[maxN], MOD-2, MOD)
	for i := maxN - 1; i >= 0; i-- {
		invPowerOf2[i] = invPowerOf2[i+1] * 2 % MOD
	}
}

func powMod(x, n, m int64) int64 {
	r := int64(1)
	x %= m
	for n > 0 {
		if n&1 == 1 {
			r = r * x % m
		}
		x = x * x % m
		n >>= 1
	}
	return r
}

type Mint = modint.StaticModint[modint.Mod998244353]

func makePolynomial(c int) []Mint {
	size := c/2 + 1
	poly := make([]Mint, size)
	for k := 0; k < size; k++ {
		// g_c(k) = c! / ((c-2k)! * k! * 2^k)
		val := fact[c] * invFact[c-2*k] % MOD * invFact[k] % MOD * invPowerOf2[k] % MOD
		poly[k] = modint.NewStatic[modint.Mod998244353](val)
	}
	return poly
}

func solvePolys(polys [][]Mint) []Mint {
	if len(polys) == 0 {
		return []Mint{modint.NewStatic[modint.Mod998244353](1)}
	}
	if len(polys) == 1 {
		return polys[0]
	}
	mid := len(polys) / 2
	left := solvePolys(polys[:mid])
	right := solvePolys(polys[mid:])
	return convolution.Convolution[modint.Mod998244353](left, right)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	h := make([]int, 2*n)
	counts := make(map[int]int)
	for i := 0; i < 2*n; i++ {
		fmt.Fscan(reader, &h[i])
		counts[h[i]]++
	}

	initFact(2 * n)

	polys := [][]Mint{}
	for _, count := range counts {
		if count >= 2 {
			polys = append(polys, makePolynomial(count))
		}
	}

	g := solvePolys(polys)

	ans := int64(0)
	for k := 0; k < len(g); k++ {
		if k > n {
			break
		}
		coeff := int64(g[k].Val())
		term := coeff * doubleFact[n-k] % MOD
		if k%2 == 1 {
			ans = (ans - term + MOD) % MOD
		} else {
			ans = (ans + term) % MOD
		}
	}

	fmt.Fprintln(writer, ans)
}
