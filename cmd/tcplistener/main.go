package main

import (
	"fmt"
	"net"
	"os"
	"tcptohttp/internal/request"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("cannot setup listener", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("cannot accept connection", err)
			return
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Cannot read from connection", err)
			return
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

	}
}

func getFromFile() *os.File {

	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("File cannot be opened", err)
		return file
	}
	return file
}

/*
func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(ch)

		line := ""
		for {

			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				line += string(data[:i])
				data = data[i+1:]
				ch <- line
				line = ""
			}

			line += string(data)

		}

		if len(line) != 0 {
			ch <- line
		}

	}()

	return ch

}
*/
