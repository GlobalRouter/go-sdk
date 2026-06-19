package globalrouter

import (
	"fmt"
	"testing"
	"time"
)

func TestVerifyWebhookSignature(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	if !VerifyWebhookSignature("secret", payload, "sha256=fe7a57ad6d1612d1e826e5cdf3a4dc11b68d86caf47e865a6c8f38fdaf81a166") {
		t.Fatal("sha256 signature should verify")
	}
	now := time.Unix(1780880000, 0)
	if !verifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, now), now) {
		t.Fatal("timestamped signature should verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}

func TestVerifyWebhookSignatureRejectsExpiredTimestamp(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	now := time.Unix(1780880000, 0)
	expired := now.Add(-(webhookSignatureTolerance + time.Second))
	if verifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, expired), now) {
		t.Fatal("expired timestamped signature should not verify")
	}
}

func TestVerifyWebhookSignatureRejectsFutureTimestamp(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	now := time.Unix(1780880000, 0)
	future := now.Add(webhookSignatureTolerance + time.Second)
	if verifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, future), now) {
		t.Fatal("future timestamped signature should not verify")
	}
}

func TestVerifyWebhookSignatureRejectsMalformedTimestamp(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	signature := "t=not-a-timestamp,v1=" + sign("secret", []byte("not-a-timestamp."+string(payload)))
	if verifyWebhookSignature("secret", payload, signature, time.Unix(1780880000, 0)) {
		t.Fatal("malformed timestamped signature should not verify")
	}
}

func timestampedSignature(secret string, payload []byte, timestamp time.Time) string {
	timestampString := fmt.Sprintf("%d", timestamp.Unix())
	signedPayload := []byte(timestampString + ".")
	signedPayload = append(signedPayload, payload...)
	return "t=" + timestampString + ",v1=" + sign(secret, signedPayload)
}
