package main

import (
	"io"

	"fmt"

	"golang.org/x/net/html"
)

type StringReader string

func (sr *StringReader) Read(p []byte) (n int, err error) {
	n = len(*sr)
	copy(p, []byte(*sr))
	err = io.EOF

	return
}

func NewReader(s string) io.Reader {
	sr := StringReader(s)
	return &sr
}

func main() {
	sr := NewReader("<html><body><h1>EXAMPLE</h1></body></html>")
	doc, err := html.Parse(sr)

	if err != nil {
		panic(fmt.Sprintf("error during parsing: %v", err))
	}

	fmt.Println(doc.FirstChild.Data)
}
