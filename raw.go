package globalrouter

import (
	"context"
)

// RequestJSON sends a JSON request to a custom GlobalRouter path and decodes a JSON response.
//
// It is useful for API surfaces that are available through GlobalRouter but do not yet have a
// typed resource helper in this SDK.
func (c *Client) RequestJSON(ctx context.Context, method string, path string, body any, out any, opts ...RequestOption) error {
	return c.doJSON(ctx, method, path, nil, body, out, opts...)
}
