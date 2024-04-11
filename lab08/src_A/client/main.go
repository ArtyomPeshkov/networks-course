package main

import (
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

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

		pack := make([]byte, 1)
		pack[0] = byte(iter % 2)
		iter += 1

		taken := min(15, rest)
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
