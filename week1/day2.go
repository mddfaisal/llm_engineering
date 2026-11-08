package week1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const (
	OPENAI_URL      = "https://api.openai.com/v1/chat/completions"
	GEMINI_BASE_URL = "https://generativelanguage.googleapis.com/v1beta/openai/"
	OLLAMA_BASE_URL = "http://127.0.0.1:11434/v1"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model           string    `json:"model"`
	ReasoningEffort string    `json:"reasoning_effort,omitempty"` // omitted from JSON if empty
	Messages        []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func Day2() {
	payload := ChatRequest{
		Model:           "gpt-5.4-mini",
		ReasoningEffort: "none",
		Messages: []Message{
			{Role: "user", Content: "Hello, GPT! This is my first ever message to you! Hi!"},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, OPENAI_URL, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+GetEnvKey("OPENAI_API_KEY"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var chat ChatResponse
	if err := json.Unmarshal(responseBody, &chat); err != nil {
		panic(err)
	}
	fmt.Println(chat.Choices[0].Message.Content)

	message := "Tell me a fun fact"
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(message),
	}
	ctx = context.Background()
	openaiClient := openai.NewClient()
	response, err := openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:           "gpt-6-luna",
		ReasoningEffort: "none",
		Messages:        messages,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("### GPT ###")
	fmt.Println(response.Choices[0].Message.Content)

	geminiClient := openai.NewClient(
		option.WithBaseURL(GEMINI_BASE_URL),
		option.WithAPIKey(GetEnvKey("GOOGLE_API_KEY")),
	)

	geminiResponse, err := geminiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:           "gemini-3.1-flash-lite",
		ReasoningEffort: "none",
		Messages:        messages,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("### GEMINI ###")
	fmt.Println(geminiResponse.Choices[0].Message.Content)

	ollamaClient := openai.NewClient(
		option.WithBaseURL(OLLAMA_BASE_URL),
		option.WithAPIKey("ollama"),
	)

	ollamaResponse, err := ollamaClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:           "llama3.2",
		ReasoningEffort: "none",
		Messages:        messages,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("### OLLAMA ###")
	fmt.Println(ollamaResponse.Choices[0].Message.Content)

}
