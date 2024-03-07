package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var activeThreads int32 = 0

func runServer(port string, concurencyLevelStr string) {
	queue := make(chan net.Conn)
	concurencyLevel, err := strconv.Atoi(concurencyLevelStr)
	if err != nil {
		log.Printf("wrong concurency level")
		concurencyLevel = 10000
	}

	go connectionQueueHandler(queue, int32(concurencyLevel))

	listen, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		queue <- conn
	}
}

func connectionQueueHandler(queue chan net.Conn, concurencyLevel int32) {
	for conn := range queue {
		for {
			if atomic.LoadInt32(&activeThreads) <= concurencyLevel {
				atomic.AddInt32(&activeThreads, 1)
				go handler(conn)
				break
			}
			time.Sleep(time.Second)
		}
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	request_storage := make([]byte, 1024)
	n, err := conn.Read(request_storage)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading from socet, ", err)
		atomic.StoreInt32(&activeThreads, atomic.LoadInt32(&activeThreads)-1)
		return
	}
	request_storage = request_storage[:n]
	request_part := strings.Split(string(request_storage), " ")
	if len(request_part) < 2 {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return
	}
	fileName := strings.Trim(request_part[1], "/")
	data, err := os.ReadFile("resources/" + fileName)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		atomic.StoreInt32(&activeThreads, atomic.LoadInt32(&activeThreads)-1)
		return
	}

	responseStr := string(data)
	conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n%s", len(responseStr), responseStr)))
	atomic.StoreInt32(&activeThreads, atomic.LoadInt32(&activeThreads)-1)
}
