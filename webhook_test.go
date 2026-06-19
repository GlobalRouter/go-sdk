package globalrouter

import "testing"

func TestVerifyWebhookSignature(t *testing.T) {
	payload := []byte(`{"task_id":"task_1"}`)
	if !VerifyWebhookSignature("secret", payload, "sha256=fe7a57ad6d1612d1e826e5cdf3a4dc11b68d86caf47e865a6c8f38fdaf81a166") {
		t.Fatal("sha256 signature should verify")
	}
	if !VerifyWebhookSignature("secret", payload, "t=1780880000,v1=fac3b9ecc70114089dcadb0e6bd74c4fbab6d614a1a1f07cb2e3f1238123c275") {
		t.Fatal("timestamped signature should verify")
	}
	if !VerifyWebhookSignature("secret", payload, "t=1780880000, v1=fac3b9ecc70114089dcadb0e6bd74c4fbab6d614a1a1f07cb2e3f1238123c275") {
		t.Fatal("timestamped signature with comma-space should verify")
	}
	if VerifyWebhookSignature("secret", payload, "sha256=bad") {
		t.Fatal("bad signature should not verify")
	}
}
