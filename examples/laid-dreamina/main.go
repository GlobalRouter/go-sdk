// 默认离线 mock；GLOBALROUTER_EXAMPLE_REAL=1 才连接海外站点并产生实际费用。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	globalrouter "github.com/GlobalRouter/go-sdk"
)

func main() {
	real := os.Getenv("GLOBALROUTER_EXAMPLE_REAL") == "1"
	apiKey := "sk-local-example"
	if real {
		apiKey = os.Getenv("GLOBALROUTER_API_KEY")
	}
	baseURL := os.Getenv("GLOBALROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.globalrouter.ai"
	}
	capture := &mockClient{}
	opts := []globalrouter.SDKOption{
		globalrouter.WithAPIKey(apiKey), globalrouter.WithBaseURL(baseURL),
		globalrouter.WithRetryConfig(globalrouter.RetryConfig{MaxRetries: 0}),
	}
	if !real {
		opts = append(opts, globalrouter.WithClient(capture))
	}
	client := globalrouter.New(opts...)
	request := globalrouter.GenerationRequest{
		Model:      "dreamina-seedance-2-0-260128",                      // 公共 ID，不是上游 doubao ID
		Provider:   &globalrouter.ProviderSelection{ProviderID: "laid"}, // Dreamina 必填
		Prompt:     "A quiet city rooftop at sunset, slow cinematic push-in",
		Resolution: "720p", Duration: globalrouter.Int(5),
	}
	if uri := os.Getenv("GLOBALROUTER_ASSET_URI"); uri != "" {
		// 必须是本租户、同 Laid 账号线路、已 Active 的平台素材 URI。
		request.InputReferences = []map[string]any{{"type": "image_url", "image_url": map[string]any{"url": uri}}}
	}
	created, err := client.Videos.Generate(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	printJSON(created)
	// Go SDK 的旧生成入口对应 Tasks.Get；使用平台 ID，不传 provider 或上游 task_id。
	queried, err := client.Tasks.Get(context.Background(), created.ID, nil)
	if err != nil {
		log.Fatal(err)
	}
	printJSON(queried)
	if !real {
		if capture.count != 2 {
			log.Fatal("expected one create and one query")
		}
		fmt.Println("Mock 创建/查询已通过；未调用上游。")
	}
}

type mockClient struct{ count int }

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	m.count++
	status, body := http.StatusOK, `{"id":"job_mock_dreamina","status":"succeeded","model":"dreamina-seedance-2-0-260128","type":"video_generation"}`
	if req.Method == http.MethodPost {
		if req.URL.Path != "/v1/videos/generations" {
			return nil, fmt.Errorf("unexpected create path: %s", req.URL.Path)
		}
		var payload globalrouter.GenerationRequest
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, err
		}
		if payload.Provider == nil || payload.Provider.ProviderID != "laid" || payload.Resolution != "720p" {
			return nil, fmt.Errorf("missing explicit provider/resolution")
		}
		status = http.StatusAccepted
		body = `{"id":"job_mock_dreamina","status":"pending","polling_url":"/v1/tasks/job_mock_dreamina"}`
	} else if req.Method != http.MethodGet || req.URL.Path != "/v1/tasks/job_mock_dreamina" || req.URL.RawQuery != "" {
		return nil, fmt.Errorf("unexpected query: %s %s", req.Method, req.URL)
	}
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
}
func printJSON(value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
