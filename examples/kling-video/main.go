package main

import (
	"context"
	"fmt"
	"log"
	"os"

	globalrouter "github.com/GlobalRouter/go-sdk"
)

func main() {
	model := os.Getenv("GLOBALROUTER_KLING_MODEL")
	if model == "" {
		log.Fatal("set GLOBALROUTER_KLING_MODEL to the configured Kling 4.0 model ID")
	}

	ctx := context.Background()
	client := globalrouter.New(
		globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY")),
	)

	task, err := client.Videos.Generate(ctx, globalrouter.GenerationRequest{
		Model:  model,
		Prompt: "A slow camera move across a clean product workspace.",
		Input: map[string]any{
			"duration_seconds": 10,
			"resolution":       "1080p",
			"aspect_ratio":     "adaptive",
			"references":       []any{},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created task %s with status %s\n", task.ID, task.Status)

	current, err := client.Tasks.Get(ctx, task.ID, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("task %s is %s\n", current.ID, current.Status)
}
