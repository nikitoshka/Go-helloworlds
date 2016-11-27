package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("ERROR: failed to open a connection: %v\n", err)
	}
	defer conn.Close()

	go readMessages(conn)

	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Println("ERROR: failed to read a piece of input")
	}
}

func readMessages(conn net.Conn) {
	if _, err := io.Copy(os.Stdout, conn); err != nil {
		log.Printf("ERROR: failed to read a message")
		return
	}
}
