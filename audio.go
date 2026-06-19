package globalrouter

import (
	"context"
	"net/http"
)

type AudioResource struct {
	client *Client
}

// CreateSpeech returns the raw response body so callers can read binary audio.
// The caller is responsible for closing the returned response body.
func (r *AudioResource) CreateSpeech(ctx context.Context, request AudioSpeechRequest, opts ...RequestOption) (*http.Response, error) {
	return r.client.do(ctx, http.MethodPost, "/v1/audio/speech", nil, request, "*/*", requestTimeoutUntilBodyClosed, opts...)
}

func (r *AudioResource) CreateTranscription(ctx context.Context, request AudioTranscriptionRequest, opts ...RequestOption) (map[string]any, error) {
	var out map[string]any
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/audio/transcriptions", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}
