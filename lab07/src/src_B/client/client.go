package main

import (
	"fmt"
	"net"
	"time"
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

	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := conn.Write([]byte(fmt.Sprintf("Ping %d %s", i+1, start.Format("15:04:05.000"))))
		if err != nil {
			continue
		}

		conn.SetReadDeadline(time.Now().Add(time.Second))
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Ping %d\n", i+1)
			fmt.Println("Request timed out")
			fmt.Println("")

			continue
		}

		fmt.Println(string(buf[:n]))
		fmt.Printf("RTT: %f seconds\n", time.Since(start).Seconds())
		fmt.Println("")
	}
}
