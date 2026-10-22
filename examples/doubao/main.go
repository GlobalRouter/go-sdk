package main

import (
	"context"
	"log"
	"os"

	globalrouter "github.com/GlobalRouter/go-sdk"
)

func main() {
	client := globalrouter.New(globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY")))
	ctx := context.Background()
	limit := 1024
	req := globalrouter.ChatRequest{
		Model: "doubao-seed-2-1-pro-260915", MaxCompletionTokens: &limit,
		Messages: []globalrouter.Message{{Role: globalrouter.RoleUser, Content: "What is 17 times 19?"}},
	}
	first, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	if len(first.Choices) == 0 || first.Choices[0].Message == nil {
		log.Fatal("missing assistant message")
	}
	// 原样保存 assistant（含 encrypted_content 和 tool_calls），不要拼接加密字符串。
	req.Messages = append(req.Messages, *first.Choices[0].Message,
		globalrouter.Message{Role: globalrouter.RoleUser, Content: "Add 7 to that result."})
	second, err := client.Chat.Create(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	if len(second.Choices) == 0 || second.Choices[0].Message == nil {
		log.Fatal("missing assistant message")
	}
	log.Println(second.Choices[0].Message.Content)
}
