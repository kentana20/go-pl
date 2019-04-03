package main

import (
	"fmt"
)

func main() {
	a := []string{"a", "b", "b", "c", "c", "d", "d"}
	a2 := dup(a)
	fmt.Println(a2)
}

func dup(s []string) []string {
	var ls []string

	ls = append(ls, s[0])
	for i := 0; i < len(s)-1; i++ {
		if s[i] != s[i+1] {
			ls = append(ls, s[i+1])
		}
	}
	return ls
}
