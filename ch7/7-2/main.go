package main

import (
	"fmt"
	"io"
	"os"
)

type countingWriter struct {
	w     io.Writer
	count int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.count += int64(n)
	return n, err
}

// CountingWriter - io.Writerとバイト数へのポインタを返す関数
func CountingWriter(w io.Writer) (io.Writer, *int64) {
	cw := countingWriter{w, 0}
	return &cw, &cw.count
}

func main() {
	writer, count := CountingWriter(os.Stdout)
	fmt.Fprint(writer, "123456789\n")
	fmt.Println(*count)
}
