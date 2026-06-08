package globalrouter

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestChatCreateSendsBearerJSONAndDecodesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer gr_test" {
			t.Fatalf("authorization header = %q", got)
		}
		if got := r.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
			t.Fatalf("content-type header = %q", got)
		}
		if got := r.Header.Get("User-Agent"); !strings.Contains(got, "globalrouter-go") {
			t.Fatalf("user-agent header = %q", got)
		}
		var body ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Model != "openai/gpt-4o-mini" {
			t.Fatalf("model = %q", body.Model)
		}
		if len(body.Messages) != 1 || body.Messages[0].Role != RoleUser {
			t.Fatalf("messages = %#v", body.Messages)
		}
		writeJSON(t, w, ChatResponse{
			ID:              "chatcmpl_123",
			Object:          "chat.completion",
			Created:         1780880000,
			Model:           body.Model,
			RouterProvider:  "openai",
			RouterRequestID: "req_123",
			Choices: []ChatChoice{{
				Index:   0,
				Message: &Message{Role: RoleAssistant, Content: "hello"},
			}},
			Usage: &Usage{PromptTokens: 3, CompletionTokens: 4, TotalTokens: 7},
		})
	}))
	defer server.Close()

	client := New(
		WithAPIKey("gr_test"),
		WithBaseURL(server.URL),
		WithRetryConfig(RetryConfig{MaxRetries: 0}),
	)

	res, err := client.Chat.Create(context.Background(), ChatRequest{
		Model: "openai/gpt-4o-mini",
		Messages: []Message{{
			Role:    RoleUser,
			Content: "ping",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ID != "chatcmpl_123" || res.Choices[0].Message.Content != "hello" {
		t.Fatalf("unexpected response: %#v", res)
	}
}

func TestModelsListBuildsFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("modality") != "text" || q.Get("capability") != "chat" || q.Get("provider") != "openai" || q.Get("available_only") != "true" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		writeJSON(t, w, ModelsResponse{
			Object: "list",
			Data: []Model{{
				ID:                  "openai/gpt-4o-mini",
				Object:              "model",
				OwnedBy:             "openai",
				DisplayName:         "GPT-4o mini",
				Category:            "chat",
				Modality:            "text",
				Capabilities:        []string{"chat"},
				InputModalities:     []string{"text"},
				OutputModalities:    []string{"text"},
				SupportedParameters: []string{"temperature"},
				ExecutionMode:       "sync",
				BillingUnit:         "token",
				Routable:            true,
			}},
		})
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL))
	res, err := client.Models.List(context.Background(), &ListModelsOptions{
		Modality:      "text",
		Capability:    "chat",
		Provider:      "openai",
		AvailableOnly: Bool(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Data) != 1 || res.Data[0].ID != "openai/gpt-4o-mini" {
		t.Fatalf("unexpected response: %#v", res)
	}
}

func TestAPIErrorParsesEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(t, w, map[string]any{
			"error": map[string]any{
				"code":       "AUTHENTICATION_FAILED",
				"message":    "Invalid API key.",
				"type":       "authentication_error",
				"request_id": "req_bad",
			},
		})
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL), WithRetryConfig(RetryConfig{MaxRetries: 0}))
	_, err := client.Chat.Create(context.Background(), ChatRequest{Model: "m"})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized || apiErr.Code != "AUTHENTICATION_FAILED" || apiErr.RequestID != "req_bad" {
		t.Fatalf("api error = %#v", apiErr)
	}
}

func TestRetryRetriesServerErrorsOnly(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(t, w, map[string]any{"error": map[string]any{"code": "UPSTREAM", "message": "try again"}})
			return
		}
		writeJSON(t, w, ChatResponse{ID: "ok", Model: "m", Choices: []ChatChoice{{Index: 0}}})
	}))
	defer server.Close()

	client := New(
		WithBaseURL(server.URL),
		WithRetryConfig(RetryConfig{MaxRetries: 1, MinDelay: time.Millisecond}),
	)
	if _, err := client.Chat.Create(context.Background(), ChatRequest{Model: "m"}); err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestChatStreamParsesServerSentEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		var body ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if !body.Stream {
			t.Fatal("stream request did not force stream=true")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"id\":\"chunk_1\",\"model\":\"m\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"he\"}}],\"router_provider\":\"openai\",\"router_request_id\":\"req_1\"}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL))
	stream, err := client.Chat.Stream(context.Background(), ChatRequest{Model: "m"})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	event, err := stream.Next()
	if err != nil {
		t.Fatal(err)
	}
	if event.Data.ID != "chunk_1" || event.Data.Choices[0].Delta.Content != "he" {
		t.Fatalf("event = %#v", event)
	}
	if _, err := stream.Next(); !errors.Is(err, io.EOF) {
		t.Fatalf("second Next error = %v, want EOF", err)
	}
}

func TestTaskEventsPreserveEventNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/task_123/events" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "event: task.running\ndata: {\"event_id\":\"evt_1\",\"task_id\":\"task_123\",\"attempt\":1,\"event_type\":\"task.running\",\"status\":\"running\",\"progress\":0.5,\"payload\":{\"stage\":\"provider\"},\"created_at\":\"2026-06-08T00:00:00Z\"}\n\n")
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL))
	stream, err := client.Tasks.Events(context.Background(), "task_123")
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	event, err := stream.Next()
	if err != nil {
		t.Fatal(err)
	}
	if event.Event != "task.running" || event.Data.TaskID != "task_123" || event.Data.Progress != 0.5 {
		t.Fatalf("event = %#v", event)
	}
}

func TestTaskAndMultimodalResourcePaths(t *testing.T) {
	seen := map[string]bool{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen[r.Method+" "+r.URL.Path] = true
		if r.URL.Path == "/v1/videos/generations" && r.Header.Get("Idempotency-Key") != "idem_1" {
			t.Fatalf("missing idempotency header")
		}
		if r.Method == http.MethodGet && r.URL.Path == "/v1/tasks" {
			if r.URL.Query().Get("metadata.project_id") != "proj_1" {
				t.Fatalf("query = %s", r.URL.RawQuery)
			}
			writeJSON(t, w, TaskListResponse{Items: []TaskResponse{}, Total: 0, Page: 1, PageSize: 20})
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/v1/tasks/batch/batch_1" {
			writeJSON(t, w, TaskBatchResponse{BatchID: "batch_1", Items: []TaskResponse{}, Total: 0, Page: 1, PageSize: 0})
			return
		}
		if strings.Contains(r.URL.Path, "/tasks/") || r.URL.Path == "/v1/tasks" || r.URL.Path == "/v1/videos/generations" || r.URL.Path == "/v1/3d/generations" {
			writeJSON(t, w, TaskResponse{ID: "task_1", Object: "task", Status: TaskStatusQueued, Progress: 0.25, Type: TaskTypeVideoGeneration, Model: "m"})
			return
		}
		writeJSON(t, w, map[string]any{"ok": true})
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL), WithRetryConfig(RetryConfig{MaxRetries: 0}))
	ctx := context.Background()
	_, _ = client.Tasks.Create(ctx, TaskCreateRequest{Type: TaskTypeVideoGeneration, Model: "m", Input: map[string]any{"prompt": "x"}})
	_, _ = client.Tasks.List(ctx, &ListTasksOptions{MetadataProjectID: "proj_1"})
	_, _ = client.Tasks.Get(ctx, "task_1", &GetTaskOptions{Wait: true})
	_, _ = client.Tasks.CreateBatch(ctx, []TaskCreateRequest{{Type: TaskTypeVideoGeneration, Model: "m", Input: map[string]any{"prompt": "x"}}})
	_, _ = client.Tasks.GetBatch(ctx, "batch_1")
	_, _ = client.Tasks.Cancel(ctx, "task_1")
	_, _ = client.Tasks.Retry(ctx, "task_1")
	_, _ = client.Images.Generate(ctx, ImageGenerationRequest{Model: "m", Prompt: "image"})
	_, _ = client.Audio.CreateSpeech(ctx, AudioSpeechRequest{Model: "m", Input: "hello", Voice: "alloy"})
	_, _ = client.Audio.CreateTranscription(ctx, AudioTranscriptionRequest{Model: "m", FileURL: "https://example.com/a.wav"})
	_, _ = client.Videos.Generate(ctx, GenerationRequest{Model: "m", Prompt: "video"}, WithIdempotencyKey("idem_1"))
	_, _ = client.ThreeD.Generate(ctx, GenerationRequest{Model: "m", Prompt: "mesh"})
	_, _ = client.Embeddings.Create(ctx, EmbeddingsRequest{Model: "m", Input: "hello"})

	for _, key := range []string{
		"POST /v1/tasks",
		"GET /v1/tasks",
		"GET /v1/tasks/task_1",
		"POST /v1/tasks/batch",
		"GET /v1/tasks/batch/batch_1",
		"POST /v1/tasks/task_1/cancel",
		"POST /v1/tasks/task_1/retry",
		"POST /v1/images/generations",
		"POST /v1/audio/speech",
		"POST /v1/audio/transcriptions",
		"POST /v1/videos/generations",
		"POST /v1/3d/generations",
		"POST /v1/embeddings",
	} {
		if !seen[key] {
			t.Fatalf("did not see %s; seen=%v", key, seen)
		}
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatal(err)
	}
}
