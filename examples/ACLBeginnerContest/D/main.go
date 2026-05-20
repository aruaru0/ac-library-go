package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aruaru0/ac-library-go/segtree"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	maxVal := 300000
	op := func(x, y int) int {
		return max(x, y)
	}
	e := func() int {
		return 0
	}

	st := segtree.NewSegTree(maxVal+1, op, e)

	for _, v := range a {
		l := max(0, v-k)
		r := min(maxVal, v+k) + 1
		best := st.Prod(l, r)
		st.Set(v, max(st.Get(v), best+1))
	}

	fmt.Fprintln(writer, st.AllProd())
}
