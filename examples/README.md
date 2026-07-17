# GlobalRouter Go SDK examples

These examples adapt request payloads from `../../llm_gateway/docs` into calls through the Go SDK.
They use a mock `http.Client` by default, so running them prints the SDK request and a mock response without calling GlobalRouter.

Run from `gr_sdk_demo/go_sdk`:

```bash
go run ./examples/images-generations
EXAMPLE_NAME=text2image go run ./examples/seedream-image
GLOBALROUTER_BASE_URL=http://127.0.0.1:8000 go run ./examples/std-video-tasks
```

The generated docs examples use `client.RequestJSON(...)` for API surfaces that do not yet have typed SDK resource helpers. `examples/images-generations` uses the typed `client.Images.Generate(...)` helper.

## Files

- `images-generations/` -> `images_generations.yaml`, via `client.Images.Generate(...)`.
- `seedance-tasks/` -> `jimeng.yaml` `/v1/tasks`, via `client.RequestJSON(...)`.
- `jimeng-legacy/` -> `jimeng.yaml` `/v1/jimeng/create` and `/v1/jimeng/query`, via `client.RequestJSON(...)`.
- `seedream-image/` -> `seedream.yaml`, via `client.RequestJSON(...)`.
- `wan-image/` -> `wan.yaml`, via `client.RequestJSON(...)`.
- `happyhorse-video/` -> `happyhorse.yaml`, via `client.RequestJSON(...)`.
- `admin-price-rules/` -> `price-rules-admin.yaml`, mock-only admin example.

Not converted here: docs without concrete JSON request examples, the multipart upload doc, invalid `models-admin.yaml`, `ali_avatar.yaml` because GlobalRouter itself does not expose avatar as an SDK capability, `claude.yaml`, `gemini.yaml`, `std_tasks.yaml`, and `sync_password.yaml` because it contains password-sync token-like sample data rather than an SDK-facing API example.
