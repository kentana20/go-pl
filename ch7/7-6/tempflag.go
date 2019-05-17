package main

import (
	"flag"
	"fmt"

	"github.com/kentana20/go-pl/ch7/7-6/tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}
