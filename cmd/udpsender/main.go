package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Println("Error while resolving udp addr")
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error while dialing udp")
	}
	defer udpConn.Close()

	rdr := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		data, _, err := rdr.ReadLine()
		if err != nil {
			fmt.Println("Error while reading data using reader", err)
		}

		_, err = udpConn.Write(data)
		if err != nil {
			fmt.Println("Error while writing data using udpConn", err)
		}
	}
}
