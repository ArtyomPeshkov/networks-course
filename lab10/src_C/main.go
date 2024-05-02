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

		switch responseData.Type {
		case ipv4.ICMPTypeDestinationUnreachable:
			switch responseData.Code {
			case 0:
				fmt.Println("ICMP Destination Network Unreachable")
				continue
			case 1:
				fmt.Println("ICMP Destination Host Unreachable")
				continue
			case 2:
				fmt.Println("ICMP Protocol Unreachable")
				continue
			case 3:
				fmt.Println("ICMP Port Unreachable")
				continue
			case 4:
				fmt.Println("Frag required")
				continue
			case 5:
				fmt.Println("Incorrect route")
				continue
			case 6:
				fmt.Println("ICMP Network Unknown")
				continue
			case 7:
				fmt.Println("ICMP Host Unknown")
				continue
			default:
				if responseData.Code > 15 {
					fmt.Println("Unknown ICMP Error Code")
				}
			}
		case ipv4.ICMPTypeTimeExceeded:
			switch responseData.Code {
			case 0:
				fmt.Println("ICMP TL")
				continue
			case 1:
				fmt.Println("ICMP Building Fragment Time Exceeded")
				continue
			default:
				fmt.Println("Unknown ICMP Error Code")
			}

		case ipv4.ICMPTypeParameterProblem:
			switch responseData.Code {
			case 0:
				fmt.Println("ICMP Pointer points to error")
				continue
			case 1:
				fmt.Println("ICMP Absence of reauired option")
				continue
			case 2:
				fmt.Println("ICMP Wrong Length")
				continue
			default:
				fmt.Println("Unknown ICMP Error Code")
			}
		case ipv4.ICMPTypeRedirect:
			switch responseData.Code {
			case 0:
				fmt.Println("Depricated")
				continue
			case 1:
				fmt.Println("ICMP Host redirection")
				continue
			case 2:
				fmt.Println("Depricated")
				continue
			case 3:
				fmt.Println("ICMP Host redirection according to type")
				continue
			default:
				fmt.Println("Unknown ICMP Error Code")
			}
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
