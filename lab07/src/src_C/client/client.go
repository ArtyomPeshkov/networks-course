package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	minRtt := 1.
	maxRtt := 0.
	avgRtt := 0.
	loss := 0
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
			loss += 1
			fmt.Printf("Ping %d\n", i+1)
			fmt.Println("Request timed out")
			fmt.Println("")

			continue
		}

		fmt.Println(string(buf[:n]))
		rtt := time.Since(start).Seconds()
		minRtt = min(minRtt, rtt)
		maxRtt = max(maxRtt, rtt)
		avgRtt += rtt
		fmt.Printf("rtt min/avg/max = %f/%f/%f seconds\n", minRtt, avgRtt/(float64(i+1-loss)), maxRtt)
		fmt.Println("")
	}
	fmt.Printf("%d%% packet loss", loss*10)
}
