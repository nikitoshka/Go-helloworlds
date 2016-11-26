package main

import (
	"fmt"

	"./counter"
)

func main() {
	text := []byte("some text that contains words\nanother line of the input\nand the last one\nbye =)")

	var b counter.ByteCounter
	b.Write(text)
	fmt.Println(b)

	var w counter.WordCounter
	w.Write(text)
	fmt.Println(w)

	var l counter.LineCounter
	l.Write(text)
	fmt.Println(l)
}
