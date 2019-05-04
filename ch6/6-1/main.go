package main

import (
	"fmt"

	"github.com/kentana20/go-pl/ch6/intset"
)

func main() {
	var x, y, z intset.IntSet
	// Add
	x.Add(1)
	x.Add(144)
	x.Add(9)
	x.Add(1020)
	fmt.Println(x.String())
	fmt.Println(x.Len())

	y.Add(9)
	y.Add(42)
	y.Add(1000)
	fmt.Println(y.String())
	fmt.Println(y.Len())
	y.Remove(100)
	fmt.Println(y.String())
	fmt.Println(y.Len())

	z = *y.Copy()
	y.Clear()
	fmt.Println(y.String())
	fmt.Println(y.Len())
	fmt.Println(z.String())
	fmt.Println(z.Len())
}
