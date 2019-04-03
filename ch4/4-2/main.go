package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
)

func main() {
	i := flag.Int("i", 256, "sha")
	s := flag.String("s", "hoge", "str")
	flag.Parse()

	switch *i {
	case 384:
		fmt.Printf("sha384: %x\n", sha512.Sum384([]byte(*s)))
	case 512:
		fmt.Printf("sha512: %x\n", sha512.Sum512([]byte(*s)))
	default:
		fmt.Printf("sha256: %x\n", sha256.Sum256([]byte(*s)))
	}
}
