package main

import (
	"math/rand"
	"net"
	"os"
)

func main() {
	conn, err := net.ListenPacket("udp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	os.Remove("test.txt")
	output, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	for {
		pack := make([]byte, 16)
		n, addr, err := conn.ReadFrom(pack)
		if err != nil {
			continue
		}
		pack = pack[:n]

		if rand.Float64() < 0.3 {
			continue
		}

		_, err = conn.WriteTo([]byte{pack[0]}, addr)
		if err != nil {
			continue
		}

		output.WriteString(string(pack[1:]))
	}
}
