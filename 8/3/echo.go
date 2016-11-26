package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const protocol = "tcp"

var host = flag.String("host", "localhost", "a host to connect to")
var port = flag.Int("port", 8000, "a port to connect to")

func main() {
	flag.Parse()

	hostname := *host + ":" + strconv.Itoa(*port)
	log.Printf("Connecting to %s...\n", hostname)

	listener, err := net.Listen(protocol, hostname)
	if err != nil {
		log.Fatalf("Failed to create a server on %s: %v\n", hostname, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to establish a connection: %v\n", err)
			continue
		}

		go handConn(conn)
	}
}

func handConn(conn net.Conn) {
	defer conn.Close()

	var wait sync.WaitGroup
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Printf("An error occurred during reading from the connection: %v\n", err)
			return
		}

		wait.Add(1)
		go echo(scanner.Text(), 1*time.Second, conn, &wait)
	}

	go func() {
		wait.Wait()

		if w, ok := conn.(*net.TCPConn); ok {
			w.CloseWrite()
		}
	}()
}

func echo(str string, duration time.Duration, conn net.Conn, wait *sync.WaitGroup) {
	defer func() {
		if wait != nil {
			wait.Done()
		}
	}()

	if _, err := io.WriteString(conn, strings.ToUpper(str)+"\n"); err != nil {
		return
	}
	time.Sleep(duration)

	if _, err := io.WriteString(conn, str+"\n"); err != nil {
		return
	}
	time.Sleep(duration)

	if _, err := io.WriteString(conn, strings.ToLower(str)+"\n"); err != nil {
		return
	}
}
