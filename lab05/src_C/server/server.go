package main

import (
	"net"
	"time"
)

func main() {
	pc, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.255:8080")
	if err != nil {
		panic(err)
	}

	for {
		_, err = pc.WriteTo([]byte(time.Now().Format(time.RFC3339)), addr)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
