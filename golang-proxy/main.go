package main

import (
	"log"
	"net/http"
)

func main() {
	client := &DefaultClient{}
	baseUrl := "https://api.openai.com/v1"
	proxyHandler := NewProxy(client, baseUrl)

	http.HandleFunc("/", proxyHandler.Handle)

	log.Println("Proxy server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
