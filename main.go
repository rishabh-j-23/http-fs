package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	linesChannel := make(chan string, 1)

	go func() {
		defer close(linesChannel)
		defer file.Close()
		dataString := ""
		for {
			buffer := make([]byte, 8)
			n, err := file.Read(buffer)
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
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error reading message file", err)
		os.Exit(1)
	}
	defer file.Close()

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
