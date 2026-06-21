package globalrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type RetryConfig struct {
	MaxRetries int
	MinDelay   time.Duration
	MaxDelay   time.Duration
}

func (r RetryConfig) withDefaults() RetryConfig {
	if r.MinDelay <= 0 {
		r.MinDelay = 250 * time.Millisecond
	}
	if r.MaxDelay <= 0 {
		r.MaxDelay = time.Second
	}
	if r.MaxDelay < r.MinDelay {
		r.MaxDelay = r.MinDelay
	}
	return r
}

type requestConfig struct {
	headers map[string]string
	timeout *time.Duration
	retry   *RetryConfig
}

type RequestOption func(*requestConfig)

type requestTimeoutMode int

const (
	requestTimeoutDisabled requestTimeoutMode = iota
	requestTimeoutUntilHeaders
	requestTimeoutUntilBodyClosed
)

func WithIdempotencyKey(key string) RequestOption {
	return WithHeader("Idempotency-Key", key)
}

func WithHeader(name, value string) RequestOption {
	return func(config *requestConfig) {
		if config.headers == nil {
			config.headers = map[string]string{}
		}
		config.headers[name] = value
	}
}

func WithRequestTimeout(timeout time.Duration) RequestOption {
	return func(config *requestConfig) {
		config.timeout = &timeout
	}
}

func WithRequestRetries(retry RetryConfig) RequestOption {
	return func(config *requestConfig) {
		normalized := retry.withDefaults()
		config.retry = &normalized
	}
}

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	params url.Values,
	body any,
	out any,
	opts ...RequestOption,
) error {
	res, err := c.do(ctx, method, path, params, body, "application/json", requestTimeoutUntilBodyClosed, opts...)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if out == nil || res.StatusCode == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, res.Body)
		return nil
	}
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return fmt.Errorf("globalrouter: decode response: %w", err)
	}
	return nil
}

func (c *Client) doStream(
	ctx context.Context,
	method string,
	path string,
	params url.Values,
	body any,
	opts ...RequestOption,
) (*http.Response, error) {
	return c.do(ctx, method, path, params, body, "text/event-stream", requestTimeoutUntilHeaders, opts...)
}

func (c *Client) do(
	ctx context.Context,
	method string,
	path string,
	params url.Values,
	body any,
	accept string,
	timeoutMode requestTimeoutMode,
	opts ...RequestOption,
) (*http.Response, error) {
	config := requestConfig{}
	for _, opt := range opts {
		opt(&config)
	}

	retryConfig := c.retry.withDefaults()
	if config.retry != nil {
		retryConfig = config.retry.withDefaults()
	}

	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("globalrouter: encode request: %w", err)
		}
	}

	var timeoutState *requestTimeoutState
	if timeoutMode != requestTimeoutDisabled {
		timeout := c.timeout
		if config.timeout != nil {
			timeout = *config.timeout
		}
		if timeout > 0 {
			ctx, timeoutState = newRequestTimeoutState(ctx, timeout, timeoutMode)
		}
	}
	cancelRequest := func() {
		if timeoutState != nil {
			timeoutState.cancel()
		}
	}

	var lastErr error
	for attempt := 0; attempt <= retryConfig.MaxRetries; attempt++ {
		req, err := c.newRequest(ctx, method, path, params, bodyBytes, accept, config.headers)
		if err != nil {
			cancelRequest()
			return nil, err
		}
		res, err := c.client.Do(req)
		if err != nil {
			if timeoutState != nil && timeoutState.deadlineExceeded() {
				err = context.DeadlineExceeded
			}
			lastErr = err
			if attempt < retryConfig.MaxRetries {
				if sleepErr := sleepRetry(ctx, retryConfig, attempt); sleepErr != nil {
					cancelRequest()
					if timeoutState != nil && timeoutState.deadlineExceeded() {
						sleepErr = context.DeadlineExceeded
					}
					return nil, sleepErr
				}
				continue
			}
			cancelRequest()
			return nil, fmt.Errorf("globalrouter: send request: %w", err)
		}
		if timeoutState != nil &&
			timeoutState.mode == requestTimeoutUntilHeaders &&
			(res.StatusCode < 500 || attempt >= retryConfig.MaxRetries) &&
			timeoutState.stopHeaderTimer() {
			_ = res.Body.Close()
			cancelRequest()
			return nil, fmt.Errorf("globalrouter: send request: %w", context.DeadlineExceeded)
		}
		if res.StatusCode >= 500 && attempt < retryConfig.MaxRetries {
			_, _ = io.Copy(io.Discard, res.Body)
			_ = res.Body.Close()
			if sleepErr := sleepRetry(ctx, retryConfig, attempt); sleepErr != nil {
				cancelRequest()
				if timeoutState != nil && timeoutState.deadlineExceeded() {
					sleepErr = context.DeadlineExceeded
				}
				return nil, sleepErr
			}
			continue
		}
		if res.StatusCode >= 400 {
			apiErr := parseAPIError(res)
			_ = res.Body.Close()
			cancelRequest()
			return nil, apiErr
		}
		if timeoutState != nil {
			if res.Body == nil {
				cancelRequest()
			} else {
				res.Body = cancelOnCloseReadCloser{ReadCloser: res.Body, cancel: timeoutState.cancel}
			}
		}
		return res, nil
	}
	cancelRequest()
	if lastErr != nil {
		return nil, fmt.Errorf("globalrouter: send request: %w", lastErr)
	}
	return nil, fmt.Errorf("globalrouter: request exhausted retries")
}

type requestTimeoutState struct {
	mode     requestTimeoutMode
	cancel   context.CancelFunc
	timer    *time.Timer
	timedOut atomic.Bool
}

func newRequestTimeoutState(ctx context.Context, timeout time.Duration, mode requestTimeoutMode) (context.Context, *requestTimeoutState) {
	if mode == requestTimeoutUntilBodyClosed {
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		return timeoutCtx, &requestTimeoutState{mode: mode, cancel: cancel}
	}
	timeoutCtx, cancel := context.WithCancel(ctx)
	state := &requestTimeoutState{mode: mode, cancel: cancel}
	state.timer = time.AfterFunc(timeout, func() {
		state.timedOut.Store(true)
		cancel()
	})
	return timeoutCtx, state
}

func (s *requestTimeoutState) deadlineExceeded() bool {
	return s.mode == requestTimeoutUntilHeaders && s.timedOut.Load()
}

func (s *requestTimeoutState) stopHeaderTimer() bool {
	if s.mode != requestTimeoutUntilHeaders || s.timer == nil {
		return false
	}
	if s.timer.Stop() {
		return false
	}
	return true
}

type cancelOnCloseReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (r cancelOnCloseReadCloser) Close() error {
	err := r.ReadCloser.Close()
	r.cancel()
	return err
}

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	path string,
	params url.Values,
	bodyBytes []byte,
	accept string,
	headers map[string]string,
) (*http.Request, error) {
	endpoint := c.baseURL + path
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	var body io.Reader
	if bodyBytes != nil {
		body = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("globalrouter: create request: %w", err)
	}
	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyBytes)), nil
		}
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func sleepRetry(ctx context.Context, config RetryConfig, attempt int) error {
	delay := config.MinDelay
	for i := 0; i < attempt; i++ {
		delay *= 2
		if delay >= config.MaxDelay {
			delay = config.MaxDelay
			break
		}
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func addString(values url.Values, key string, value string) {
	if value != "" {
		values.Set(key, value)
	}
}

func addInt(values url.Values, key string, value *int) {
	if value != nil {
		values.Set(key, strconv.Itoa(*value))
	}
}

func addBool(values url.Values, key string, value *bool) {
	if value != nil {
		values.Set(key, strconv.FormatBool(*value))
	}
}

func cleanPathValue(value string) string {
	return url.PathEscape(strings.TrimSpace(value))
}
