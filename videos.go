package globalrouter

import (
	"context"
	"net/http"
)

type VideosResource struct {
	client *Client
}

func (r *VideosResource) Generate(ctx context.Context, request GenerationRequest, opts ...RequestOption) (*VideoGenerationResponse, error) {
	var out VideoGenerationResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/api/v1/videos", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
