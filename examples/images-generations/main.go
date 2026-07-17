package main

import (
	"context"
	"log"

	globalrouter "github.com/GlobalRouter/go-sdk"
	"github.com/GlobalRouter/go-sdk/examples/internal/exampleutil"
)

func main() {
	ctx := context.Background()
	capture := &exampleutil.CaptureClient{}
	client := exampleutil.NewClient(capture)

	matched := false
	if !exampleutil.RunSelected("create_image", &matched) {
		exampleutil.EnsureSelection(matched)
		return
	}

	before := len(capture.Requests)
	response, err := client.Images.Generate(ctx, globalrouter.ImageGenerationRequest{
		Model:  "gpt-image-1",
		Prompt: "A cute baby sea otter",
		N:      globalrouter.Int(1),
		Size:   "1024x1024",
	})
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.PrintCaptured("OpenAI-compatible image generation", "create_image", capture, before, response)
	exampleutil.EnsureSelection(matched)
}
