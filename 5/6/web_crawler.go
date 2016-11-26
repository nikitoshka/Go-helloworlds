package main

import (
	"log"
	"os"
	"the_go_progr_lang/5_6/links"
)

func main() {
	goThrough(crawler, os.Args[1:])
}

func crawler(url string) []string {
	log.Println(url)

	links, err := links.Exract(url)
	if err != nil {
		log.Println(err)
	}

	return links
}

func goThrough(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)

	for len(worklist) > 0 {
		items := worklist
		worklist = nil

		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}
