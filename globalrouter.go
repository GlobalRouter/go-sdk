package globalrouter

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	SDKVersion     = "0.1.0"
	defaultBaseURL = "https://api.globalrouter.ai"
	defaultUA      = "globalrouter-go/" + SDKVersion
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseURL   string
	apiKey    string
	client    HTTPClient
	timeout   time.Duration
	retry     RetryConfig
	retrySet  bool
	userAgent string

	Models     *ModelsResource
	Chat       *ChatResource
	Embeddings *EmbeddingsResource
	Tasks      *TasksResource
	Images     *ImagesResource
	Audio      *AudioResource
	Videos     *VideosResource
	ThreeD     *ThreeDResource
}

type SDKOption func(*Client)

func New(opts ...SDKOption) *Client {
	c := &Client{
		baseURL:   defaultBaseURL,
		apiKey:    os.Getenv("GLOBALROUTER_API_KEY"),
		timeout:   60 * time.Second,
		retry:     RetryConfig{MaxRetries: 2, MinDelay: 250 * time.Millisecond, MaxDelay: time.Second},
		userAgent: defaultUA,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.baseURL = strings.TrimRight(c.baseURL, "/")
	if c.client == nil {
		c.client = &http.Client{}
	}
	c.Models = &ModelsResource{client: c}
	c.Chat = &ChatResource{client: c}
	c.Embeddings = &EmbeddingsResource{client: c}
	c.Tasks = &TasksResource{client: c}
	c.Images = &ImagesResource{client: c}
	c.Audio = &AudioResource{client: c}
	c.Videos = &VideosResource{client: c}
	c.ThreeD = &ThreeDResource{client: c}
	return c
}

func WithAPIKey(apiKey string) SDKOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

func WithBaseURL(baseURL string) SDKOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithClient(client HTTPClient) SDKOption {
	return func(c *Client) {
		c.client = client
	}
}

func WithTimeout(timeout time.Duration) SDKOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithRetryConfig(config RetryConfig) SDKOption {
	return func(c *Client) {
		c.retry = config.withDefaults()
		c.retrySet = true
	}
}

func WithUserAgent(userAgent string) SDKOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func String(v string) *string { return &v }

func Bool(v bool) *bool { return &v }

func Int(v int) *int { return &v }

func Int64(v int64) *int64 { return &v }

func Float64(v float64) *float64 { return &v }

func Pointer[T any](v T) *T { return &v }

type SecuritySource func(context.Context) (string, error)
