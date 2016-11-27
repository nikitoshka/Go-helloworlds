package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	fileNameCounter = make(map[string]uint)
	urlToFile       = make(map[string]string)
	fileNameMut     sync.Mutex
)

var cancelch = make(chan struct{})
var donech = make(chan struct{})

var ioMut sync.Mutex

var defaultFileName = "index.html"

var (
	urlCounter    int
	urlCounterMut sync.RWMutex
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("##### wrong usage")
	}

	for i := 1; i < len(os.Args); i++ {
		urlCounterMut.Lock()
		urlCounter++
		urlCounterMut.Unlock()

		go fetch(os.Args[i])
	}

	<-donech
}

func fetch(url string) {
	log.Printf("----- [%s] getting\n", url)

	defer func() {
		select {
		case <-donech:
			log.Printf("----- [%s:defer] done\n", url)
			return

		default:
			log.Printf("----- [%s:defer] checking for exhaustion of the requests\n", url)
			urlCounterMut.Lock()
			defer urlCounterMut.Unlock()
			urlCounter--

			if urlCounter <= 0 {
				log.Printf("----- [%s:defer] closing the chan\n", url)
				close(donech)
			}
		}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("##### failed to make a request for %s: %v\n", url, err)
		return
	}

	req.Cancel = cancelch

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("##### failed to get a response for %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("----- [%s] locking the ioMut\n", url)
	ioMut.Lock()
	defer ioMut.Unlock()

	if isDone() {
		log.Printf("----- [%s] already done\n", url)
		return
	}

	log.Printf("----- [%s] proceeding to saving\n", url)

	fileName := getFileName(resp.Request.URL.Path)
	if len(fileName) == 0 {
		log.Printf("##### failed to get a filename for %s\n", resp.Request.URL.String())
		return
	}

	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("##### failed to create a file %s for %s\n", fileName, resp.Request.URL.String())
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Printf("##### failed to copy a response of %s into %s\n", resp.Request.URL.String(), fileName)
		return
	}

	log.Printf("----- [%s] done\n", url)
	close(cancelch)
	close(donech)
}

func getFileName(url string) string {
	fileNameMut.Lock()
	defer fileNameMut.Unlock()

	var fileName string

	if strings.HasSuffix(url, "/") {
		fileName = defaultFileName
	} else {
		fileName = url[strings.LastIndex(url, "/"):]
	}

	if strings.HasPrefix(fileName, "/") {
		fileName = fileName[strings.Index(fileName, "/")+1:]
	}

	fileNameCounter[fileName]++

	if count := fileNameCounter[fileName]; count > 1 {
		dotPosition := strings.LastIndex(fileName, ".")

		if dotPosition < 0 {
			fileName = fileName + "_" + strconv.Itoa(int(count))
		} else {
			fileName = fileName[:dotPosition] + "_" + strconv.Itoa(int(count)) + fileName[dotPosition:]
		}
	}

	return fileName
}

func isDone() bool {
	select {
	case <-donech:
		return true

	default:
		return false
	}
}
