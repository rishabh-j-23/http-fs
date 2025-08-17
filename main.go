package main

import (
	"fmt"
	"os"
)

func main() {
	buffer := make([]byte, 8)
	file, err := os.Open("messages.txt")
	defer file.Close()

	if err != nil {
		fmt.Println("Error reading message file", err)
		os.Exit(1)
	}

	for {
		bytes, err := file.Read(buffer)
		if err != nil {
			os.Exit(0)
		}
		if bytes > 0 {
			fmt.Println("read:", string(buffer[:bytes]))
		}
	}
}
