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

func (r *ImagesResource) CreateTask(ctx context.Context, request ImageTaskCreateRequest, opts ...RequestOption) (*ImageTaskResponse, error) {
	if request.IdempotencyKey != "" {
		opts = append([]RequestOption{WithIdempotencyKey(request.IdempotencyKey)}, opts...)
	}
	var out ImageTaskResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/api/v1/image-tasks", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ImagesResource) GetTask(ctx context.Context, imageTaskID string, opts ...RequestOption) (*ImageTaskResponse, error) {
	var out ImageTaskResponse
	if err := r.client.doJSON(ctx, http.MethodGet, "/api/v1/image-tasks/"+cleanPathValue(imageTaskID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
