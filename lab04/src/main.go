package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("resources/journal.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	flag.Parse()

	if nil != StartProxy("localhost", "8080") {
		panic("failed to start")
	}
}
