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
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/images/generations", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}
