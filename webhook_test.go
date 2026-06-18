package globalrouter

import (
	"strconv"
	"testing"
	"time"
)

func TestVerifyWebhookSignature(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	if !VerifyWebhookSignature("secret", payload, "sha256=fe7a57ad6d1612d1e826e5cdf3a4dc11b68d86caf47e865a6c8f38fdaf81a166") {
		t.Fatal("sha256 signature should verify")
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signedPayload := append([]byte(timestamp+"."), payload...)
	if !VerifyWebhookSignature("secret", payload, "t="+timestamp+",v1="+sign("secret", signedPayload)) {
		t.Fatal("timestamped signature should verify")
	}
	staleTimestamp := strconv.FormatInt(time.Now().Add(-webhookSignatureTolerance-time.Second).Unix(), 10)
	stalePayload := append([]byte(staleTimestamp+"."), payload...)
	if VerifyWebhookSignature("secret", payload, "t="+staleTimestamp+",v1="+sign("secret", stalePayload)) {
		t.Fatal("stale timestamped signature should not verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}
