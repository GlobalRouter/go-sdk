package main

import (
	"context"
	"log"
	"os"

	globalrouter "github.com/GlobalRouter/go-sdk"
)

func main() {
	client := globalrouter.New(
		globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY")),
	)

	res, err := client.Chat.Create(context.Background(), globalrouter.ChatRequest{
		Model: "openai/gpt-4o-mini",
		Messages: []globalrouter.Message{{
			Role:    globalrouter.RoleUser,
			Content: "Explain GlobalRouter in one sentence.",
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(res.Choices) > 0 && res.Choices[0].Message != nil {
		log.Println(res.Choices[0].Message.Content)
	}
}
