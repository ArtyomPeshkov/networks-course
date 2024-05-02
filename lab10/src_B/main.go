package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	addr, err := net.ResolveIPAddr("ip4", os.Args[1])
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 2 && os.Args[2] == "un" {
		addr.IP = net.IP{64, 232, 161, 102}
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var rttMin, rttMax, rttSum time.Duration
	seqNum, recv := 0, 0
	for {
		start := time.Now()

		packet := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   1,
				Seq:  seqNum,
				Data: []byte(""),
			},
		}
		data, _ := (packet).Marshal(nil)
		_, err = conn.WriteTo(data, addr)
		if err != nil {
			panic(err)
		}
		seqNum += 1

		sum := 0
		sumData, _ := packet.Body.Marshal(0)
		for _, elem := range sumData {
			sum += int(elem)
		}

		resp := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(time.Second))
		_, _, err := conn.ReadFrom(resp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		responseData, err := icmp.ParseMessage(1, resp)
		if err != nil {
			panic(err)
		}

		if sum+responseData.Checksum != 65535 && packet.Checksum != 0 {
			fmt.Println("Checksum error")
			continue
		} else if responseData.Checksum == 0 {
			fmt.Println("Checksum ignored")
		}

		recv++
		rtt := time.Since(start)
		if recv == 1 {
			rttMin = rtt
		}

		rttMin = min(rttMin, rtt)
		rttMax = max(rttMax, rtt)
		rttSum += rtt

		fmt.Printf("Echo: Loss=%.f%%, Min=%v, Max=%v, Avg=%v\n", float64(seqNum-recv)*100/float64(seqNum), rttMin, rttMax, rttSum/time.Duration(recv))

		time.Sleep(1 * time.Second)
	}
}
