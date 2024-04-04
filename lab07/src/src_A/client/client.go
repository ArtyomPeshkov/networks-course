package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		fmt.Print(">> ")
		var data string
		fmt.Scanln(&data)

		_, err = conn.Write([]byte(data))
		if err != nil {
			continue
		}

		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		fmt.Println("Resp: ", string(buf[:n]))
	}
}
