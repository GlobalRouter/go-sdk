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
	if !VerifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, time.Now().Unix())) {
		t.Fatal("timestamped signature should verify")
	}
	if VerifyWebhookSignature("secret", payload, timestampedSignature("secret", payload, time.Now().Add(-webhookTimestampTolerance-time.Second).Unix())) {
		t.Fatal("expired timestamped signature should not verify")
	}
	if VerifyWebhookSignature("secret", payload, "t=bad,v1=fac3b9ecc70114089dcadb0e6bd74c4fbab6d614a1a1f07cb2e3f1238123c275") {
		t.Fatal("malformed timestamped signature should not verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}

func timestampedSignature(secret string, payload []byte, timestamp int64) string {
	timestampString := fmt.Sprint(timestamp)
	signedPayload := []byte(timestampString + ".")
	signedPayload = append(signedPayload, payload...)
	return "t=" + timestampString + ",v1=" + sign(secret, signedPayload)
}
