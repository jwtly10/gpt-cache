package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestQueryIndex(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &MockClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				response := QueryResponse{
					Id:       1,
					Distance: 0.5,
				}
				respBody, _ := json.Marshal(response)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBody)),
				}, nil
			},
		}

		config := IndexConfig{
			BaseUrl: "http://example.com",
			Client:  client,
		}

		indexClient := NewIndexClient(config)
		resp, err := indexClient.QueryIndex("test context", 0.5)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if resp.Id != 1 {
			t.Errorf("Expected Id 1, got %v", resp.Id)
		}
	})

	t.Run("no content", func(t *testing.T) {
		client := &MockClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}, nil
			},
		}
		config := IndexConfig{
			BaseUrl: "http://example.com",
			Client:  client,
		}

		indexClient := NewIndexClient(config)
		resp, err := indexClient.QueryIndex("test context", 0.5)
		if err != nil {
			t.Errorf("Expect no error, just empty resp, got %v", err)
		}

		expected := QueryResponse{}

		if resp != expected {
			t.Errorf("Expected empty QueryResponse{}, got %v", resp)
		}
	})

	t.Run("error", func(t *testing.T) {
		client := &MockClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				response := ErrResponse{
					Detail: "Indexing service error",
				}

				respBody, _ := json.Marshal(response)

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader(respBody)),
				}, nil
			},
		}

		config := IndexConfig{
			BaseUrl: "http://example.com",
			Client:  client,
		}

		indexClient := NewIndexClient(config)
		_, err := indexClient.QueryIndex("test context", 0.5)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		expected := "Indexing service error"

		if err.Error() != expected {
			t.Errorf("Expected %v error, got %v", expected, err.Error())
		}

	})
}

func TestAddIndex(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &MockClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				response := AddResponse{
					Status:  "ok",
					Message: "added",
				}
				respBody, _ := json.Marshal(response)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
				}, nil
			},
		}

		config := IndexConfig{
			BaseUrl: "http://example.com",
			Client:  client,
		}

		indexClient := NewIndexClient(config)
		err := indexClient.AddIndex(1, "test context")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
