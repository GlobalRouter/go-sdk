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
			TextPrompt: "Warm acoustic guitar and soft piano, calm, instrumental",
			AudioConfig: &globalrouter.SeedAudioConfig{
				Format: "mp3",
			},
			Watermark: globalrouter.Bool(false),
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
