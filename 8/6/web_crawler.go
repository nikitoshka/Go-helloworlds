package main

import (
	"log"
	"os"

	"flag"
	"fmt"
	"sync"
	"the_go_progr_lang/8/6/links"
	"time"
)

type workLinks struct {
	links []string
	depth int
}

var workch = make(chan workLinks)
var restrch = make(chan struct{}, 20)
var wg sync.WaitGroup
var depth = flag.Int("depth", 0, "the depth of crawling")

func main() {
	flag.Parse()

	now := time.Now()

	if *depth == 0 {
		fmt.Println("infinite crawling")
	} else {
		fmt.Printf("crawling with the depth = %d\n", *depth)
	}

	var args []string
	if os.Args[1] == "-depth" {
		args = os.Args[3:]
	} else {
		args = os.Args[1:]
	}

	wg.Add(1)
	go func() {
		workch <- workLinks{args, 0}
	}()

	go goThrough(crawler)

	wg.Wait()
	close(workch)

	fmt.Printf("time elapsed: %s\n", time.Since(now))
}

func crawler(url string) []string {
	restrch <- struct{}{}

	log.Println(url)

	links, err := links.Exract(url)
	if err != nil {
		log.Println(err)
	}

	<-restrch
	return links
}

func goThrough(f func(item string) []string) {
	seen := make(map[string]bool)

	for worklist := range workch {
		if *depth != 0 && worklist.depth > *depth {
			wg.Done()
			continue
		}

		for _, link := range worklist.links {
			if seen[link] {
				continue
			}

			seen[link] = true

			wg.Add(1)
			go func(link string, depth int) {
				workch <- workLinks{f(link), depth + 1}
			}(link, worklist.depth)
		}

		wg.Done()
	}

	// for len(worklist) > 0 {
	// 	items := worklist
	// 	worklist = nil

	// 	for _, item := range items {
	// 		if !seen[item] {
	// 			seen[item] = true
	// 			worklist = append(worklist, f(item)...)
	// 		}
	// 	}
	// }
}
