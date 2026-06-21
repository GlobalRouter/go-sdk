# GlobalRouter Go SDK

Go SDK for the GlobalRouter public API.

GlobalRouter provides one interface for model access, routing, transparent
pricing, async multimodal tasks, billing, logs, and administration. This SDK
focuses on the public model, chat, embeddings, async task, image, audio, video,
and 3D generation APIs.

## Installation

```bash
go get github.com/GlobalRouter/go-sdk
```

## Usage

```go
package main

import (
	"context"
	"log"
	"os"

	globalrouter "github.com/GlobalRouter/go-sdk"
)

func main() {
	ctx := context.Background()
	client := globalrouter.New(
		globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY")),
	)

	res, err := client.Chat.Create(ctx, globalrouter.ChatRequest{
		Model: "openai/gpt-4o-mini",
		Messages: []globalrouter.Message{{
			Role:    globalrouter.RoleUser,
			Content: "Say hello in one sentence.",
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Choices[0].Message.Content)
}
```

## Streaming

```go
stream, err := client.Chat.Stream(ctx, globalrouter.ChatRequest{
	Model: "openai/gpt-4o-mini",
	Messages: []globalrouter.Message{{
		Role:    globalrouter.RoleUser,
		Content: "Stream a short answer.",
	}},
})
if err != nil {
	log.Fatal(err)
}
defer stream.Close()

for {
	event, err := stream.Next()
	if errors.Is(err, io.EOF) {
		break
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println(event.Data.Choices)
}
```

## Resources

- `client.Models.List`
- `client.Chat.Create`
- `client.Chat.Stream`
- `client.Embeddings.Create`
- `client.Tasks.Create`
- `client.Tasks.List`
- `client.Tasks.Get`
- `client.Tasks.CreateBatch`
- `client.Tasks.GetBatch`
- `client.Tasks.Events`
- `client.Tasks.EventsMany`
- `client.Tasks.Cancel`
- `client.Tasks.Retry`
- `client.Images.Generate`
- `client.Audio.CreateSpeech`
- `client.Audio.CreateTranscription`
- `client.Videos.Generate`
- `client.ThreeD.Generate`

## Configuration

```go
client := globalrouter.New(
	globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY")),
	globalrouter.WithBaseURL("http://127.0.0.1:8000"),
	globalrouter.WithRetryConfig(globalrouter.RetryConfig{MaxRetries: 2}),
)
```

The SDK reads `GLOBALROUTER_API_KEY` by default when `WithAPIKey` is omitted.

Per-request options:

- `WithIdempotencyKey`
- `WithHeader`
- `WithRequestTimeout`
- `WithRequestRetries`

## Error Handling

GlobalRouter error envelopes are returned as `*globalrouter.APIError`.

```go
var apiErr *globalrouter.APIError
if errors.As(err, &apiErr) {
	log.Printf("request %s failed with %s: %s", apiErr.RequestID, apiErr.Code, apiErr.Message)
}
```

## Webhook Signatures

```go
ok := globalrouter.VerifyWebhookSignature(secret, payload, signature)
```

The verifier supports both `sha256=<hex>` and `t=<timestamp>,v1=<hex>` formats. Timestamped signatures must be within the default 5 minute tolerance. Legacy `sha256=<hex>` signatures do not include a timestamp and are verified for compatibility without a freshness check.

## Development

```bash
go test ./...
```
