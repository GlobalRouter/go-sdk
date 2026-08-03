package main

import "testing"

func TestRealRequestBodyUsesServerSelectedProvider(t *testing.T) {
	request := realRequestBody()

	if request.Model != "doubao-seed-audio-1-0" {
		t.Fatalf("model = %q", request.Model)
	}
	if len(request.References) != 0 {
		t.Fatalf("real request must not use placeholder references: %#v", request.References)
	}
}
