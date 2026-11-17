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

	globalrouter "github.com/GlobalRouter/go-sdk"
)

const model = "gemini-nano-banana-2.1"

func main() {
	ctx := context.Background()
	if os.Getenv("GLOBALROUTER_EXAMPLE_REAL") == "1" {
		opts := []globalrouter.SDKOption{globalrouter.WithAPIKey(os.Getenv("GLOBALROUTER_API_KEY"))}
		if baseURL := os.Getenv("GLOBALROUTER_BASE_URL"); baseURL != "" {
			opts = append(opts, globalrouter.WithBaseURL(baseURL))
		}
		client := globalrouter.New(opts...)
		response, err := client.Images.GenerateContent(ctx, model, requestBody())
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
	response, err := client.Images.GenerateContent(ctx, model, requestBody())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("# POST /v1beta/models/%s:generateContent\n", model)
	fmt.Println("\n# Request JSON")
	printJSONFromBytes(capture.bodies[0])
	fmt.Println("\n# Mock response")
	printJSON(response)
}

func requestBody() globalrouter.GeminiGenerateContentRequest {
	return globalrouter.GeminiGenerateContentRequest{
		Contents: []globalrouter.GeminiImageContent{{
			Role: "user",
			Parts: []globalrouter.GeminiImagePart{{
				Text: "Generate a minimal 16:9 product illustration on a white background.",
			}},
		}},
		GenerationConfig: &globalrouter.GeminiImageGenerationConfig{
			ResponseModalities: []string{"TEXT", "IMAGE"},
			ImageConfig: &globalrouter.GeminiImageConfig{
				AspectRatio: "16:9",
				ImageSize:   "2K",
			},
		},
	}
}

type captureClient struct {
	bodies [][]byte
}

func (c *captureClient) Do(req *http.Request) (*http.Response, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	_ = req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	c.bodies = append(c.bodies, bodyBytes)

	response := globalrouter.GeminiGenerateContentResponse{
		Candidates: []globalrouter.GeminiImageCandidate{{
			Content: globalrouter.GeminiImageContent{
				Role: "model",
				Parts: []globalrouter.GeminiImagePart{
					{Text: "Here is the generated image."},
					{InlineData: &globalrouter.GeminiImageInlineData{
						MimeType: "image/png",
						Data:     "iVBORw0KGgoAAAANSUhEUg...",
					}},
				},
			},
			FinishReason: "STOP",
		}},
		ModelVersion: model,
		ResponseID:   "response_xxx",
		UsageMetadata: &globalrouter.GeminiImageUsageMetadata{
			PromptTokenCount:     globalrouter.Int(18),
			CandidatesTokenCount: globalrouter.Int(1290),
			TotalTokenCount:      globalrouter.Int(1308),
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
