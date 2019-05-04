package main

import (
	"fmt"

	"github.com/kentana20/go-pl/ch6/intset"
)

func main() {
	var x intset.IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	x.Add(1020)
	fmt.Println(x.String())
	fmt.Println(x.Len())
	x.AddAll(5, 7, 10)
	fmt.Println(x.String())
	fmt.Println(x.Len())
}
