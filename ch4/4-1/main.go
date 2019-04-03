package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	str1 := "abc"
	str2 := "xyz"
	hash1 := sha256.Sum256([]byte(str1))
	hash2 := sha256.Sum256([]byte(str2))

	fmt.Printf("hash1: %x \n", hash1)
	fmt.Printf("hash2: %x \n", hash2)
	dif := 0
	for i := 0; i < 32; i++ {
		// hash1とhash2のXORを数える
		dif += popcount(hash1[i] ^ hash2[i])
	}
	fmt.Printf("dif: %d\n", dif)
}

func popcount(x byte) int {
	count := 0
	for x != 0 {
		count++
		x &= x - 1
	}
	return count
}
