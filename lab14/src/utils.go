package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("JSON: %s\n", string(js))

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func sendImage(w http.ResponseWriter, iconPath string) {
	fileBytes, err := os.ReadFile(iconPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}
