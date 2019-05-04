package main

import (
	"fmt"

	"github.com/kentana20/go-pl/ch6/intset"
)

func main() {
	var x, y intset.IntSet
	x.AddAll(1, 3, 5)
	fmt.Println(x.String())

	y.AddAll(2, 4, 6)
	fmt.Println(y.String())

	x.UnionWith(&y)
	fmt.Printf("UnionWith: %s\n", x.String())

	x.IntersectWith(&y)
	fmt.Printf("IntersectWith: %s\n", x.String())
}
