package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	ch  chan string
	who string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("ERROR: failed to start listening: %v\n", err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("WARNING: failed to establish a connection: %v\n", err)
			continue
		}

		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)

	for {
		select {
		case cli := <-entering:
			var online string
			for c := range clients {
				online += " " + c.who
			}
			clients[cli] = true

			if len(online) > 0 {
				online = "Online users: " + online
				go func() {
					cli.ch <- online
				}()
			}

		case msg := <-messages:
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-leaving:
			close(cli.ch)
			delete(clients, cli)
		}
	}
}

func handleConn(conn net.Conn) {
	cli := client{make(chan string), conn.RemoteAddr().String()}

	go clientWriter(cli, conn)

	entering <- cli

	messages <- cli.who + " has connected"

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if err := input.Err(); err != nil {
			log.Printf("ERROR: error while reading from %s: %v\n", cli.who, err)
			break
		}

		messages <- "[" + cli.who + "]" + ": " + input.Text()
	}

	leaving <- cli
	messages <- cli.who + " has left"
	conn.Close()
}

func clientWriter(cli client, conn net.Conn) {
	for msg := range cli.ch {
		if _, err := fmt.Fprintln(conn, msg); err != nil {
			log.Printf("ERROR: error while writing to %s: %v\n", cli.who, err)
			continue
		}
	}
}
