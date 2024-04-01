package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type Proxy struct {
	client    Client
	targetURL string
}

func NewProxy(client Client, targetURL string) *Proxy {
	return &Proxy{
		client:    client,
		targetURL: targetURL,
	}
}

func (p *Proxy) Handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("Target URL: %s", p.targetURL)
	log.Printf("Original request: %s %s", r.Method, r.URL.String())

	// Logging original headers
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("Request header: %s=%s", name, value)
		}
	}

	// Logging original body
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "go-gpt-cache: Failed to read original request body", http.StatusInternalServerError)
		return
	}
	log.Printf("Original Request body: %s", string(body))

	// IMPORTANT: Restore the original body for further use
	// This is required as io.ReadAll consumes the body, and it can't be read again
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	newRequestBody := bytes.NewBuffer(body)

	newReq, err := http.NewRequest(r.Method, p.targetURL+r.URL.String(), newRequestBody)
	if err != nil {
		http.Error(w, "go-gpt-cache: Failed to create new request", http.StatusInternalServerError)
		return
	}

	copyHeaders(r.Header, newReq.Header)
	resp, err := p.client.Do(newReq)
	if err != nil {
		http.Error(w, "go-gpt-cache: Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeaders(src http.Header, dst http.Header) {
	for name, values := range src {
		for _, value := range values {
			dst.Add(name, value)
		}
	}
}
