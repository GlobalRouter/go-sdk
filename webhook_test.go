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
	signature := timestampedSignature("secret", payload, now)
	if !verifyWebhookSignatureAt("secret", payload, signature, now) {
		t.Fatal("timestamped signature should verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}

func TestVerifyWebhookSignatureRejectsExpiredTimestamp(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	now := time.Unix(1780880000, 0)
	expired := now.Add(-webhookTimestampTolerance - time.Second)
	signature := timestampedSignature("secret", payload, expired)

	if verifyWebhookSignatureAt("secret", payload, signature, now) {
		t.Fatal("expired timestamped signature should not verify")
	}
}

func TestVerifyWebhookSignatureRejectsMalformedTimestamp(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	signature := "t=not-a-timestamp,v1=" + sign("secret", append([]byte("not-a-timestamp."), payload...))

	if verifyWebhookSignatureAt("secret", payload, signature, time.Unix(1780880000, 0)) {
		t.Fatal("malformed timestamped signature should not verify")
	}
}

func timestampedSignature(secret string, payload []byte, timestamp time.Time) string {
	timestampValue := strconv.FormatInt(timestamp.Unix(), 10)
	signedPayload := append([]byte(timestampValue+"."), payload...)
	return "t=" + timestampValue + ",v1=" + sign(secret, signedPayload)
}
