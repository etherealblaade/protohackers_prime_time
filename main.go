package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
)

type ExpectedJson struct {
	Method *string `json:"method"`
	Number *int    `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
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

func handleMessage(conn net.Conn, message ExpectedJson) {
	resp := Response{Method: "isPrime", Prime: true}
	jresp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("error marshaling response", "error", err)
	}

	if validateJson(message) {
		conn.Write(append(jresp, '\n'))
	} else {
		conn.Write([]byte("malformed"))
	}
}

func validateJson(json ExpectedJson) bool {
	if json.Method == nil || json.Number == nil {
		return false
	}

	if *json.Method == "isPrime" {
		return true
	}

	return false
}

func main() {
	ch := make(chan ExpectedJson)
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		slog.Error("Error listening socket", "error", err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("Error accepting connection", "error", err)
		}

		go handleConnection(conn, ch)

		for message := range ch {
			fmt.Println("New incoming message:", message)
			handleMessage(conn, message)
		}
	}
}
