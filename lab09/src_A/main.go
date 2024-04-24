package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	fmt.Println(localAddress.IP)
	for i, bt := range localAddress.IP.DefaultMask() {
		fmt.Print(strconv.Itoa(int(bt)))
		if i < 3 {
			fmt.Print(".")
		}
	}
	fmt.Println()
}
