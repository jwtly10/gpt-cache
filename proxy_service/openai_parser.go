package main

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

type Parser interface {
	// Parse takes a request body as input and returns a slice of strings and an error.
	// It parses the LLM Request and extracts the content of each message or "context"
	Parse(reqBody string) ([]string, error)
}

type OpenAIParser struct{}

func NewOpenAiParser() *OpenAIParser {
	return &OpenAIParser{}
}

// Parse parses the request body and extracts the content of each message.
// It returns a slice of strings containing the content of each message and an error, if any.
func (r *OpenAIParser) Parse(reqBody string) ([]string, error) {
	var req openai.ChatCompletionRequest
	err := json.Unmarshal([]byte(reqBody), &req)
	if err != nil {
		return []string{}, err
	}

	var messages []string
	for _, message := range req.Messages {
		messages = append(messages, message.Content)
	}

	return messages, nil
}
