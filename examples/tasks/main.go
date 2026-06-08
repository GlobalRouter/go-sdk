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

	task, err := client.Videos.Generate(context.Background(), globalrouter.GenerationRequest{
		Model:  "doubao-seedance-1-0-lite-t2v",
		Prompt: "A calm product demo shot of a model routing console.",
	}, globalrouter.WithIdempotencyKey("example-video-001"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("task %s status=%s progress=%.2f", task.ID, task.Status, task.Progress)
}
