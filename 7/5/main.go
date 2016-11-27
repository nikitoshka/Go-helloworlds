package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type limitReader struct {
	r io.Reader
	n int64
}

func (r *limitReader) Read(p []byte) (n int, err error) {
	b := make([]byte, r.n)
	n, err = r.r.Read(b)

	if err != nil && err != io.EOF {
		return
	}

	copy(p, b)

	return n, io.EOF
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &limitReader{r, n}
}

func main() {
	r := strings.NewReader("test string; kinda cool")
	lr := LimitReader(r, 11)

	b, err := ioutil.ReadAll(lr)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))
}
