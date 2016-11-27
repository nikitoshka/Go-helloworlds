package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"sort"
	"strings"
)

const defaultHost = "localhost:8000"

func main() {
	var hosts []string

	if len(os.Args) == 1 {
		hosts = append(hosts, defaultHost)
	} else {
		for _, host := range os.Args[1:] {
			hosts = append(hosts, host)
		}
	}

	conns := make(map[string]net.Conn)

	for _, host := range hosts {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			log.Printf("Failed to connect to %s: %v\n", host, err)
			continue
		}

		defer conn.Close()
		conns[host] = conn
	}

	for len(conns) != 0 {
		var row string
		s := getSorted(conns)

		for _, host := range s {
			if len(host) == 0 {
				continue
			}

			scanner := bufio.NewScanner(conns[host])
			if ok := scanner.Scan(); ok {
				t := strings.TrimSpace(scanner.Text())
				row += "\t" + "[" + host + "]" + t
			}

			if err := scanner.Err(); err != nil {
				delete(conns, host)
			}
		}

		log.Println(row)
	}
}

func getSorted(m map[string]net.Conn) []string {
	s := make([]string, len(m))

	for k := range m {
		s = append(s, k)
	}

	sort.Strings(s)
	return s
}
