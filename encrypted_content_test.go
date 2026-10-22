package globalrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncryptedContentRoundtrip(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		t.Run(fmt.Sprint(streaming), func(t *testing.T) {
			model := "doubao-seed-2-1-pro-260915"
			opaque := "opaque+/state=="
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req ChatRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatal(err)
				}
				if len(req.Messages) > 1 && req.Messages[1].EncryptedContent != opaque {
					t.Fatal("encrypted state was lost in the next request")
				}
				if req.Stream {
					w.Header().Set("Content-Type", "text/event-stream")
					fmt.Fprintf(w, "data: {\"id\":\"c1\",\"model\":%q,\"choices\":[{\"index\":0,\"delta\":{\"encrypted_content\":%q,\"content\":\"\"}}]}\n\ndata: [DONE]\n\n", model, opaque)
					return
				}
				writeJSON(t, w, ChatResponse{
					ID: "c1", Model: model, Choices: []ChatChoice{{
						Message: &Message{Role: RoleAssistant, Content: "", EncryptedContent: opaque},
					}},
				})
			}))
			defer server.Close()
			client := New(WithAPIKey("test-key"), WithBaseURL(server.URL))
			req := ChatRequest{Model: model, Messages: []Message{{Role: RoleUser, Content: "hello"}}}
			var assistant Message
			if streaming {
				stream, err := client.Chat.Stream(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				defer stream.Close()
				chunk, err := stream.Next()
				if err != nil {
					t.Fatal(err)
				}
				assistant = Message{Role: RoleAssistant, EncryptedContent: chunk.Data.Choices[0].Delta.EncryptedContent}
			} else {
				res, err := client.Chat.Create(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				assistant = *res.Choices[0].Message
			}
			if assistant.EncryptedContent != opaque {
				t.Fatal("encrypted state was lost in response parsing")
			}
			req.Messages = append(req.Messages, assistant, Message{Role: RoleUser, Content: "continue"})
			if _, err := client.Chat.Create(context.Background(), req); err != nil {
				t.Fatal(err)
			}
		})
	}
}
