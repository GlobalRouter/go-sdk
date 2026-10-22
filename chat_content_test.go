package globalrouter

import (
	"encoding/json"
	"testing"
)

func TestChatContentPartsPreserveMediaFields(t *testing.T) {
	parts := []ContentPart{
		{Type: "video", Video: map[string]any{"url": "https://example.test/a.mp4"}},
		{Type: "video_url", VideoURL: map[string]any{"url": "https://example.test/b.mp4", "fps": 2}},
		{Type: "input_audio", InputAudio: map[string]any{"data": "UklGRg==", "format": "wav"}},
		{Type: "file", File: map[string]any{"file_id": "file-test"}},
	}
	body, err := json.Marshal(ChatRequest{
		Model:    "gemini-3.8-flash",
		Messages: []Message{{Role: RoleUser, Content: parts}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var payload struct {
		Messages []struct {
			Content []map[string]any `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatal(err)
	}
	fields := []string{"video", "video_url", "input_audio", "file"}
	for i, field := range fields {
		part := payload.Messages[0].Content[i]
		if _, ok := part[field]; !ok || len(part) != 2 {
			t.Fatalf("content part %d lost its original format: %#v", i, part)
		}
	}
}
