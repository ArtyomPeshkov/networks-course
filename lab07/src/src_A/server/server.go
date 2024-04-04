package main

import (
	"math/rand"
	"net"
	"strings"
)

func main() {
	conn, err := net.ListenPacket("udp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		if rand.Intn(10) < 2 {
			continue
		}

		_, err = conn.WriteTo([]byte(strings.ToUpper(string(buf[:n]))), addr)
		if err != nil {
			panic(err)
		}
	}
}
