package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var verbose = flag.Bool("v", false, "verbose mode")

var wg sync.WaitGroup
var sema = make(chan struct{}, 20)

type rootSize struct {
	root string
	size int64 // bytes
}

func main() {
	flag.Parse()

	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	filesizes := make(chan rootSize)
	var ticker <-chan time.Time

	if *verbose {
		ticker = time.Tick(1 * time.Second)
	}

	for _, dirname := range roots {
		wg.Add(1)
		go walkDir(dirname, dirname, filesizes)
	}

	go func() {
		wg.Wait()
		close(filesizes)
	}()

	rootFiles := make(map[string]int64)
	rootSizes := make(map[string]int64)

loop:
	for {
		select {
		case <-ticker:
			printSize(rootFiles, rootSizes)

		case size, ok := <-filesizes:
			if !ok {
				break loop
			}

			rootFiles[size.root]++
			rootSizes[size.root] += size.size
		}
	}

	printSize(rootFiles, rootSizes)
}

func walkDir(rootname string, dirname string, filesize chan<- rootSize) {
	defer wg.Done()

	for _, info := range getDirInfo(dirname) {
		if info.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dirname, info.Name())
			walkDir(rootname, subdir, filesize)
		} else {
			filesize <- rootSize{rootname, info.Size()}
		}
	}
}

func getDirInfo(dirName string) []os.FileInfo {
	sema <- struct{}{}
	defer func() {
		<-sema
	}()

	info, err := ioutil.ReadDir(dirName)

	if err != nil {
		return nil
	}

	return info
}

func printSize(files, sizes map[string]int64) {
	for name, count := range files {
		size, ok := sizes[name]
		if !ok {
			continue
		}

		fmt.Printf("%s (%d files): %.3f MB\n", name, count, float64(size)/1024/1024)
	}
}
