package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

var indent = 0

type Node interface{}
type CharData string
type Element struct {
	Type     string
	Attr     map[string]string
	Children []Node
}

func buildElement(dec *xml.Decoder, element *Element) {
	if dec == nil || element == nil {
		panic("nullptr")
	}

	token, err := dec.Token()
	if err != nil {
		panic(fmt.Sprintf("error while getting a token: %s", err))
	}

	for {
		switch token := token.(type) {
		case xml.StartElement:
			elem := &Element{Attr: make(map[string]string)}
			elem.Type = token.Name.Local
			for _, v := range token.Attr {
				elem.Attr[v.Name.Local] = v.Value
			}

			if len(element.Type) == 0 {
				*element = *elem
			} else {
				buildElement(dec, elem)
				element.Children = append(element.Children, *elem)
			}

		case xml.CharData:
			cd := CharData(token)
			element.Children = append(element.Children, cd)

		case xml.EndElement:
			return
		}

		token, err = dec.Token()
		if err != nil {
			panic(fmt.Sprintf("error while getting a token: %s", err))
		}
	}
}

func printElement(element *Element) {
	if element == nil {
		panic("nullptr")
	}

	indent++
	fmt.Printf("%*s<%s", 2*indent, " ", element.Type)

	for k, v := range element.Attr {
		fmt.Printf(" %s=%s", k, v)
	}

	fmt.Printf(">\n")

	for _, child := range element.Children {
		switch child := child.(type) {
		case CharData:
			indent++
			fmt.Printf("%*s%s\n", 2*indent, " ", child)
			indent--

		case Element:
			fmt.Printf("\n")
			printElement(&child)
		}
	}

	fmt.Printf("%*s</%s>\n", 2*indent, " ", element.Type)
	indent--
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("wrong number of arguments")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	dec := xml.NewDecoder(f)
	rootElement := &Element{Attr: make(map[string]string)}

	buildElement(dec, rootElement)
	printElement(rootElement)
}
