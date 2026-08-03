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
	client := globalrouter.New(
		globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY")),
	)
	response, err := client.Audio.CreateSeedAudio(
		context.Background(),
		globalrouter.SeedAudioRequest{
			Model:      "doubao-seed-audio-1-0",
			TextPrompt: "Use @音频1 as a style reference for a calm piano passage",
			References: []globalrouter.SeedAudioReference{{
				AudioURL: "https://example.com/reference.mp3",
			}},
			AudioConfig: &globalrouter.SeedAudioConfig{
				Format:         "mp3",
				EnableSubtitle: globalrouter.Bool(true),
			},
		},
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
