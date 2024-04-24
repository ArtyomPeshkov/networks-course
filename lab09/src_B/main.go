package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	ip := os.Args[1]
	start, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	end, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	for port := start; port < end; port++ {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, fmt.Sprint(port)), time.Second)
		if err != nil {
			fmt.Println(fmt.Sprint(port) + " closed")
		} else {
			fmt.Println(fmt.Sprint(port) + " opened")
			conn.Close()
		}
	}
}
