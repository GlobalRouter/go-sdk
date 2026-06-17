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
	timestampedPayload := []byte(timestamp + ".")
	timestampedPayload = append(timestampedPayload, payload...)
	if !verifyWebhookSignatureAt("secret", payload, "t="+timestamp+",v1="+sign("secret", timestampedPayload), now) {
		t.Fatal("timestamped signature should verify")
	}
	if VerifyWebhookSignature("secret", payload, "t=bad,v1=fac3b9ecc70114089dcadb0e6bd74c4fbab6d614a1a1f07cb2e3f1238123c275") {
		t.Fatal("malformed timestamp should not verify")
	}
	staleTimestamp := strconv.FormatInt(now.Add(-webhookSignatureTolerance-time.Second).Unix(), 10)
	stalePayload := []byte(staleTimestamp + ".")
	stalePayload = append(stalePayload, payload...)
	if verifyWebhookSignatureAt("secret", payload, "t="+staleTimestamp+",v1="+sign("secret", stalePayload), now) {
		t.Fatal("stale timestamped signature should not verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}
