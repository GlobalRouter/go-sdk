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
		if err := client.RequestJSON(ctx, http.MethodPost, "/v1/tasks", payload, &response); err != nil {
			log.Fatal(err)
		}
		exampleutil.PrintCaptured("Seedance native video tasks", item.Name, capture, before, response)
	}
	exampleutil.EnsureSelection(matched)
}

var examples = []example{
	{Name: "text2video", Body: `{
  "model": "doubao-seedance-1-0-pro-250528",
  "content": [
    {
      "type": "text",
      "text": "写实风格，晴朗的蓝天之下，一大片白色的雏菊花田，镜头逐渐拉近，最终定格在一朵雏菊花的特写上，花瓣上有几颗晶莹的露珠"
    }
  ],
  "duration": 5,
  "ratio": "16:9"
}`},
	{Name: "single_frame2video", Body: `{
  "model": "doubao-seedance-1-0-pro-250528",
  "content": [
    {
      "type": "text",
      "text": "女孩抱着狐狸，女孩睁开眼，温柔地看向镜头，狐狸友善地抱着，镜头缓缓拉出，女孩的头发被风吹动"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/i2v_foxrgirl.png"
      }
    }
  ],
  "duration": 5,
  "ratio": "adaptive"
}`},
	{Name: "first_last_frames2video", Body: `{
  "model": "doubao-seedance-1-0-pro-250528",
  "content": [
    {
      "type": "text",
      "text": "360度环绕运镜"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/seepro_first_frame.jpeg"
      },
      "role": "first_frame"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/seepro_last_frame.jpeg"
      },
      "role": "last_frame"
    }
  ],
  "duration": 5,
  "ratio": "adaptive"
}`},
	{Name: "full_reference2video", Body: `{
  "model": "doubao-seedance-2-0-260128",
  "content": [
    {
      "type": "text",
      "text": "全程使用视频1的第一视角构图，全程使用音频1作为背景音乐。第一人称视角果茶宣传广告，seedance牌「苹苹安安」苹果果茶限定款"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/r2v_tea_pic1.jpg"
      },
      "role": "reference_image"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/r2v_tea_pic2.jpg"
      },
      "role": "reference_image"
    },
    {
      "type": "video_url",
      "video_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_video/r2v_tea_video1.mp4"
      },
      "role": "reference_video"
    },
    {
      "type": "audio_url",
      "audio_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_audio/r2v_tea_audio1.mp3"
      },
      "role": "reference_audio"
    }
  ],
  "duration": 11,
  "ratio": "16:9",
  "generate_audio": true
}`},
	{Name: "reference_image2video", Body: `{
  "model": "doubao-seedance-1-0-lite-i2v-250428",
  "content": [
    {
      "type": "text",
      "text": "[图1]戴着眼镜穿着蓝色T恤的男生和[图2]的柯基小狗，坐在[图3]的草坪上，3D卡通风格"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/seelite_ref_1.png"
      },
      "role": "reference_image"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/seelite_ref_2.png"
      },
      "role": "reference_image"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/seelite_ref_3.png"
      },
      "role": "reference_image"
    }
  ],
  "duration": 5,
  "ratio": "16:9"
}`},
	{Name: "generate_audio", Body: `{
  "model": "doubao-seedance-1-5-pro-251215",
  "content": [
    {
      "type": "text",
      "text": "女孩抱着狐狸，女孩睁开眼，温柔地看向镜头，狐狸友善地抱着，镜头缓缓拉出，女孩的头发被风吹动，可以听到风声"
    },
    {
      "type": "image_url",
      "image_url": {
        "url": "https://ark-project.tos-cn-beijing.volces.com/doc_image/i2v_foxrgirl.png"
      }
    }
  ],
  "duration": 5,
  "ratio": "adaptive",
  "generate_audio": true
}`},
}
