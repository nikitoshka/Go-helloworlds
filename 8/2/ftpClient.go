package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

var host = flag.String("host", "localhost", "a host name to connect to")
var port = flag.Int("port", 8000, "a port to connect to")

func main() {
	flag.Parse()

	fullHostname := *host + ":" + strconv.Itoa(*port)
	log.Printf("Connecting to %s...\n", fullHostname)

	conn, err := net.Dial("tcp", fullHostname)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v\n", fullHostname, err)
	}

	defer conn.Close()

	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	conn.Write([]byte(scanner.Text()))
	// }

	// go func() {
	// 	scanner := bufio.NewScanner(conn)
	// 	for scanner.Scan() {
	// 		io.WriteString(os.Stdout, scanner.Text())
	// 	}
	// }()

	go mustCopy(os.Stdout, conn)
	mustCopy(conn, os.Stdin)
}

func mustCopy(w io.Writer, r io.Reader) {
	if _, err := io.Copy(w, r); err != nil {
		log.Fatalf("Failed to get a piece of data on the connection: %v\n", err)
	}
}
