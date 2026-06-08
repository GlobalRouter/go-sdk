package globalrouter

import (
	"context"
	"net/http"
)

type ThreeDResource struct {
	client *Client
}

func (r *ThreeDResource) Generate(ctx context.Context, request GenerationRequest, opts ...RequestOption) (*TaskResponse, error) {
	var out TaskResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/3d/generations", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
