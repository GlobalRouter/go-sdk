package globalrouter

import (
	"context"
	"net/http"
)

type AudioResource struct {
	client *Client
}

func (r *AudioResource) CreateSpeech(ctx context.Context, request AudioSpeechRequest, opts ...RequestOption) (map[string]any, error) {
	var out map[string]any
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/audio/speech", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AudioResource) CreateTranscription(ctx context.Context, request AudioTranscriptionRequest, opts ...RequestOption) (map[string]any, error) {
	var out map[string]any
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/audio/transcriptions", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}
