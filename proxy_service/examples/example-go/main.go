package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))

	config.BaseURL = "http://localhost:8080"
	c := openai.NewClientWithConfig(config)

	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 300,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Show me how to write to file in python",
				},
			},
		},
	)

	if err != nil {
		log.Fatalf("Error making OpenAI Request: %v", err)
	}

	log.Printf("Response: %v\n", resp.Choices[0].Message.Content)
}
