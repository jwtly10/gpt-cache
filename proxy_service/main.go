package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

func main() {
	client := &DefaultClient{}
	baseUrl := "https://api.openai.com/v1"

	parser := NewOpenAiParser()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	log.Printf("TESTING. Flushing redis db.")
	err := rdb.FlushDB().Err()
	if err != nil {
		log.Fatalf("Failed to flush redis db: %v", err)
		return
	}

	cache := NewRedisCache(rdb)

	IndexClient := NewIndexClient("http://localhost:8000", client)

	cacheService := NewCacheService(cache, parser, *IndexClient)

	proxyHandler := NewProxy(client, baseUrl, *cacheService)

	http.HandleFunc("/", proxyHandler.Handle)

	log.Println("Proxy server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
