package globalrouter

import (
	"context"
	"net/http"
)

type EmbeddingsResource struct {
	client *Client
}

func (r *EmbeddingsResource) Create(ctx context.Context, request EmbeddingsRequest, opts ...RequestOption) (map[string]any, error) {
	var out map[string]any
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/embeddings", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}
