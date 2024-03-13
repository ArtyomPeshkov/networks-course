package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func StartProxy(address string, port string) error {
	fmt.Printf("Started proxy on http://%s:%s\n", address, port)

	http.HandleFunc("/", httpHandler)
	return http.ListenAndServe(address+":"+port, nil)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet && req.Method != http.MethodPost {
		fmt.Printf("Only allowed methods are %v and %v, got %v\n", http.MethodGet, http.MethodPost, req.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqCopy, err := http.NewRequest(req.Method, strings.Trim(req.RequestURI, "/"), req.Body)
	if err != nil {
		fmt.Printf("Cloning request failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	reqCopy.Header = req.Header.Clone()

	client := http.Client{}
	response, err := client.Do(reqCopy)
	if err != nil {
		fmt.Printf("Sending request failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Reading body failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for k, v := range response.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(response.StatusCode)
	w.Write(body)

	log.Printf("URL: %v, Status: %v, Type: %v", response.Request.URL, response.Status, reqCopy.Method)

}
