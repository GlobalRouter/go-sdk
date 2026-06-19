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
	signature := "t=" + timestamp + ",v1=" + sign("secret", append([]byte(timestamp+"."), payload...))
	if !VerifyWebhookSignature("secret", payload, signature) {
		t.Fatal("timestamped signature should verify")
	}
	staleTimestamp := strconv.FormatInt(time.Now().Add(-webhookTimestampTolerance-time.Second).Unix(), 10)
	staleSignature := "t=" + staleTimestamp + ",v1=" + sign("secret", append([]byte(staleTimestamp+"."), payload...))
	if VerifyWebhookSignature("secret", payload, staleSignature) {
		t.Fatal("stale timestamped signature should not verify")
	}
	if VerifyWebhookSignature("secret", payload, "t=bad,v1=bad") {
		t.Fatal("invalid timestamped signature should not verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}
