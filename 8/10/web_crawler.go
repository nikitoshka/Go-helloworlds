package main

import (
	"log"
	"os"

	"flag"
	"fmt"
	"the_go_progr_lang/8/10/links"
	"time"
)

type workLinks struct {
	links []string
	depth int
}

var workch = make(chan workLinks)
var restrch = make(chan struct{}, 20)
var closech = make(chan struct{})
var donech = make(chan struct{})
var workCounter int
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

	workCounter++
	go func() {
		workch <- workLinks{args, 0}
	}()

	go terminateByInput()
	go waitForCompletion()
	go goThrough(crawler)

	<-donech

	fmt.Printf("time elapsed: %s\n", time.Since(now))
}

func crawler(url string) []string {
	restrch <- struct{}{}

	log.Println(url)

	links, err := links.Exract(url, closech)
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
			workCounter--
			continue
		}

		if isCanceled() {
			return
		}

		for _, link := range worklist.links {
			if seen[link] {
				continue
			}

			if isCanceled() {
				return
			}

			seen[link] = true

			workCounter++
			go func(link string, depth int) {
				if isCanceled() {
					return
				}

				workch <- workLinks{f(link), depth + 1}
			}(link, worklist.depth)
		}

		workCounter--
	}
}

func terminateByInput() {
	os.Stdin.Read(make([]byte, 1))

	select {
	case <-donech:
		return

	default:
		log.Println("########## CANCELED ##########")
		close(closech)
		close(workch)
		close(donech)
	}
}

func waitForCompletion() {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("!!!!!!!!!!\tPANIC OCCURRED: %v\t!!!!!!!!!!\n", p)
		}
	}()

	for {
		select {
		case <-donech:
			return

		default:
			if workCounter > 0 {
				continue
			}

			log.Println("########## DONE ##########")

			close(closech)
			close(workch)
			close(donech)
			return
		}
	}
}

func isCanceled() bool {
	select {
	case <-donech:
		return true
	default:
		return false
	}
}
