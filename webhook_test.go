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
	timestamp := time.Now().Unix()
	timestampedPayload := []byte(fmt.Sprintf("%d.", timestamp))
	timestampedPayload = append(timestampedPayload, payload...)
	if !VerifyWebhookSignature("secret", payload, fmt.Sprintf("t=%d,v1=%s", timestamp, sign("secret", timestampedPayload))) {
		t.Fatal("timestamped signature should verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}

func TestVerifyWebhookSignatureRejectsStaleTimestamp(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	timestamp := time.Now().Add(-DefaultWebhookSignatureTolerance - time.Second).Unix()
	timestampedPayload := []byte(fmt.Sprintf("%d.", timestamp))
	timestampedPayload = append(timestampedPayload, payload...)
	signature := fmt.Sprintf("t=%d,v1=%s", timestamp, sign("secret", timestampedPayload))

	if VerifyWebhookSignature("secret", payload, signature) {
		t.Fatal("stale timestamped signature should not verify")
	}
	if !VerifyWebhookSignatureWithTolerance("secret", payload, signature, DefaultWebhookSignatureTolerance+2*time.Second) {
		t.Fatal("timestamped signature should verify with a larger caller-provided tolerance")
	}
}

func TestVerifyWebhookSignatureRejectsInvalidTimestamp(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	timestampedPayload := []byte("not-a-timestamp.")
	timestampedPayload = append(timestampedPayload, payload...)
	signature := fmt.Sprintf("t=not-a-timestamp,v1=%s", sign("secret", timestampedPayload))

	if VerifyWebhookSignature("secret", payload, signature) {
		t.Fatal("timestamped signature with invalid timestamp should not verify")
	}
}
