package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
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

	newInput := make(chan string)
	inputBroke := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(conn)

		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Printf("An error occurred during reading from the connection: %v\n", err)
				inputBroke <- struct{}{}
			}

			newInput <- scanner.Text()
		}
	}()

	for {
		select {
		case str := <-newInput:
			go echo(str, 1*time.Second, conn)

		case <-inputBroke:
			log.Println("Closing a connectioin because of input trouble")
			return

		case <-time.After(10 * time.Second):
			io.WriteString(conn, "the connection is closed because of 10s inactivity\n")
			return
		}
	}
}

func echo(str string, duration time.Duration, conn net.Conn) {
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
