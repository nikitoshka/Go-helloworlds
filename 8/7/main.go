package main

import (
	"flag"
	"log"

	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"the_go_progr_lang/8/7/links"
)

var urlToMirror = flag.String("url", "", "an urlToMirror to mirror")

func main() {
	flag.Parse()

	if len(*urlToMirror) == 0 {
		log.Fatalln("wrong usage")
	}

	links, err := links.Exract(*urlToMirror)
	if err != nil {
		log.Fatalf("extract error: %s\n", err)
	}

	origurlToMirror, err := url.Parse(*urlToMirror)
	if err != nil {
		log.Fatalf("failed to parse urlToMirror: %v\n", err)
	}

	linkschan := make(chan string, len(links)+1)
	resultchan := make(chan struct{})

	go getLink(origurlToMirror.Host, linkschan, resultchan)

	for _, link := range links {
		linkschan <- link
	}

	close(linkschan)

	for range links {
		<-resultchan
	}
}

func getLink(host string, linkschan <-chan string, resultchan chan<- struct{}) {
	for link := range linkschan {
		resp, err := http.Get(link)
		if err != nil {
			log.Printf("error while getting \"%s\": %v\n", link, err)
			resultchan <- struct{}{}
			continue
		}

		go save(host, resp, resultchan)
	}
}

func save(host string, resp *http.Response, resultchan chan<- struct{}) {
	defer func() {
		resp.Body.Close()
		resultchan <- struct{}{}
	}()

	u := resp.Request.URL

	if u.Host != host {
		log.Printf("foreign host: %s\n", u.Host)
		return
	}

	path := u.Path[1:]
	log.Printf("path: %s\n", path)

	if len(path) != 0 {
		if err := os.MkdirAll(path, 0777); err != nil {
			log.Printf("failed to create a foder %s: %v\n", path, err)
			return
		}
	}

	filename := u.Path[strings.LastIndex(u.Path, "/")+1:]
	if len(filename) == 0 {
		filename = "index.html"
	}

	log.Printf("filename: %s\n", filename)

	fullPath := path + filename
	f, err := os.Create(fullPath)
	if err != nil {
		log.Printf("failed to open a file %s: %v\n", fullPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Printf("failed to write to the file %s: %v\n", fullPath, err)
		return
	}
}
