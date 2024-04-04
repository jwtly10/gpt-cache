package main

import (
	"log"
	"net/http"
)

func main() {

	config := NewConfig()
	log.Printf("Config Setup: %+v", config)

	parser := NewOpenAiParser()

	// rdb := redis.NewClient(&redis.Options{
	// 	// Addr:     "127.0.0.1:6379",
	// 	Addr:     "redis:6379",
	// 	Password: "",
	// 	DB:       0,
	// })

	rdb := NewRedisClient(config.Redis)

	log.Printf("TESTING. Flushing redis db.")
	err := rdb.FlushDB().Err()
	if err != nil {
		log.Fatalf("Failed to flush redis db: %v", err)
		return
	}

	cache := NewRedisCache(rdb)

	IndexClient := NewIndexClient(config.Index)

	cacheService := NewCacheService(cache, parser, *IndexClient, config.Cache)

	proxyHandler := NewProxy(config.Proxy, *cacheService)

	http.HandleFunc("/", proxyHandler.Handle)

	log.Println("Proxy server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
