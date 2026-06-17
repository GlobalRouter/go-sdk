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
	if !VerifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, time.Now())) {
		t.Fatal("timestamped signature should verify")
	}
	if VerifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, time.Now().Add(-webhookSignatureTolerance-time.Second))) {
		t.Fatal("expired timestamped signature should not verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}

func timestampedSignature(secret string, payload []byte, timestamp time.Time) string {
	timestampValue := strconv.FormatInt(timestamp.Unix(), 10)
	signedPayload := []byte(timestampValue + ".")
	signedPayload = append(signedPayload, payload...)
	return "t=" + timestampValue + ",v1=" + sign(secret, signedPayload)
}
