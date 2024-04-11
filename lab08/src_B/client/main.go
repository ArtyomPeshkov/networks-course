package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func send(conn net.Conn) {
	file, err := os.ReadFile("test.txt")
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
		pack[1] = byte('s')
		iter += 1
		fmt.Println(string(pack))

		taken := min(16, rest)
		pack = append(pack, file[pos:pos+taken]...)
		rest -= taken
		pos += taken

		if rand.Float64() >= 0.3 {
			_, err = conn.Write(pack)
			if err != nil {
				continue
			}
		}

		ack := make([]byte, 1)
		for {
			err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			_, err = conn.Read(ack)
			if err != nil || ack[0] != pack[0] {
				_, err = conn.Write(pack)
				if err != nil {
					continue
				}
			} else {
				break
			}
		}
	}
}

func recv(conn net.Conn) {

	os.Remove("server_test.txt")
	output, err := os.OpenFile("server_test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	for {
		pack := []byte{0, 'r'}

		if rand.Float64() >= 0.3 {
			_, err := conn.Write(pack)
			if err != nil {
				continue
			}
		}

		ack := make([]byte, 1)
		for {
			err := conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			_, err = conn.Read(ack)
			if err != nil || ack[0] != pack[0] {
				_, err = conn.Write(pack)
				if err != nil {
					continue
				}
			} else {
				break
			}
		}
		break
	}
	conn.SetReadDeadline(time.Time{})
	fmt.Println("started")
	for {
		pack := make([]byte, 18)
		n, err := conn.Read(pack)
		if err != nil {
			fmt.Println("Packet lost ")
			fmt.Println(err)
			continue
		}
		pack = pack[:n]

		if rand.Float64() < 0.3 {
			fmt.Println("Lost")
			continue
		}

		_, err = conn.Write([]byte{pack[0]})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("recieved")

		if pack[1] == 'r' {
			output.WriteString(string(pack[2:]))
		} else if pack[1] == 'l' {
			output.WriteString(string(pack[2:]))
			break
		} else {
			break
		}
	}
}

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if os.Args[1] == "send" {
		send(conn)
	} else {
		recv(conn)
	}
}
