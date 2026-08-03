# GlobalRouter Go SDK examples

These examples follow the public GlobalRouter docs at:

```text
https://test-global-router.xiaoniuds.com/zh-CN/docs
```

They use a mock `http.Client` by default, so running them prints the SDK request and a mock response without calling GlobalRouter.

Run from `gr_sdk_demo/go_sdk`:

```bash
go run ./examples/create-chat-completion
go run ./examples/create-images
go run ./examples/create-image-task
go run ./examples/create-video
go run ./examples/create-seed-audio
```

## Files

- `create-chat-completion/` -> `POST /api/v1/chat/completions`, via `client.Chat.Create(...)`.
- `create-images/` -> `POST /api/v1/images`, via `client.Images.Generate(...)`.
- `create-image-task/` -> `POST /api/v1/image-tasks`, via `client.Images.CreateTask(...)`.
- `create-video/` -> `POST /api/v1/videos`, via `client.Videos.Generate(...)`.
- `create-seed-audio/` -> `POST /doubao/api/v3/tts/create`, via `client.Audio.CreateSeedAudio(...)`.

The examples intentionally include only the public `/api/v1` request shapes shown in the GlobalRouter docs.
