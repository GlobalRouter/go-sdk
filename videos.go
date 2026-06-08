package globalrouter

import (
	"context"
	"net/http"
)

type VideosResource struct {
	client *Client
}

func (r *VideosResource) Generate(ctx context.Context, request GenerationRequest, opts ...RequestOption) (*TaskResponse, error) {
	var out TaskResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/videos/generations", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
