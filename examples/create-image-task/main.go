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
		response, err := client.Images.CreateTask(
			ctx,
			requestBody(),
			globalrouter.WithIdempotencyKey("client-image-task-001"),
		)
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

	response, err := client.Images.CreateTask(ctx, requestBody(), globalrouter.WithIdempotencyKey("client-image-task-001"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("# POST /api/v1/image-tasks")
	fmt.Println("\n# Request JSON")
	printJSONFromBytes(capture.bodies[0])
	fmt.Println("\n# cURL")
	fmt.Println(curlForRequest(capture.requests[0], capture.bodies[0]))
	fmt.Println("\n# Mock response")
	printJSON(response)
}

func requestBody() globalrouter.ImageTaskCreateRequest {
	return globalrouter.ImageTaskCreateRequest{
		Model:  "jimeng_t2i_v31",
		Prompt: "生成 4 张电商商品图，白底，高级感",
		InputReferences: []globalrouter.ImageTaskReference{{
			Type:     "image_url",
			ImageURL: map[string]any{"url": "https://example.com/reference.png"},
		}},
		N:    globalrouter.Int(4),
		Size: "1024x1024",
		Provider: &globalrouter.ProviderSelection{
			ProviderID: "doubao_gr",
			Options: map[string]map[string]any{
				"doubao_gr": {
					"watermark":    false,
					"force_single": false,
				},
			},
		},
		Metadata: map[string]any{"client_request_id": "client-001"},
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
		"id":          "imgtask_xxx",
		"object":      "image.task",
		"status":      "queued",
		"model":       "jimeng_t2i_v31",
		"provider":    "doubao_gr",
		"polling_url": "/api/v1/image-tasks/imgtask_xxx",
		"created_at":  1770000000,
		"updated_at":  1770000000,
		"data":        []any{},
		"usage":       map[string]any{"billing_status": "unreserved"},
		"metadata":    map[string]any{"client_request_id": "client-001"},
	}
	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
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
