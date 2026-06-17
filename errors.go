package globalrouter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Type       string
	RequestID  string
	Body       string
	Response   *http.Response
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Code != "" {
		return fmt.Sprintf("globalrouter: %s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("globalrouter: http %d: %s", e.StatusCode, e.Message)
}

func parseAPIError(res *http.Response) *APIError {
	raw, _ := io.ReadAll(res.Body)
	errOut := &APIError{
		StatusCode: res.StatusCode,
		Code:       "ROUTER_SDK_HTTP_ERROR",
		Type:       "router_error",
		Body:       string(raw),
		Response:   res,
	}
	var envelope struct {
		Error struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			Type      string `json:"type"`
			RequestID string `json:"request_id"`
		} `json:"error"`
	}
	if json.Unmarshal(raw, &envelope) == nil {
		if envelope.Error.Code != "" {
			errOut.Code = envelope.Error.Code
		}
		if envelope.Error.Message != "" {
			errOut.Message = envelope.Error.Message
		}
		if envelope.Error.Type != "" {
			errOut.Type = envelope.Error.Type
		}
		errOut.RequestID = envelope.Error.RequestID
	}
	if errOut.Message == "" {
		errOut.Message = string(raw)
	}
	if errOut.Message == "" {
		errOut.Message = http.StatusText(res.StatusCode)
	}
	return errOut
}
