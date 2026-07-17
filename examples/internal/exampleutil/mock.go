package exampleutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	globalrouter "github.com/GlobalRouter/go-sdk"
)

type CaptureClient struct {
	Requests []*http.Request
	Bodies   [][]byte
}

func (c *CaptureClient) Do(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	c.Requests = append(c.Requests, req)
	c.Bodies = append(c.Bodies, bodyBytes)

	var payload map[string]any
	_ = json.Unmarshal(bodyBytes, &payload)
	response := map[string]any{
		"id":     "example_mock",
		"object": "example.response",
		"status": "mocked",
		"model":  payload["model"],
		"data": map[string]any{
			"path": req.URL.Path,
		},
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

func NewClient(capture *CaptureClient) *globalrouter.Client {
	apiKey := os.Getenv("GLOBALROUTER_API_KEY")
	if apiKey == "" {
		apiKey = "sk-local-example"
	}
	baseURL := os.Getenv("GLOBALROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8000"
	}
	return globalrouter.New(
		globalrouter.WithAPIKey(apiKey),
		globalrouter.WithBaseURL(baseURL),
		globalrouter.WithClient(capture),
		globalrouter.WithRetryConfig(globalrouter.RetryConfig{MaxRetries: 0}),
	)
}

func RunSelected(name string, matched *bool) bool {
	selected := os.Getenv("EXAMPLE_NAME")
	if selected != "" && selected != name {
		return false
	}
	*matched = true
	return true
}

func EnsureSelection(matched bool) {
	selected := os.Getenv("EXAMPLE_NAME")
	if selected != "" && !matched {
		panic(fmt.Sprintf("unknown EXAMPLE_NAME=%q", selected))
	}
}

func MustObject(raw string) map[string]any {
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		panic(err)
	}
	return payload
}

func PrintCaptured(title string, name string, capture *CaptureClient, index int, response any) {
	req := capture.Requests[index]
	body := capture.Bodies[index]
	fmt.Printf("\n## %s: %s\n", title, name)
	fmt.Println("\n# Request JSON")
	PrintJSONFromBytes(body)
	fmt.Println("\n# cURL")
	fmt.Println(CurlForRequest(req, body))
	fmt.Println("\n# Mock response")
	PrintJSON(response)
}

func PrintJSON(value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func PrintJSONFromBytes(data []byte) {
	var pretty bytes.Buffer
	if len(data) == 0 {
		fmt.Println("null")
		return
	}
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		panic(err)
	}
	fmt.Println(pretty.String())
}

func CurlForRequest(req *http.Request, body []byte) string {
	var pretty bytes.Buffer
	if len(body) > 0 {
		if err := json.Indent(&pretty, body, "", "  "); err != nil {
			panic(err)
		}
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
	if pretty.Len() > 0 {
		lines = append(lines, fmt.Sprintf("  --data-raw %s", shellQuote(pretty.String())))
	}
	return strings.Join(lines, " \\\n")
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
