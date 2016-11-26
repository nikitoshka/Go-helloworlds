package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type cmdType int

const (
	nan cmdType = iota
	ls
	close
	cd
)

var availableCommands = map[cmdType]string{
	ls:    "ls",
	close: "close",
	cd:    "cd",
}

var host = flag.String("host", "localhost", "a host to establish the server on")
var port = flag.Int("port", 8000, "a port to establish the server on")

func main() {
	flag.Parse()
	hostname := *host + ":" + strconv.Itoa(*port)

	listener, err := net.Listen("tcp", hostname)
	if err != nil {
		log.Fatalf("Failed to establish the server: %v\n", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to establish a connection: %v\n", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	io.WriteString(conn, "The connection's established\n")

	curDir := getCommandPrompt()
	io.WriteString(conn, curDir+" ")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Printf("An error during reading occured: %v\n", err)
			return
		}

		cmd := scanner.Text()

		curDirChanged, connClosed := execCmd(cmd, conn)
		if connClosed {
			log.Println("The connection's closed")
			return
		}
		if curDirChanged {
			curDir = getCommandPrompt()
		}

		io.WriteString(conn, curDir+" ")
	}
}

func getCommandPrompt() string {
	if path, err := os.Getwd(); err != nil {
		return "[~]"
	} else {
		return "[" + path + "]"
	}
}

func execCmd(cmd string, conn net.Conn) (curDirChanged, connClosed bool) {
	var err error
	comType := getCmdType(cmd)

	switch comType {
	case nan:
		err = fmt.Errorf("command not found\n")
	case ls:
		err = execLS(cmd, conn)
	case cd:
		err = execCD(cmd, conn)
		if err == nil {
			curDirChanged = true
		}
	case close:
		conn.Close()
		connClosed = true

	default:
		err = fmt.Errorf("command not found\n")
	}

	if err != nil {
		io.WriteString(conn, err.Error())
	}

	return
}

func execLS(cmd string, conn net.Conn) error {
	comType := getCmdType(cmd)
	if comType != ls {
		return fmt.Errorf("command not found: %s", cmd)
	}

	args := strings.Fields(cmd)
	args = args[1:]

	if len(args) == 0 {
		currentDir := "."
		contents, err := getDirContents(currentDir)

		if err != nil {
			return err
		}

		sendLSResult(currentDir, contents, conn)
		return nil
	}

	for _, dir := range args {
		contents, err := getDirContents(dir)
		if err != nil {
			return err
		}

		sendLSResult(dir, contents, conn)
	}

	return nil
}

func execCD(cmd string, conn net.Conn) error {
	comType := getCmdType(cmd)
	if comType != cd {
		return fmt.Errorf("command not found: %s", cmd)
	}

	args := strings.Fields(cmd)
	if len(args) != 2 {
		return fmt.Errorf("%s: wrong syntax\n", availableCommands[comType])
	}

	return os.Chdir(args[1])
}

func getDirContents(dir string) (map[string]int64, error) {
	contents, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	descr := make(map[string]int64)
	for _, entity := range contents {
		descr[entity.Name()] = entity.Size()
	}

	return descr, nil
}

func sendLSResult(dir string, contents map[string]int64, conn net.Conn) {
	io.WriteString(conn, dir+":\n")
	for name, size := range contents {
		io.WriteString(conn, name+":\t"+strconv.FormatInt(size, 10)+"\n")
	}
}

func getCmdType(cmd string) cmdType {
	for k, v := range availableCommands {
		if strings.HasPrefix(cmd, v) {
			return k
		}
	}

	return nan
}
