package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	globalrouter "github.com/GlobalRouter/go-sdk"
)

func main() {
	ctx := context.Background()

	if os.Getenv("GLOBALROUTER_EXAMPLE_REAL") == "1" {
		opts := []globalrouter.SDKOption{globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY"))}
		if baseURL := os.Getenv("GLOBALROUTER_BASE_URL"); baseURL != "" {
			opts = append(opts, globalrouter.WithBaseURL(baseURL))
		}
		client := globalrouter.New(opts...)
		response, err := client.Videos.Generate(ctx, requestBody())
		if err != nil {
			log.Fatal(err)
		}
		printJSON(response)
		return
	}

	capture := &captureClient{}
	client := globalrouter.New(
		globalrouter.WithAPIKey("sk-local-example"),
		globalrouter.WithBaseURL("http://127.0.0.1:8000"),
		globalrouter.WithClient(capture),
		globalrouter.WithRetryConfig(globalrouter.RetryConfig{MaxRetries: 0}),
	)

	response, err := client.Videos.Generate(ctx, requestBody())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# POST /api/v1/videos")
	fmt.Println("\n# Request JSON")
	printJSONFromBytes(capture.bodies[0])
	fmt.Println("\n# cURL")
	fmt.Println(curlForRequest(capture.requests[0], capture.bodies[0]))
	fmt.Println("\n# Mock response")
	printJSON(response)
}

func requestBody() globalrouter.GenerationRequest {
	return globalrouter.GenerationRequest{
		Model:         "doubao-seedance-1-0-pro-fast-251015",
		Prompt:        "A quiet city rooftop at sunset, slow cinematic push-in",
		AspectRatio:   "16:9",
		Duration:      globalrouter.Int(5),
		Resolution:    "720p",
		SR:            &globalrouter.SuperResolutionRequest{Resolution: "1080p"},
		Seed:          globalrouter.Int(12345),
		GenerateAudio: globalrouter.Bool(true),
		CallbackURL:   "https://example.com/webhooks/video",
		FrameImages: []globalrouter.VideoFrameImage{{
			Type:      "image_url",
			ImageURL:  map[string]any{"url": "https://example.com/assets/opening-frame.png"},
			FrameType: "first_frame",
		}},
		InputReferences: []map[string]any{{
			"type":      "image_url",
			"image_url": map[string]any{"url": "https://example.com/assets/reference.png"},
		}},
		Provider: &globalrouter.ProviderSelection{
			ProviderID: "doubao",
			Options: map[string]map[string]any{
				"doubao": {"some_provider_option": "value"},
			},
		},
	}
}

type captureClient struct {
	requests []*http.Request
	bodies   [][]byte
}

func (c *captureClient) Do(req *http.Request) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	_ = req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	c.requests = append(c.requests, req)
	c.bodies = append(c.bodies, bodyBytes)

	response := map[string]any{
		"id":          "job_123",
		"polling_url": "/api/v1/videos/job_123",
		"status":      "pending",
	}
	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusAccepted,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(data)),
		Request:    req,
	}, nil
}

func printJSON(value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}

func printJSONFromBytes(data []byte) {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		log.Fatal(err)
	}
	fmt.Println(pretty.String())
}

func curlForRequest(req *http.Request, body []byte) string {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "", "  "); err != nil {
		log.Fatal(err)
	}
	lines := []string{fmt.Sprintf("curl -X %s %s", req.Method, shellQuote(req.URL.String()))}
	keys := make([]string, 0, len(req.Header))
	for key := range req.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if strings.EqualFold(key, "Content-Length") {
			continue
		}
		for _, value := range req.Header.Values(key) {
			if strings.EqualFold(key, "Authorization") {
				value = "Bearer ${GLOBALROUTER_API_KEY}"
			}
			lines = append(lines, fmt.Sprintf("  -H %s", shellQuote(key+": "+value)))
		}
	}
	lines = append(lines, fmt.Sprintf("  --data-raw %s", shellQuote(pretty.String())))
	return strings.Join(lines, " \\\n")
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
