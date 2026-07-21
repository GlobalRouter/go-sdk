package globalrouter

import (
	"context"
	"net/http"
)

// SeedanceCompatibilityResource provides typed access to the additional
// Seedance video-generation and asset compatibility endpoints. The existing
// Videos resource remains the primary GlobalRouter video API.
type SeedanceCompatibilityResource struct {
	client *Client
}

func (r *SeedanceCompatibilityResource) CreateVideoGeneration(
	ctx context.Context,
	request SeedanceVideoGenerationRequest,
	opts ...RequestOption,
) (*SeedanceVideoGenerationResponse, error) {
	var out SeedanceVideoGenerationResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/v1/video/generations", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *SeedanceCompatibilityResource) GetVideoGeneration(
	ctx context.Context,
	taskID string,
	opts ...RequestOption,
) (*SeedanceVideoGenerationResponse, error) {
	var out SeedanceVideoGenerationResponse
	if err := r.client.doJSON(ctx, http.MethodGet, "/v1/video/generations/"+cleanPathValue(taskID), nil, nil, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *SeedanceCompatibilityResource) CreateAssetGroup(
	ctx context.Context,
	request SeedanceAssetGroupCreateRequest,
	opts ...RequestOption,
) (*SeedanceAssetGroupResponse, error) {
	var out SeedanceAssetGroupResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/api/v3/assets/groups", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *SeedanceCompatibilityResource) CreateAsset(
	ctx context.Context,
	request SeedanceAssetCreateRequest,
	opts ...RequestOption,
) (*SeedanceAssetResponse, error) {
	var out SeedanceAssetResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/api/v3/assets", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *SeedanceCompatibilityResource) GetAsset(
	ctx context.Context,
	request SeedanceAssetGetRequest,
	opts ...RequestOption,
) (*SeedanceAssetResponse, error) {
	var out SeedanceAssetResponse
	if err := r.client.doJSON(ctx, http.MethodPost, "/api/v3/assets/get", nil, request, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
