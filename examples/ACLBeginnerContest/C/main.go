package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aruaru0/ac-library-go/dsu"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	d := dsu.NewDSU(n)
	for i := 0; i < m; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return
		}
		u--
		v--
		d.Merge(u, v)
	}

	// 連結成分の数 - 1
	ans := len(d.Groups()) - 1
	fmt.Fprintln(writer, ans)
}
