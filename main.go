package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net"
)

type ExpectedRequest struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Resp struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

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
	defer conn.Close()
	buffer := bufio.NewReader(conn)
	malformedResponse, err := json.Marshal("malformed")
	if err != nil {
		slog.Error("marshaling malformed response", "error", err)
	}

	for {
		request, err := buffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			slog.Error("reading from buffer: %w", "error", err)
			return
		}

		res, err := validateJson(request)
		if err != nil {
			conn.Write(append(malformedResponse, '\n'))
			return
		}

		resp, _ := json.Marshal(Resp{Method: "isPrime", Prime: isPrime(*res.Number)})

		conn.Write(append(resp, '\n'))
	}
}

func validateJson(requestString string) (ExpectedRequest, error) {
	obj := ExpectedRequest{}
	if err := json.Unmarshal([]byte(requestString), &obj); err != nil {
		slog.Error("parsing request: %v", "error", err)
		return ExpectedRequest{}, fmt.Errorf("validating json: %w", err)
	}

	if obj.Method == "isPrime" {
		if obj.Number == nil {
			return ExpectedRequest{}, fmt.Errorf("missing number field")
		}
		return obj, nil
	}

	return ExpectedRequest{}, fmt.Errorf("isPrime field is not correct: %s", obj.Method)

}

func isPrime(num float64) bool {
	if num != math.Trunc(num) {
		return false
	}

	x := int64(num)

	if x < 2 {
		return false
	}

	if x == 2 || x == 3 {
		return true
	}

	if x%2 == 0 {
		return false
	}

	limit := int64(math.Sqrt(float64(x)))
	for d := int64(3); d <= limit; d += 2 {
		if x%d == 0 {
			return false
		}
	}

	return true
}
