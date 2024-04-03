package main

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

type OpenAIParser struct {
}

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
