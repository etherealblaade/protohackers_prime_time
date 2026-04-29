package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
)

type ExpectedJson struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

func handleConnection(conn net.Conn, ch chan string) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ch <- scanner.Text()
	}

	close(ch)
}

func main() {
	ch := make(chan string)
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
