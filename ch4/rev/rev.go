package main

import "fmt"

func main() {
	a := [...]int{0, 1, 2, 3, 4, 5}

	reverse(a[:])
	fmt.Println(a) // [5, 4, 3, 2, 1, 0]

	s := []int{0, 1, 2, 3, 4, 5} // スライスリテラル
	reverse(s[:2])               // [0, 1] -> [1, 0]
	reverse(s[2:])               // [2, 3, 4, 5] -> [5, 4, 3, 2]
	reverse(s)                   // [2, 3, 4, 5, 0, 1]
	fmt.Println(s)               // [2, 3, 4, 5, 0, 1]
}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
