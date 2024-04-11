package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func send(conn net.PacketConn, addr net.Addr) {
	file, err := os.ReadFile("server_test.txt")
	if err != nil {
		panic(err)
	}

	rest := len(file)
	pos := 0
	iter := 0

	for {
		if rest == 0 {
			break
		}

		pack := make([]byte, 2)
		pack[0] = byte(iter % 2)
		pack[1] = byte('r')
		iter += 1
		fmt.Println(string(pack))

		taken := min(16, rest)
		pack = append(pack, file[pos:pos+taken]...)
		rest -= taken
		pos += taken

		if rest == 0 {
			pack[1] = byte('l')
		}

		if rand.Float64() >= 0.3 {
			_, err = conn.WriteTo(pack, addr)
			if err != nil {
				continue
			}
		}

		ack := make([]byte, 1)
		for {
			err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			_, _, err = conn.ReadFrom(ack)
			if err != nil || ack[0] != pack[0] {
				_, err = conn.WriteTo(pack, addr)
				fmt.Println("Lost")
			} else {
				break
			}
		}
	}
	conn.SetReadDeadline(time.Time{})
}

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
		fmt.Println("Ready to get pack")
		pack := make([]byte, 18)
		n, addr, err := conn.ReadFrom(pack)
		if err != nil {
			continue
		}
		fmt.Println("Got pack")
		println(pack[:n])
		pack = pack[:n]

		if rand.Float64() < 0.3 {
			continue
		}

		_, err = conn.WriteTo([]byte{pack[0]}, addr)
		if err != nil {
			continue
		}

		if pack[1] == 's' {
			output.WriteString(string(pack[2:]))
		} else if pack[1] == 'r' {
			send(conn, addr)
			fmt.Println("sent data")
		}
	}
}
