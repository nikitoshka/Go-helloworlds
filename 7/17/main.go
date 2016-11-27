package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("wrong number of arguments")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	dec := xml.NewDecoder(f)
	var stack []string

	for {
		token, err := dec.Token()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalln(err)
			}
		}

		switch token := token.(type) {
		case xml.StartElement:
			stack = append(stack, token.Name.Local)

			for _, attr := range token.Attr {
				stack = append(stack, attr.Name.Local)

				if containsAll(stack, os.Args[2:]) {
					fmt.Printf("%s:\t%s\n", strings.Join(stack, " "), attr.Value)
				}

				stack = stack[:len(stack)-1]
			}

		case xml.EndElement:
			stack = stack[:len(stack)-1]

		case xml.CharData:
			if containsAll(stack, os.Args[2:]) {
				fmt.Printf("%s:\t%s\n", strings.Join(stack, " "), token)
			}
		}
	}
}

func containsAll(x, y []string) bool {
	if len(y) > len(x) {
		return false
	}

	coincedences := 0

	for i := 0; i < len(x); i++ {
		if x[i] == y[coincedences] {
			coincedences++
		} else {
			coincedences = 0
		}
	}

	return coincedences == len(y)
}
