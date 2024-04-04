package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func makeRequest(c *openai.Client, message string) (time.Duration, error) {
	start := time.Now() // Start timing before the request

	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 300,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	elapsed := time.Since(start) // Calculate the elapsed time after the request

	if err != nil {
		log.Printf("Error making OpenAI Request: %v", err)
		return elapsed, err
	}

	// Commented for brevity
	// log.Printf("Response: %v\n", resp.Choices[0].Message.Content)
	log.Printf("Response Length: %d\n", len(resp.Choices[0].Message.Content))
	log.Printf("Request took: %s", elapsed)

	return elapsed, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
	config.BaseURL = "http://localhost:8080"
	c := openai.NewClientWithConfig(config)

	var benchmarks []time.Duration

	log.Println("Making the first request (No Cache hit) ...")
	message := "How do you write to file in Java"
	fmt.Println("Asking ChatGPT: ", message)
	res1, _ := makeRequest(c, message)
	benchmarks = append(benchmarks, res1)

	fmt.Println()

	log.Println("Making the second request (Cache hit) ...")
	fmt.Println("Asking ChatGPT: ", message)
	message = "Show me how to write to file in Java"
	res2, _ := makeRequest(c, message)
	benchmarks = append(benchmarks, res2)

	fmt.Println()

	log.Println("****************** Benchmark Results ******************")
	for i, benchmark := range benchmarks {
		log.Printf("Request %d took: %s\n", i+1, benchmark)
	}
}
