package globalrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageGenerationPreservesFalseWatermark(t *testing.T) {
	for _, direct := range []bool{false, true} {
		t.Run(map[bool]string{false: "native", true: "direct"}[direct], func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				wantPath := "/v1/images/generations"
				if direct {
					wantPath = "/api/v1/images"
				}
				if r.URL.Path != wantPath {
					t.Fatalf("path = %s, want %s", r.URL.Path, wantPath)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body["watermark"] != false {
					t.Fatalf("watermark = %v, want false", body["watermark"])
				}
				if direct {
					provider, ok := body["provider"].(map[string]any)
					if !ok || provider["provider_id"] != "doubao" {
						t.Fatalf("provider = %v", body["provider"])
					}
				}
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte(`{"error":{"message":"unavailable"}}`))
			}))
			defer server.Close()

			client := New(WithBaseURL(server.URL), WithRetryConfig(RetryConfig{MaxRetries: 2}))
			watermark := false
			if direct {
				_, _ = client.Images.GenerateOpenRouter(context.Background(), map[string]any{
					"model":     "doubao-seedream-5-0-pro-260628",
					"prompt":    "一只像素风格的花猫坐在草坪上",
					"provider":  map[string]any{"provider_id": "doubao"},
					"watermark": false,
				})
			} else {
				_, _ = client.Images.Generate(context.Background(), ImageGenerationRequest{
					Model: "doubao-seedream-5-0-pro-260628", Prompt: "cat", Watermark: &watermark,
				})
			}
			if requests != 1 {
				t.Fatalf("requests = %d, want 1", requests)
			}
		})
	}
}
