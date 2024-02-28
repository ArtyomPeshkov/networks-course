package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := NewStoreServer()
	mux.HandleFunc("/product/", server.TaskHandler)
	mux.HandleFunc("/products/", server.TaskHandler)

	if err := http.ListenAndServe("localhost:"+"8080", mux); err != http.ErrServerClosed {
		log.Print("Failed to run server at port 8080\nError: " + err.Error() + "\n")
	}
}
