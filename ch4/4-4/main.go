package main

import "fmt"

func main() {
	s := []int{0, 1, 2, 3, 4, 5} // スライスリテラル
	r := rotate(s, 3)
	fmt.Println(r) // [3, 4, 5, 0, 1, 2]
}

func rotate(s []int, i int) []int {
	s = append(s, s[:i]...)
	fmt.Println(s) // [0, 1, 2, 3, 4, 5, 0, 1, 2]

	// i番目の要素を起点にスライスして返す
	return s[i:]
}
