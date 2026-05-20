package main

import (
	"fmt"
	"strings"
)

func main() {
	var k int
	if _, err := fmt.Scan(&k); err != nil {
		return
	}
	fmt.Println(strings.Repeat("ACL", k))
}
