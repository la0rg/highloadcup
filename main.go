package main

import (
	"log"
	"net/http"
)

func DumbHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Hello, I'm Web Server!"))
}

func main() {
	http.HandleFunc("/", DumbHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalf("FATAL: %v", err)
	}
}
