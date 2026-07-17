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
		exampleutil.PrintCaptured("Admin price rules", item.Name, capture, before, response)
	}
	exampleutil.EnsureSelection(matched)
}

var examples = []example{
	{Name: "update_price_rule", Method: "POST", Path: "/admin/price-rules/price_rule_123", Body: `{
  "input_price_micros_per_million": 200000,
  "output_price_micros_per_million": 800000,
  "status": 1
}`},
	{Name: "create_price_rule", Method: "POST", Path: "/admin/price-rules/create", Body: `{
  "model_key": "llm_model_gpt4o",
  "billing_family": "text_token",
  "billing_strategy": "text_io",
  "input_token_min": 0,
  "input_price_micros_per_million": 150000,
  "output_price_micros_per_million": 600000,
  "effective_at": "2026-06-02 00:00:00",
  "status": 1
}`},
	{Name: "business_override_price_rule", Method: "POST", Path: "/admin/price-rules/business-override", Body: `{
  "model_key": "llm_model_gpt4o",
  "billing_family": "text_token",
  "billing_strategy": "text_io",
  "input_price_micros_per_million": 150000,
  "output_price_micros_per_million": 600000
}`},
}
