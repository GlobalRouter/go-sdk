package main

import (
	"context"
	"log"

	"github.com/GlobalRouter/go-sdk/examples/internal/exampleutil"
)

type example struct {
	Name   string
	Method string
	Path   string
	Body   string
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
		if err := client.RequestJSON(ctx, item.Method, item.Path, payload, &response); err != nil {
			log.Fatal(err)
		}
		exampleutil.PrintCaptured("Jimeng legacy video tasks", item.Name, capture, before, response)
	}
	exampleutil.EnsureSelection(matched)
}

var examples = []example{
	{Name: "text2video_720p", Method: "POST", Path: "/v1/jimeng/create", Body: `{
  "req_key": "jimeng_ti2v_v30_pro",
  "prompt": "千军万马在草原上奔跑的壮观场景",
  "frames": 121,
  "aspect_ratio": "16:9"
}`},
	{Name: "text2video_1080p", Method: "POST", Path: "/v1/jimeng/create", Body: `{
  "req_key": "jimeng_ti2v_v30_1080",
  "prompt": "千军万马在草原上奔跑的壮观场景",
  "frames": 121,
  "aspect_ratio": "16:9"
}`},
	{Name: "single_frame2video_720p", Method: "POST", Path: "/v1/jimeng/create", Body: `{
  "req_key": "jimeng_i2v_first_v30",
  "prompt": "千军万马",
  "image_urls": [
    "https://img.alicdn.com/imgextra/i1/O1CN01gDEY8M1W114Hi3XcN_!!6000000002727-0-tps-1024-406.jpg"
  ],
  "frames": 121,
  "aspect_ratio": "16:9"
}`},
	{Name: "single_frame2video_1080p", Method: "POST", Path: "/v1/jimeng/create", Body: `{
  "req_key": "jimeng_i2v_first_v30_1080",
  "prompt": "千军万马",
  "image_urls": [
    "https://img.alicdn.com/imgextra/i1/O1CN01gDEY8M1W114Hi3XcN_!!6000000002727-0-tps-1024-406.jpg"
  ],
  "frames": 121,
  "aspect_ratio": "16:9"
}`},
	{Name: "first_last_frames2video_720p", Method: "POST", Path: "/v1/jimeng/create", Body: `{
  "req_key": "jimeng_i2v_first_tail_v30",
  "prompt": "千军万马",
  "image_urls": [
    "https://img.alicdn.com/imgextra/i1/O1CN01gDEY8M1W114Hi3XcN_!!6000000002727-0-tps-1024-406.jpg",
    "https://img.alicdn.com/imgextra/i2/O1CN01ktT8451iQutqReELT_!!6000000004408-0-tps-689-487.jpg"
  ],
  "frames": 121,
  "aspect_ratio": "16:9"
}`},
	{Name: "first_last_frames2video_1080p", Method: "POST", Path: "/v1/jimeng/create", Body: `{
  "req_key": "jimeng_i2v_first_tail_v30_1080",
  "prompt": "千军万马",
  "image_urls": [
    "https://img.alicdn.com/imgextra/i1/O1CN01gDEY8M1W114Hi3XcN_!!6000000002727-0-tps-1024-406.jpg",
    "https://img.alicdn.com/imgextra/i2/O1CN01ktT8451iQutqReELT_!!6000000004408-0-tps-689-487.jpg"
  ],
  "frames": 121,
  "aspect_ratio": "16:9"
}`},
	{Name: "query", Method: "POST", Path: "/v1/jimeng/query", Body: `{
  "task_id": "7491596536074305586"
}`},
}
