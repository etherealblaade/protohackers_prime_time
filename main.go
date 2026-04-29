package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
)

type ExpectedJson struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

func handleConnection(conn net.Conn, ch chan ExpectedJson) {
	defer conn.Close()

	var ob ExpectedJson

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Bytes()
		if err := json.Unmarshal(message, &ob); err != nil {
			slog.Error("Invalid json", "error", err)
			conn.Write([]byte("Invalid JSON input\n"))
			continue
		}

		ch <- ob

	}

	close(ch)
}

func main() {
	ch := make(chan ExpectedJson)
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		slog.Error("Error listening socket", "error", err)
	}

	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		slog.Error("Error accepting connection", "error", err)
	}

	go handleConnection(conn, ch)

	for message := range ch {
		fmt.Println("New incoming message:", message)
	}

}
