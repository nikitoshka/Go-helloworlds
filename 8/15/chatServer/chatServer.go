package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
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

const idleDuration = 300
const clientMessagesCapacity = 5

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
				select {
				case cli.ch <- msg:
				default:
				}
			}

		case cli := <-leaving:
			close(cli.ch)
			delete(clients, cli)
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	who, err := makeIntroduction(conn)
	if err != nil {
		log.Printf("ERROR: failed to make an introductioin with %s: %v\n",
			conn.RemoteAddr().String(), err)
		return
	}

	cli := client{make(chan string, clientMessagesCapacity), who}
	activity := make(chan struct{})

	go disconnectIdle(conn, activity)
	go clientWriter(cli, conn)

	entering <- cli

	messages <- cli.who + " has connected"

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if err := input.Err(); err != nil {
			log.Printf("ERROR: error while reading from %s: %v\n", cli.who, err)
			break
		}

		activity <- struct{}{}
		messages <- "[" + cli.who + "]" + ": " + input.Text()
	}

	leaving <- cli
	messages <- cli.who + " has left"
}

func clientWriter(cli client, conn net.Conn) {
	for msg := range cli.ch {
		if _, err := fmt.Fprintln(conn, msg); err != nil {
			log.Printf("ERROR: error while writing to %s: %v\n", cli.who, err)
			continue
		}
	}
}

func makeIntroduction(conn net.Conn) (who string, err error) {
	if _, err = fmt.Fprintf(conn, "Enter yout name: "); err != nil {
		return
	}

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return
		}

		who = scanner.Text()
	} else {
		err = fmt.Errorf("ERROR: failed to read a name")
	}

	return
}

func disconnectIdle(conn net.Conn, activity <-chan struct{}) {
	ticker := time.NewTicker(idleDuration * time.Second)

	for {
		select {
		case <-ticker.C:
			log.Printf("WARNING: The time is over for %s\n", conn.RemoteAddr().String())
			conn.Close()
			ticker.Stop()
			return

		case <-activity:
			ticker.Stop()
			ticker = time.NewTicker(idleDuration * time.Second)
		}
	}
}
