package main

import "fmt"

var x uint8 = 1<<1 | 1<<5
var y uint8 = 1<<1 | 1<<2

var apples int32 = 1
var oranges int16 = 2
var compote = int(apples) + int(oranges)

func main() {
	fmt.Printf("%08b\n", x)
	fmt.Printf("%08b\n", y)

	for i := uint(0); i < 8; i++ {
		if x&(1<<i) != 0 {
			fmt.Printf("%08b\n", x<<1)
			fmt.Println(i)
		}
	}
}
