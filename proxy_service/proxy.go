package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Proxy struct {
	config       ProxyConfig
	cacheService CacheService
}

func NewProxy(config ProxyConfig, cacheService CacheService) *Proxy {
	return &Proxy{
		config:       config,
		cacheService: cacheService,
	}
}

func (p *Proxy) Handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("Target URL: %s", p.config.TargetUrl)
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
		http.Error(w, "gpt-cache: Failed to read original request body", http.StatusInternalServerError)
		return
	}
	log.Printf("Original Request body: %s", string(body))

	// Clean up the request body, for better indexing
	context, err := p.cacheService.Parse(string(body))
	if err != nil {
		log.Print("Failed to parse request body")
	}

	// Successfully parsed the request body
	if len(context) != 0 {
		// Now try to find cache hit
		cache, err := p.cacheService.Get(r.Context(), context)
		if err != nil {
			log.Printf("Cache miss due to error: %v", err)
		}

		if len(cache) > 0 {
			log.Printf("Cache hit, returning cached response")

			// Emulate response headers
			// TODO: Undersand this more, should they be set based off of the original request?
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(cache)))
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Accept-Encoding", "gzip")

			w.Header().Set("Authorization", r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

			// Write response body
			var buf bytes.Buffer
			tee := io.TeeReader(bytes.NewReader(cache), &buf)
			if _, err := io.Copy(w, tee); err != nil {
				log.Printf("Failed to write response body to client: %v", err)
			}

			return
		}

	}

	log.Println("Cache miss, forwarding request to target")
	log.Printf("Forwarding request to %s", p.config.TargetUrl+r.URL.String())

	// IMPORTANT: Restore the original body for further use
	// This is required as io.ReadAll consumes the body, and it can't be read again
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	newRequestBody := bytes.NewBuffer(body)

	newReq, err := http.NewRequest(r.Method, p.config.TargetUrl+r.URL.String(), newRequestBody)
	if err != nil {
		http.Error(w, "gpt-cache: Failed to create new request", http.StatusInternalServerError)
		return
	}

	copyHeaders(r.Header, newReq.Header)
	resp, err := p.config.Client.Do(newReq)
	if err != nil {
		http.Error(w, "gpt-cache: Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Immediately copy headers and status code to the response writer
	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)

	// Stream the body directly to the client and capture it for caching
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)
	if _, err := io.Copy(w, tee); err != nil {
		log.Printf("Failed to write response body to client: %v", err)
		// Since we've already started sending the response, we can't send a new HTTP status code to the client
	}

	// Cache the response body
	if err := p.cacheService.Save(r.Context(), context, buf.Bytes()); err != nil {
		log.Printf("Failed to save response to cache: %v", err)
	}
}

func copyHeaders(src http.Header, dst http.Header) {
	for name, values := range src {
		for _, value := range values {
			dst.Add(name, value)
		}
	}
}
