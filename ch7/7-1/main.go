package main

import (
	"bufio"
	"bytes"
	"fmt"
)

// ByteCounter 型
type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // intをByteCounterへ変換
	return len(p), nil
}

// WordCounter 型
type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	reader := bytes.NewReader(p)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		*c++
	}
	return len(p), nil
}

func main() {
	var c ByteCounter
	c.Write([]byte("hello"))
	fmt.Println(c)

	c = 0
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c)

	var w WordCounter
	w.Write([]byte("hello kentana"))
	fmt.Println(w)
}
