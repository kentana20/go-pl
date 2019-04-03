package main

import "fmt"

func main() {
	a := [...]int{0, 1, 2, 3, 4, 5}

	reverse(a[:])
	fmt.Println(a) // [5, 4, 3, 2, 1, 0]

	reverse2(&a)
	fmt.Println(a)
}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// 要素数を指定すると動作するけど、これであってるのだろうか
func reverse2(s *[6]int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
