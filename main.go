package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalf("error creating listening socket: %v", err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("error accepting connection", "error", err)
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	buffer := bufio.NewReader(conn)

	for {
		response, err := buffer.ReadString('\n')
		if err != nil {
			slog.Error("reading from buffer: %w", "error", err)
		}

		fmt.Fprintf(conn, "Echo: %s", response)
	}
}
