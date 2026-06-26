package main

import (
	"bufio"
	"fmt"
	"miniRedis/internal/resp"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("New client connected.")
	defer conn.Close()
	reader := bufio.NewReader(conn) // buffered reader because raw network sockets read in tiny chunks
	parser := resp.NewParser(reader)
	for {
		args, err := parser.Parse()
		if err != nil {
			fmt.Println("Client disconnected on error: ", err)
			break
		}
		fmt.Println("Received Command: %q\n", args)
	}
}
