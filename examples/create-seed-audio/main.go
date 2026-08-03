package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	globalrouter "github.com/GlobalRouter/go-sdk"
)

func main() {
	opts := []globalrouter.SDKOption{
		globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY")),
	}
	if baseURL := os.Getenv("GLOBALROUTER_BASE_URL"); baseURL != "" {
		opts = append(opts, globalrouter.WithBaseURL(baseURL))
	}
	client := globalrouter.New(opts...)
	response, err := client.Audio.CreateSeedAudio(
		context.Background(),
		realRequestBody(),
	)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}

func realRequestBody() globalrouter.SeedAudioRequest {
	return globalrouter.SeedAudioRequest{
		Model:      "doubao-seed-audio-1-0",
		TextPrompt: "Generate a calm piano passage with a soft, relaxing atmosphere.",
		AudioConfig: &globalrouter.SeedAudioConfig{
			Format:         "mp3",
			EnableSubtitle: globalrouter.Bool(true),
		},
	}
}
