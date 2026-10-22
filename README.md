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
		Model:               "qwen3.8-max",
		MaxCompletionTokens: globalrouter.Int(1024),
		Reasoning: &globalrouter.ReasoningConfig{
			Effort:  "high",
			Exclude: globalrouter.Bool(true),
		},
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

## Chat 内容透传

`ContentPart` 支持 `Video`、`VideoURL`、`InputAudio`、`File`。
通用 Chat 链路按传入的类型和字段透传，上游模型支持情况由上游判断；
原生协议适配器无法转换的类型会明确报错。

```go
content := []globalrouter.ContentPart{
    {Type: "video_url", VideoURL: map[string]any{"url": "https://example.com/video.mp4"}},
    {Type: "input_audio", InputAudio: map[string]any{"data": "<base64>", "format": "wav"}},
    {Type: "file", File: map[string]any{"file_id": "<上游文件 ID>"}},
}
```

将 `content` 赋给 `Message.Content`。任意上游扩展字段可以通过
`Message.Content: []map[string]any{...}` 表达。文件 ID 不会自动跨 provider 搬运。

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

图片生成的 `ImageGenerationRequest.Watermark` 使用 `*bool`，显式传入 `false`
会发送 JSON `false`。对象形式的 `provider`（例如 `provider_id`）使用
`client.Images.GenerateOpenRouter(ctx, map[string]any{...})` 调用 `/api/v1/images`。
两种图片提交均只发送一次；超时或结果未知时应核对原请求。

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

The verifier supports both `sha256=<hex>` and `t=<timestamp>,v1=<hex>` formats.

## Development

```bash
go test ./...
```
