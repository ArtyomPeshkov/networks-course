package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a command: ")
	command, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading command from user:", err)
		return
	}
	command = strings.Trim(command, "\n")
	command = command[:len(command)-1]
	conn.Write([]byte(command))

	output := make([]byte, 1024)
	n, err := conn.Read(output)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading output:", err)
		return
	}

	fmt.Println(string(output[:n]))
}
