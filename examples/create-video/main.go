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
	if !exampleutil.RunSelected("seedance_video", &matched) {
		exampleutil.EnsureSelection(matched)
		return
	}

	before := len(capture.Requests)
	response, err := client.Videos.Generate(ctx, globalrouter.GenerationRequest{
		Model:         "doubao-seedance-1-0-pro-fast-251015",
		Prompt:        "A quiet city rooftop at sunset, slow cinematic push-in",
		AspectRatio:   "16:9",
		Duration:      globalrouter.Int(5),
		Resolution:    "720p",
		SR:            &globalrouter.SuperResolutionRequest{Resolution: "1080p"},
		Seed:          globalrouter.Int(12345),
		GenerateAudio: globalrouter.Bool(true),
		CallbackURL:   "https://example.com/webhooks/video",
		FrameImages: []globalrouter.VideoFrameImage{{
			Type:      "image_url",
			ImageURL:  map[string]any{"url": "https://example.com/assets/opening-frame.png"},
			FrameType: "first_frame",
		}},
		InputReferences: []map[string]any{{
			"type":      "image_url",
			"image_url": map[string]any{"url": "https://example.com/assets/reference.png"},
		}},
		Provider: &globalrouter.ProviderSelection{
			ProviderID: "doubao",
			Options: map[string]map[string]any{
				"doubao": {"some_provider_option": "value"},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.PrintCaptured("Create video (/api/v1/videos)", "seedance_video", capture, before, response)
	exampleutil.EnsureSelection(matched)
}
