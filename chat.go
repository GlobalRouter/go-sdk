package globalrouter

import (
	"context"
	"net/http"
)

type ChatResource struct {
	client *Client
}

func (r *ChatResource) Create(ctx context.Context, request ChatRequest, opts ...RequestOption) (*ChatResponse, error) {
	var out ChatResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/chat/completions", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ChatResource) Stream(ctx context.Context, request ChatRequest, opts ...RequestOption) (*SSEStream[ChatResponse], error) {
	request.Stream = true
	res, err := r.client.doStream(ctx, http.MethodPost, "/v1/chat/completions", nil, request, opts...)
	if err != nil {
		return nil, err
	}
	return newSSEStream[ChatResponse](res), nil
}
