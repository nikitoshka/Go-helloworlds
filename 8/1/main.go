package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

var host = flag.String("host", "localhost", "a host name to connect to")
var port = flag.Int("port", 8000, "a port to connect to")
var timeZone = flag.String("timeZone", "Europe/Moscow", "a timezine to get a tiem in")

func main() {
	flag.Parse()

	fullHostname := *host + ":" + strconv.Itoa(*port)
	fmt.Printf("Creating a server on %s\n", fullHostname)

	listener, err := net.Listen("tcp", fullHostname)
	if err != nil {
		log.Fatalf("Failed to create a server: %v\n", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Failed to establish a connection: %v", err)
			continue
		}
		log.Println("A connection occurred")

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	layout := "15:04:05\n"
	currentTime := time.Now()
	var timeString string

	loc, _ := time.LoadLocation(*timeZone)
	if loc != nil {
		currentTime = currentTime.In(loc)
	}
	timeString = currentTime.Format(layout)

	for _, err := io.WriteString(conn, timeString); err == nil; _, err = io.WriteString(conn, timeString) {
		time.Sleep(1 * time.Second)

		currentTime = time.Now()
		if loc != nil {
			currentTime = currentTime.In(loc)
		}
		timeString = currentTime.Format(layout)
	}
}
