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
	timestamp := strconv.FormatInt(now.Unix(), 10)
	if !verifyWebhookSignature("secret", payload, timestampedSignature("secret", timestamp, payload), now, webhookSignatureTolerance) {
		t.Fatal("timestamped signature should verify")
	}
	expired := strconv.FormatInt(now.Add(-webhookSignatureTolerance-time.Second).Unix(), 10)
	if verifyWebhookSignature("secret", payload, timestampedSignature("secret", expired, payload), now, webhookSignatureTolerance) {
		t.Fatal("expired timestamped signature should not verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}

func timestampedSignature(secret, timestamp string, payload []byte) string {
	signedPayload := []byte(timestamp + ".")
	signedPayload = append(signedPayload, payload...)
	return "t=" + timestamp + ",v1=" + sign(secret, signedPayload)
}
