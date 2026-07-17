package globalrouter

import (
	"context"
	"net/http"
	"strings"
)

type ImagesResource struct {
	client *Client
}

func (r *ImagesResource) Generate(ctx context.Context, request ImageGenerationRequest, opts ...RequestOption) (map[string]any, error) {
	var out map[string]any
	if err := r.client.doJSON(ctx, http.MethodPost, "/api/v1/images", nil, sanitizeImageGenerationRequest(request), &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func sanitizeImageGenerationRequest(request ImageGenerationRequest) ImageGenerationRequest {
	if request.Provider == nil {
		return request
	}
	provider := *request.Provider
	if strings.TrimSpace(provider.ProviderID) == "" {
		provider.ProviderID = ""
	}
	if provider.ProviderID == "" && len(provider.Options) == 0 {
		request.Provider = nil
		return request
	}
	request.Provider = &provider
	return request
}

func (r *ImagesResource) CreateTask(ctx context.Context, request ImageTaskCreateRequest, opts ...RequestOption) (*ImageTaskResponse, error) {
	if request.IdempotencyKey != "" {
		opts = append([]RequestOption{WithIdempotencyKey(request.IdempotencyKey)}, opts...)
		request.IdempotencyKey = ""
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
