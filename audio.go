package globalrouter

import (
	"context"
	"net/http"
	"strings"
)

type AudioResource struct {
	client *Client
}

func (r *AudioResource) CreateSpeech(ctx context.Context, request AudioSpeechRequest, opts ...RequestOption) (*http.Response, error) {
	return r.client.doBinary(ctx, http.MethodPost, "/v1/audio/speech", nil, request, audioSpeechAccept(request.ResponseFormat), opts...)
}

func (r *AudioResource) CreateTranscription(ctx context.Context, request AudioTranscriptionRequest, opts ...RequestOption) (map[string]any, error) {
	var out map[string]any
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/audio/transcriptions", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func audioSpeechAccept(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "mp3":
		return "audio/mpeg"
	case "opus":
		return "audio/opus"
	case "aac":
		return "audio/aac"
	case "flac":
		return "audio/flac"
	case "wav":
		return "audio/wav"
	case "pcm":
		return "audio/pcm"
	default:
		return "*/*"
	}
}
