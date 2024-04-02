package main

import "net/http"

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type DefaultClient struct{}

func (c *DefaultClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
