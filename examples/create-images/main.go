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
	if !exampleutil.RunSelected("seedream_with_reference", &matched) {
		exampleutil.EnsureSelection(matched)
		return
	}

	before := len(capture.Requests)
	response, err := client.Images.Generate(ctx, globalrouter.ImageGenerationRequest{
		Model:  "seedream-image",
		Prompt: "生成一张简洁的模型网关架构图。",
		InputReferences: []globalrouter.ImageTaskReference{{
			Type:     "image_url",
			ImageURL: map[string]any{"url": "https://example.com/reference.png"},
		}},
		Provider: &globalrouter.ProviderSelection{
			ProviderID: "doubao",
			Options: map[string]map[string]any{
				"doubao": {"some_provider_option": "value"},
			},
		},
		Background:        "transparent",
		AspectRatio:       "1:1",
		Resolution:        "2K",
		OutputCompression: globalrouter.Int(90),
		OutputFormat:      "png",
		Quality:           "high",
		Seed:              globalrouter.Int(42),
		Stream:            globalrouter.Bool(false),
		N:                 globalrouter.Int(1),
	})
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.PrintCaptured("Create images (/api/v1/images)", "seedream_with_reference", capture, before, response)
	exampleutil.EnsureSelection(matched)
}
