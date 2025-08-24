package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/rishabh/http-fs/internal/request"
)

const PORT = ":42069"

func main() {

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		slog.Error("Error creating server on", PORT, err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Error while accepting conn", "err", err)
		} else {
			fmt.Println("Connection has being accepted")
		}
		defer conn.Close()

		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error reading the request")
			break
		}
		fmt.Printf("Request Line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)

	}
}
