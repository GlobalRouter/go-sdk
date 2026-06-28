package globalrouter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

const defaultWebhookSignatureTolerance = 5 * time.Minute

func VerifyWebhookSignature(secret string, payload []byte, signature string) bool {
	return verifyWebhookSignature(secret, payload, signature, time.Now(), defaultWebhookSignatureTolerance)
}

func verifyWebhookSignature(secret string, payload []byte, signature string, now time.Time, tolerance time.Duration) bool {
	if strings.HasPrefix(signature, "sha256=") {
		expected := sign(secret, payload)
		provided := strings.TrimPrefix(signature, "sha256=")
		return secureCompareHex(expected, provided)
	}
	parts := map[string]string{}
	for _, item := range strings.Split(signature, ",") {
		key, value, ok := strings.Cut(item, "=")
		if ok {
			parts[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	timestamp, okT := parts["t"]
	provided, okV := parts["v1"]
	if !okT || !okV {
		return false
	}
	timestampUnix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	signedAt := time.Unix(timestampUnix, 0)
	age := now.Sub(signedAt)
	if age < 0 {
		age = -age
	}
	if age > tolerance {
		return false
	}
	signedPayload := []byte(timestamp + ".")
	signedPayload = append(signedPayload, payload...)
	return secureCompareHex(sign(secret, signedPayload), provided)
}

func sign(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func secureCompareHex(expected, provided string) bool {
	expectedBytes, errExpected := hex.DecodeString(expected)
	providedBytes, errProvided := hex.DecodeString(provided)
	if errExpected != nil || errProvided != nil {
		return false
	}
	return hmac.Equal(expectedBytes, providedBytes)
}
