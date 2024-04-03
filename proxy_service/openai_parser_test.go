package main

import "testing"

func TestOpenAiParser(t *testing.T) {
	parser := NewOpenAiParser()

	reqBody := `{"model":"gpt-3.5-turbo","messages":[{"role":"system","content":"You are a helpful assistant."},{"role":"user","content":"Hello!"}]}`

	messages, err := parser.Parse(reqBody)
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	if messages[0] != "You are a helpful assistant." {
		t.Errorf("Expected 'You are a helpful assistant.', got %s", messages[0])
	}

	if messages[1] != "Hello!" {
		t.Errorf("Expected 'Hello!', got %s", messages[1])
	}
}
