package main

import (
	"context"
	"fmt"
	"llm_engineering/scrapper"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
	apiKey := os.Getenv("OPENAI_API_KEY")

	switch {
	case apiKey == "":
		panic("No API key was found")
	case !strings.HasPrefix(apiKey, "sk-proj-"):
		panic("Invalid API key format")
	case strings.TrimSpace(apiKey) != apiKey:
		panic("An API key was found, but it looks like it might have space or tab characters at the start or end - please remove them - see troubleshooting notebook")
	default:
		fmt.Println("API key found and looks good so far!")
	}

	message := "Hello, GPT! This is my first ever message to you! Hi!"
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(message),
	}
	ctx := context.Background()
	client := openai.NewClient(option.WithAPIKey(apiKey))
	response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:           "gpt-5.4-mini",
		ReasoningEffort: "none",
		Messages:        messages,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Choices[0].Message.Content)

	systemPrompt := `You are a snarky assistant that analyzes the contents of a website,
and provides a short, snarky, humorous summary, ignoring text that might be navigation related.
Respond in markdown. Do not wrap the markdown in a code block - respond just with the markdown.`

	userPrompt := `Here are the contents of a website.
Provide a short summary of this website.
If it includes news or announcements, then summarize these too.`
	websiteContents, err := scrapper.FetchWebsiteContents("https://edwarddonner.com")
	if err != nil {
		panic(err)
	}

	response, err = client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:           "gpt-5.4-mini",
		ReasoningEffort: "none",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt + "\n\n" + websiteContents),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Choices[0].Message.Content)

}
