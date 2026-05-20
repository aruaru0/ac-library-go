package main

import (
	"fmt"
)

func main() {
	var a, b, c, d int64
	if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
		return
	}
	if a <= d && c <= b {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
