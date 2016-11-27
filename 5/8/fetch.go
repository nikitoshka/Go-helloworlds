package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func fetch(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}

	f, err := os.Create(local)
	defer func() {
		log.Println("closing...")
		e := f.Close()
		if err == nil {
			err = e
		}
	}()

	_, err = io.Copy(f, resp.Body)

	return local, err
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("wrong usage")
	}

	if name, err := fetch(os.Args[1]); err != nil {
		log.Fatalln(err)
	} else {
		log.Println(name)
	}
}
