package main

import (
	"context"
	"log"
	"net/http"

	"github.com/GlobalRouter/go-sdk/examples/internal/exampleutil"
)

type example struct {
	Name string
	Body string
}

func main() {
	ctx := context.Background()
	capture := &exampleutil.CaptureClient{}
	client := exampleutil.NewClient(capture)

	matched := false
	for _, item := range examples {
		if !exampleutil.RunSelected(item.Name, &matched) {
			continue
		}
		payload := exampleutil.MustObject(item.Body)
		before := len(capture.Requests)
		var response map[string]any
		if err := client.RequestJSON(ctx, http.MethodPost, "/v1/wan/image/create", payload, &response); err != nil {
			log.Fatal(err)
		}
		exampleutil.PrintCaptured("Wan image generation and editing", item.Name, capture, before, response)
	}
	exampleutil.EnsureSelection(matched)
}

var examples = []example{
	{Name: "text2image", Body: `{
  "model": "wan2.7-image-pro",
  "input": {
    "messages": [
      {
        "role": "user",
        "content": [
          {
            "text": "一间有着精致窗户的花店，漂亮的木质门，摆放着花朵"
          }
        ]
      }
    ]
  },
  "parameters": {
    "size": "2K",
    "n": 1,
    "watermark": false,
    "thinking_mode": true
  }
}`},
	{Name: "image_edit", Body: `{
  "model": "wan2.7-image-pro",
  "input": {
    "messages": [
      {
        "role": "user",
        "content": [
          {
            "image": "https://example.com/car.webp"
          },
          {
            "image": "https://example.com/paint.webp"
          },
          {
            "text": "把图2的涂鸦喷绘在图1的汽车上"
          }
        ]
      }
    ]
  },
  "parameters": {
    "size": "2K",
    "n": 1,
    "watermark": false
  }
}`},
}
