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
		if err := client.RequestJSON(ctx, http.MethodPost, "/v1/seedream/image/create", payload, &response); err != nil {
			log.Fatal(err)
		}
		exampleutil.PrintCaptured("Seedream image generation", item.Name, capture, before, response)
	}
	exampleutil.EnsureSelection(matched)
}

var examples = []example{
	{Name: "text2image", Body: `{
  "model": "doubao-seedream-5-0-260128",
  "prompt": "充满活力的特写编辑肖像，模特眼神犀利，头戴雕塑感帽子，色彩拼接丰富，眼部焦点锐利，景深较浅，具有Vogue杂志封面的美学风格，采用中画幅拍摄，工作室灯光效果强烈。",
  "size": "2K",
  "output_format": "png",
  "watermark": false
}`},
	{Name: "image2image", Body: `{
  "model": "doubao-seedream-5-0-260128",
  "prompt": "保持模特姿势和液态服装的流动形状不变。将服装材质从银色金属改为完全透明的清水（或玻璃）。透过液态水流，可以看到模特的皮肤细节。光影从反射变为折射。",
  "image": "https://ark-project.tos-cn-beijing.volces.com/doc_image/seedream4_5_imageToimage.png",
  "size": "2K",
  "output_format": "png",
  "watermark": false
}`},
}
