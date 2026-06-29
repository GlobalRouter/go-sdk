# GlobalRouter Go SDK v1 Design

## Background

GlobalRouter is a multi-tenant AI model aggregation platform with an
OpenAI-compatible public API plus async multimodal task APIs. The Go SDK should
make the public API easy to call from Go services without exposing provider
credentials or admin-only surfaces.

The SDK intentionally references the shape of `OpenRouterTeam/go-sdk`: a
top-level `New(...)` constructor, resource groups, Bearer authentication,
retry configuration, custom HTTP clients, typed request/response models, and
server-sent event streaming. It does not copy the full Speakeasy-generated
surface because GlobalRouter v1 only needs the product's core public API.

## Scope

v1 covers the same core public operations as the in-repo Python helper:

- `GET /v1/models`
- `POST /v1/chat/completions`
- `POST /v1/embeddings`
- `POST /v1/tasks`
- `GET /v1/tasks`
- `POST /v1/tasks/batch`
- `GET /v1/tasks/batch/{batch_id}`
- `GET /v1/tasks/{task_id}`
- `GET /v1/tasks/{task_id}/events`
- `GET /v1/tasks/events`
- `POST /v1/tasks/{task_id}/cancel`
- `POST /v1/tasks/{task_id}/retry`
- `POST /v1/images/generations`
- `POST /v1/audio/speech`
- `POST /v1/audio/transcriptions`
- `POST /v1/videos/generations`
- `POST /v1/3d/generations`

Out of scope for v1: admin APIs, portal APIs, billing management APIs, provider
configuration APIs, and OpenRouter compatibility management endpoints such as
keys, credits, generations, guardrails, workspaces, and organization resources.

## Public Shape

The module path is `github.com/GlobalRouter/go-sdk` and the package name is
`globalrouter`.

Client creation:

```go
client := globalrouter.New(
    globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY")),
    globalrouter.WithBaseURL("https://api.globalrouter.com"),
)
```

Top-level resources:

- `client.Models`
- `client.Chat`
- `client.Embeddings`
- `client.Tasks`
- `client.Images`
- `client.Audio`
- `client.Videos`
- `client.ThreeD`

The SDK provides typed structs for stable GlobalRouter fields and uses
`map[string]any` for intentionally extensible payloads such as task input,
metadata, routing, response mode, provider-specific results, and multimodal
provider responses.

## Request Flow

All resource methods call the shared request helper. The helper:

1. Builds URLs from the configured base URL and resource path.
2. Adds `Authorization: Bearer <key>` when an API key is configured.
3. Adds `Content-Type: application/json` for request bodies.
4. Adds a GlobalRouter user agent.
5. Encodes request bodies as JSON.
6. Retries transport errors and HTTP 5xx responses according to `RetryConfig`.
7. Parses GlobalRouter error envelopes into `*APIError`.
8. Decodes JSON responses into typed result structs or `map[string]any`.

Per-request options support idempotency keys, headers, timeout overrides, and
retry overrides.

## Streaming

Chat streaming and task events use the same SSE parser. The parser preserves the
SSE `event:` name and decodes each `data:` payload into a typed `StreamEvent[T]`.
`data: [DONE]` terminates the stream with `io.EOF`, matching Go iterator
conventions.

## Error Handling

HTTP 4xx and final HTTP 5xx responses return `*APIError`:

```go
type APIError struct {
    StatusCode int
    Code       string
    Message    string
    Type       string
    RequestID  string
    Body       string
}
```

This mirrors the GlobalRouter unified error envelope and keeps request IDs
available for debugging and support.

## Verification

The test suite uses `httptest` and does not depend on live provider keys. It
verifies:

- Bearer authentication and JSON request encoding.
- Model filter query parameters.
- Error envelope parsing.
- Retry behavior for HTTP 5xx responses.
- SSE parsing for chat and task events.
- Task, image, audio, video, 3D, and embeddings endpoint paths.
- Idempotency headers.
- Webhook HMAC signature verification.

Live provider verification is intentionally not claimed by SDK unit tests.
