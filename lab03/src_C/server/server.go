package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func runServer(port string) {
	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	request_storage := make([]byte, 1024)
	n, err := conn.Read(request_storage)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading from socet, ", err)
		return
	}
	request_storage = request_storage[:n]
	request_part := strings.Split(string(request_storage), " ")
	if len(request_part) < 2 {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}
	fileName := strings.Trim(request_part[1], "/")
	data, err := os.ReadFile("resources/" + fileName)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}

	responseStr := string(data)
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n%s", len(responseStr), responseStr)))
}
