package main

import (
	"bytes"
	"fmt"

	"./countingWriter"
)

func main() {
	var buf bytes.Buffer
	text := []byte("some text that contains words\nanother line of the input\nand the last one\nbye =)")
	cw, count := countingWriter.CountingWriter(&buf)

	cw.Write(text)

	fmt.Println(buf.String(), "\n", *count)
}
