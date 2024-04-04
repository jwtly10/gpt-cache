package main

import (
	"context"
	"log"
	"strings"
)

// CacheService represents a service that interacts with a cache, parser, and index client.
type CacheService struct {
	cache       Cache
	parser      Parser
	indexClient IndexClient
	config      CacheServiceConfig
}

// NewCacheService creates a new instance of CacheService with the provided cache, parser, and index client.
func NewCacheService(cache Cache, parser Parser, indexClient IndexClient, config CacheServiceConfig) *CacheService {
	return &CacheService{
		cache:       cache,
		parser:      parser,
		indexClient: indexClient,
		config:      config,
	}
}

// CacheService.Parse parses the request body using the parser.
// It returns a slice of strings and an error.
func (c *CacheService) Parse(reqBody string) ([]string, error) {
	return c.parser.Parse(reqBody)
}

// CacheService.Get retrieves the cached response for the given context.
// It first queries the cache key from the index service using the provided request context.
// If the cache key is found, it retrieves the cached response from the cache.
// If the cache key is not found or an error occurs during the process, it returns an empty byte slice and the error.
// Otherwise, it returns the cached response as a byte slice and nil error.
func (c *CacheService) Get(ctx context.Context, context []string) ([]byte, error) {
	// data should already be in a nice format

	// Step 1 - Get the cache key from index service
	// TODO: make threshold configurable via injected service config
	idRes, err := c.indexClient.QueryIndex(strings.Join(context, " "), c.config.Threshold)
	if err != nil {
		return []byte{}, err
	}

	// If its empty then return
	if (idRes == QueryResponse{}) {
		return []byte{}, nil
	}

	log.Printf("Response from index service with ID: %d and Distance: %f", idRes.Id, idRes.Distance)

	// TODO - Add a check for distance threshold
	// Here we will need some validation or perhaps some additional logic regarding the distance threshold

	// Step 2 - Get the cached response from cache
	resp, err := c.cache.Get(ctx, idRes.Id)
	if err != nil {
		return []byte{}, nil
	}

	return resp, nil
}

// CacheService.Save saves the response to the cache and adds the index to the index service.
// It takes the request context, context, and response as input.
// It returns an error if any.
func (c *CacheService) Save(ctx context.Context, context []string, response []byte) error {
	// We have now just had a cache miss so we need to add the request to index, and then cache the response

	// Step 1 - Save response to redis and get new ID
	id, err := c.cache.Save(ctx, response)
	if err != nil {
		return err
	}

	log.Printf("Saved response to cache with ID: %d", id)

	// Step 2 - Add the index to the index service with original request
	err = c.indexClient.AddIndex(id, strings.Join(context, " "))
	if err != nil {
		return err
	}

	return nil
}
