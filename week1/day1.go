package week1

import (
	"context"
	"fmt"
	"llm_engineering/scrapper"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func GetEnvKey(key string) string {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
	return os.Getenv(key)
}

func Day1() {

	message := "Hello, GPT! This is my first ever message to you! Hi!"
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(message),
	}
	ctx := context.Background()
	client := openai.NewClient(option.WithAPIKey(GetEnvKey("OPENAI_API_KEY")))
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
