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
	now := time.Unix(1780880000, 0)
	if !verifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, now), now, defaultWebhookSignatureTolerance) {
		t.Fatal("timestamped signature should verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
	staleTimestamp := now.Add(-defaultWebhookSignatureTolerance - time.Second)
	if verifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, staleTimestamp), now, defaultWebhookSignatureTolerance) {
		t.Fatal("stale timestamped signature should not verify")
	}
	if verifyWebhookSignature("secret", payload, "t=not-a-timestamp,v1=bad", now, defaultWebhookSignatureTolerance) {
		t.Fatal("malformed timestamped signature should not verify")
	}
}

func timestampedSignature(secret string, payload []byte, timestamp time.Time) string {
	timestampText := strconv.FormatInt(timestamp.Unix(), 10)
	signedPayload := []byte(timestampText + ".")
	signedPayload = append(signedPayload, payload...)
	return "t=" + timestampText + ",v1=" + sign(secret, signedPayload)
}
