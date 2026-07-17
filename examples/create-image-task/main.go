package main

import (
	"context"
	"log"

	globalrouter "github.com/GlobalRouter/go-sdk"
	"github.com/GlobalRouter/go-sdk/examples/internal/exampleutil"
)

func main() {
	ctx := context.Background()
	capture := &exampleutil.CaptureClient{}
	client := exampleutil.NewClient(capture)

	matched := false
	if !exampleutil.RunSelected("jimeng_image_task", &matched) {
		exampleutil.EnsureSelection(matched)
		return
	}

	before := len(capture.Requests)
	response, err := client.Images.CreateTask(ctx, globalrouter.ImageTaskCreateRequest{
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
		Metadata: map[string]any{
			"client_request_id": "client-001",
		},
		IdempotencyKey: "client-image-task-001",
	})
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.PrintCaptured("Create image task (/api/v1/image-tasks)", "jimeng_image_task", capture, before, response)
	exampleutil.EnsureSelection(matched)
}
