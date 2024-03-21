package main

import (
	"fmt"
	"net"
)

func main() {
	pc, err := net.ListenPacket("udp4", ":8080")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Time: %s\n", buf[:n])
	}
}
