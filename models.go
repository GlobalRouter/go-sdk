package globalrouter

import (
	"context"
	"net/http"
	"net/url"
)

type ModelsResource struct {
	client *Client
}

type ListModelsOptions struct {
	Modality      string
	Capability    string
	Provider      string
	AvailableOnly *bool
}

func (r *ModelsResource) List(ctx context.Context, options *ListModelsOptions, opts ...RequestOption) (*ModelsResponse, error) {
	params := url.Values{}
	if options != nil {
		addString(params, "modality", options.Modality)
		addString(params, "capability", options.Capability)
		addString(params, "provider", options.Provider)
		addBool(params, "available_only", options.AvailableOnly)
	}
	var out ModelsResponse
	if err := r.client.doJSON(ctx, http.MethodGet, "/v1/models", params, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
