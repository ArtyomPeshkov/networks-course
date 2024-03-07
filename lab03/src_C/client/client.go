package main

import (
	"fmt"
	"net"
)

func runClient(host string, port string, filename string) {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\n\r\n", filename, host+":"+port)
	_, err = conn.Write([]byte(request))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}

	response := string(buf)
	fmt.Println(response)
}
