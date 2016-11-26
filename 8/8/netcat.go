package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

const protocol = "tcp"

var host = flag.String("host", "localhost", "a host to connect to")
var port = flag.Int("port", 8000, "a port to connect to")

func main() {
	flag.Parse()

	hostname := *host + ":" + strconv.Itoa(*port)
	log.Printf("Connecting to %s...\n", hostname)

	conn, err := net.Dial(protocol, hostname)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v\n", hostname, err)
	}

	defer conn.Close()

	ch := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn)
		io.WriteString(os.Stdout, "done")
		ch <- struct{}{}
	}()

	mustCopy(conn, os.Stdin)

	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.CloseWrite()
	}

	<-ch
}

func mustCopy(w io.Writer, r io.Reader) {
	if _, err := io.Copy(w, r); err != nil {
		log.Fatalf("Failed to get a piece of data on the connection: %v\n", err)
	}
}
