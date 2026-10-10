package globalrouter

import (
	"context"
	"net/http"
)

type ImagesResource struct {
	client *Client
}

func (r *ImagesResource) Generate(ctx context.Context, request ImageGenerationRequest, opts ...RequestOption) (map[string]any, error) {
	var out map[string]any
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/images/generations", nil, request, &out, append(opts, WithRequestRetries(RetryConfig{MaxRetries: 0}))...); err != nil {
		return nil, err
	}
	return out, nil
}

// GenerateOpenRouter accepts the direct image contract, including object-form provider selection.
func (r *ImagesResource) GenerateOpenRouter(ctx context.Context, request map[string]any, opts ...RequestOption) (map[string]any, error) {
	var out map[string]any
	if err := r.client.doJSON(ctx, http.MethodPost, "/api/v1/images", nil, request, &out, append(opts, WithRequestRetries(RetryConfig{MaxRetries: 0}))...); err != nil {
		return nil, err
	}
	return out, nil
}
