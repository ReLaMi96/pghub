package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("failed to resolve udp address", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Failed to make connection", err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read line", err)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println("failed to write to conn", err)
		}
	}

}
