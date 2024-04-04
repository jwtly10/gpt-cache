package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type IndexConfig struct {
	BaseUrl string
	Client  Client
}

type CacheServiceConfig struct {
	Threshold float32
}

type ProxyConfig struct {
	Client    Client
	TargetUrl string
}

type Config struct {
	Redis RedisConfig
	Index IndexConfig
	Cache CacheServiceConfig
	Proxy ProxyConfig
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	defaultConfig := &Config{
		Cache: CacheServiceConfig{
			Threshold: 0.2,
		},
		Redis: RedisConfig{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		},
		Index: IndexConfig{
			BaseUrl: "http://localhost:8000",
			Client:  &DefaultClient{},
		},
		Proxy: ProxyConfig{
			TargetUrl: "https://api.openai.com/v1",
			Client:    &DefaultClient{},
		},
	}

	// Override defaults if environment variables are set
	if thresholdStr := os.Getenv("CACHE_THRESHOLD"); thresholdStr != "" {
		if threshold, err := strconv.ParseFloat(thresholdStr, 32); err == nil {
			defaultConfig.Cache.Threshold = float32(threshold)
		} else {
			log.Printf("Failed to parse CACHE_THRESHOLD: %v, using default", err)
		}
	}

	if host := os.Getenv("REDIS_HOST_URL"); host != "" {
		if port := os.Getenv("REDIS_PORT"); port != "" {
			defaultConfig.Redis.Addr = fmt.Sprintf("%s:%s", host, port)
		}
	}

	if pw := os.Getenv("REDIS_PASSWORD"); pw != "" {
		defaultConfig.Redis.Password = pw
	}

	if db := os.Getenv("REDIS_DB"); db != "" {
		if dbInt, err := strconv.Atoi(db); err == nil {
			defaultConfig.Redis.DB = dbInt
		} else {
			log.Printf("Failed to parse REDIS_DB: %v, using default", err)
		}
	}

	if host := os.Getenv("INDEX_HOST_URL"); host != "" {
		if port := os.Getenv("INDEX_PORT"); port != "" {
			defaultConfig.Index.BaseUrl = fmt.Sprintf("%s:%s", host, port)
		}
	}

	return defaultConfig
}
