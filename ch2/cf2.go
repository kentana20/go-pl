package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/kentana20/go-pl/ch2/tempconv"
)

func main() {
	for _, arg := range os.Args[1:] {
		value := arg
		if value == "" {
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()
			value = stdin.Text()
		}

		t, err := strconv.ParseFloat(value, 64)

		if err != nil {
			fmt.Fprintf(os.Stderr, "cf: %v\n", err)
			os.Exit(1)
		}

		f := tempconv.Fahrenheit(t)
		c := tempconv.Celsius(t)

		fmt.Printf("%s = %s, %s = %s\n", f, tempconv.FToC(f), c, tempconv.CToF(c))
	}
}
