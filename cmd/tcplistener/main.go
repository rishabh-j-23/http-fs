package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net"
)

const PORT = ":42069"

func getLinesChannel(conn net.Conn) <-chan string {
	linesChannel := make(chan string, 1)

	go func() {
		defer close(linesChannel)
		defer conn.Close()
		dataString := ""
		for {
			buffer := make([]byte, 8)
			n, err := conn.Read(buffer)
			if err != nil {
				break
			}

			buffer = buffer[:n]

			if i := bytes.IndexByte(buffer, '\n'); i != -1 {
				dataString += string(buffer[:i])
				linesChannel <- dataString
				buffer = buffer[i+1:]
				dataString = ""
			}

			dataString += string(buffer)
		}

		if len(dataString) > 0 {
			linesChannel <- dataString
		}
	}()

	return linesChannel
}

func main() {

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		slog.Error("Error creating server on", PORT, err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Error while accepting conn", "err", err)
		} else {
			fmt.Println("Connection has being accepted")
		}
		defer conn.Close()

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
	}
	defer listener.Close()
}
