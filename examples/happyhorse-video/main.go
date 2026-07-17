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
		if err := client.RequestJSON(ctx, http.MethodPost, "/v1/happyhorse/create", payload, &response); err != nil {
			log.Fatal(err)
		}
		exampleutil.PrintCaptured("Happyhorse video tasks", item.Name, capture, before, response)
	}
	exampleutil.EnsureSelection(matched)
}

var examples = []example{
	{Name: "text2video", Body: `{
  "model": "happyhorse-1.0-t2v",
  "input": {
    "prompt": "一座由硬纸板和瓶盖搭建的微型城市，在夜晚焕发出生机。一列硬纸板火车缓缓驶过，小灯点缀其间，照亮前路。"
  },
  "parameters": {
    "resolution": "720P",
    "ratio": "16:9",
    "duration": 5
  }
}`},
	{Name: "single_frame2video", Body: `{
  "model": "happyhorse-1.0-i2v",
  "input": {
    "prompt": "一只猫在草地上奔跑",
    "media": [
      {
        "type": "first_frame",
        "url": "https://cdn.translate.alibaba.com/r/wanx-demo-1.png"
      }
    ]
  },
  "parameters": {
    "resolution": "720P",
    "duration": 5
  }
}`},
	{Name: "reference2video", Body: `{
  "model": "happyhorse-1.0-r2v",
  "input": {
    "prompt": "[Image 1]中身着红色旗袍的女性，镜头先以侧面中景勾勒旗袍修身剪裁与S型曲线，随即切换至低角度仰拍，捕捉她轻抬玉手展开[Image 2]中的折扇的同时，[Image 3]中的流苏耳坠随头部转动轻盈摆动的细节，最后推近至面部特写，定格在她指尖轻点扇骨、眼波流转间的含蓄风情，多视角全方位展现东方韵味。",
    "media": [
      {
        "type": "reference_image",
        "url": "https://example.com/image1.jpg"
      },
      {
        "type": "reference_image",
        "url": "https://example.com/image2.jpg"
      },
      {
        "type": "reference_image",
        "url": "https://example.com/image3.jpg"
      }
    ]
  },
  "parameters": {
    "resolution": "720P",
    "ratio": "16:9",
    "duration": 5
  }
}`},
	{Name: "video_edit", Body: `{
  "model": "happyhorse-1.0-video-edit",
  "input": {
    "prompt": "让视频中的马头人身角色穿上图片中的条纹毛衣",
    "media": [
      {
        "type": "video",
        "url": "https://example.com/video.mp4"
      },
      {
        "type": "reference_image",
        "url": "https://example.com/clothes.jpg"
      }
    ]
  },
  "parameters": {
    "resolution": "720P"
  }
}`},
}
