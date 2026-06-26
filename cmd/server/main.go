package main

import (
	"bufio"
	"fmt"
	"miniRedis/internal/resp"
	"miniRedis/internal/storage"
	"net"
	"strings"
)

func main() {
	// init central storage engine once
	db := storage.NewEngine()

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Mini-Redis listening on port 6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			continue
		}
		// pass the shared engine to every client
		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *storage.Engine) {
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
		// ignore empty commands
		if len(args) == 0 {
			continue
		}

		// convert command to uppercase to make it case insensitive
		command := strings.ToUpper(args[0])

		switch command {
		case "PING":
			conn.Write(([]byte("+PONG\r\n")))

		case "SET":
			if len(args) != 3 { // [SET, "hello", "vedansh"]
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
				continue
			}
			db.Set(args[1], args[2])
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(args) != 2 {
				conn.Write([]byte("-ERR wrong number of argumnets for 'get' command\r\n"))
				continue
			}
			val, ok := db.Get(args[1])
			if !ok {
				// "$-1\r\n" is the exact RESP protocol for a Null String (key not found)
				conn.Write([]byte("$-1\r\n"))
				continue
			}
			// respond with a bult string containing the value
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))

		case "DEL":
			if len(args) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'del' command\r\n"))
				continue
			}
			db.Delete(args[1])
			conn.Write([]byte(":1\r\n")) // RESP for Integer 1(success)

		default:
			conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command)))

		}
		fmt.Printf("Received Command: %q\n", args)
	}
}
