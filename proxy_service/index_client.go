package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type IndexClient struct {
	config IndexConfig
}

type QueryRequest struct {
	Context           string  `json:"context"`
	DistanceThreshold float32 `json:"distance_threshold"`
}

type AddRequest struct {
	Id      int64  `json:"id"`
	Context string `json:"context"`
}

type QueryResponse struct {
	Id       int64   `json:"id"`
	Distance float32 `json:"distance"`
}

type AddResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ErrResponse struct {
	Detail string `json:"detail"`
}

func NewIndexClient(config IndexConfig) *IndexClient {
	return &IndexClient{
		config: config,
	}
}

// QueryIndex sends a query to the index microservice and returns the query response.
// It takes a context string and a distance threshold as parameters.
// The indexing microservice returns 204 if no match is found.
// If a match is found, it will return 200, and the response will contain the ID and distance of the match.
// If an error occurs during the request or response handling, it returns an empty QueryResponse and the error.
//
// Note: It is important for callers to check both the returned QueryResponse and error. An empty QueryResponse
// with a nil error indicates a successful query with no matches found.
func (i *IndexClient) QueryIndex(context string, distanceThreshold float32) (QueryResponse, error) {
	endpoint := "/queryIndex"
	query := QueryRequest{
		Context:           context,
		DistanceThreshold: distanceThreshold,
	}

	reqBody, err := json.Marshal(query)
	if err != nil {
		return QueryResponse{}, err
	}

	req, err := http.NewRequest("POST", i.config.BaseUrl+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return QueryResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := i.config.Client.Do(req)
	if err != nil {
		return QueryResponse{}, err
	}

	defer resp.Body.Close()

	// The indexing microservice only returns 204 if no match is found
	// If match was found, it will return 200
	// All other cases are some kind of error whihc needs to be handled

	// If no match found, return empty response
	if resp.StatusCode == http.StatusNoContent {
		return QueryResponse{}, nil
	}

	// If not OK, return error
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrResponse

		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			return QueryResponse{}, err
		}

		return QueryResponse{}, errors.New(errorResp.Detail)
	}

	// If OK, return the response
	var queryResp QueryResponse
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	if err != nil {
		return QueryResponse{}, err
	}

	return queryResp, nil
}

// AddIndex adds an index with the specified ID and context to the index service.
// It returns an `AddResponse` containing the result of the operation and an error, if any.
func (i *IndexClient) AddIndex(id int64, context string) error {
	endpoint := "/addIndex"
	query := AddRequest{
		Id:      id,
		Context: context,
	}

	reqBody, err := json.Marshal(query)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", i.config.BaseUrl+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := i.config.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// If not OK, return error
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrResponse

		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		if err != nil {
			return err
		}

		return errors.New(errorResp.Detail)
	}

	var addResp AddResponse
	err = json.NewDecoder(resp.Body).Decode(&addResp)
	if err != nil {
		return err
	}

	return nil
}
