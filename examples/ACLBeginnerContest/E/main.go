package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aruaru0/ac-library-go/lazysegtree"
)

const MOD = 998244353
const inv9 = 443664157

type S struct {
	val   int64
	pow10 int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	op := func(a, b S) S {
		return S{
			val:   (a.val*b.pow10 + b.val) % MOD,
			pow10: (a.pow10 * b.pow10) % MOD,
		}
	}

	e := func() S {
		return S{val: 0, pow10: 1}
	}

	mapping := func(f int, x S) S {
		if f == 0 {
			return x
		}
		// (pow10 - 1) * inv9 * f
		num := (x.pow10 - 1 + MOD) % MOD
		val := num * inv9 % MOD * int64(f) % MOD
		return S{
			val:   val,
			pow10: x.pow10,
		}
	}

	composition := func(f, g int) int {
		if f == 0 {
			return g
		}
		return f
	}

	id := func() int {
		return 0
	}

	// 初期値はすべて1
	initV := make([]S, n)
	for i := 0; i < n; i++ {
		initV[i] = S{val: 1, pow10: 10}
	}

	st := lazysegtree.NewLazySegTreeFromSlice(initV, op, e, mapping, composition, id)

	for i := 0; i < q; i++ {
		var l, r, d int
		fmt.Fscan(reader, &l, &r, &d)
		l-- // 0-indexedにする
		st.ApplyRange(l, r, d)
		fmt.Fprintln(writer, st.AllProd().val)
	}
}
