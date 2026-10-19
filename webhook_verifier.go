package sevk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const webhookToleranceSeconds = 5 * 60

// ErrWebhookVerification is returned when a webhook delivery fails verification.
type ErrWebhookVerification struct {
	Reason string
}

func (e *ErrWebhookVerification) Error() string {
	return "webhook verification failed: " + e.Reason
}

// WebhookVerifier verifies the signature on incoming sevk webhook deliveries.
//
// sevk uses the Svix-compatible signature scheme (svix-id, svix-timestamp,
// svix-signature headers), so this verifier or the official svix-go package
// can be used. sevk issues canonical whsec_ base64 secrets that can be passed
// directly to either verifier.
type WebhookVerifier struct {
	key []byte
}

// NewWebhookVerifier constructs a verifier from the endpoint signing secret.
func NewWebhookVerifier(secret string) (*WebhookVerifier, error) {
	if secret == "" {
		return nil, &ErrWebhookVerification{Reason: "webhook secret is required"}
	}
	var key []byte
	if strings.HasPrefix(secret, "whsec_") {
		decoded, err := base64.StdEncoding.DecodeString(secret[len("whsec_"):])
		if err != nil || len(decoded) == 0 {
			return nil, &ErrWebhookVerification{Reason: "invalid base64 secret"}
		}
		key = decoded
	} else {
		key = []byte(secret)
	}
	return &WebhookVerifier{key: key}, nil
}

// Verify checks the headers and signature on a webhook delivery and returns
// the parsed JSON payload. Returns *ErrWebhookVerification on missing
// headers, timestamp drift, or signature mismatch.
func (w *WebhookVerifier) Verify(payload []byte, headers http.Header) (any, error) {
	msgID := headerValue(headers, "svix-id")
	msgTimestamp := headerValue(headers, "svix-timestamp")
	msgSignature := headerValue(headers, "svix-signature")

	if msgID == "" || msgTimestamp == "" || msgSignature == "" {
		return nil, &ErrWebhookVerification{Reason: "missing required svix headers"}
	}

	timestamp, err := strconv.ParseInt(msgTimestamp, 10, 64)
	if err != nil || !asciiDigits(msgTimestamp) {
		return nil, &ErrWebhookVerification{Reason: "invalid timestamp header"}
	}
	now := time.Now().Unix()
	if timestamp < now-webhookToleranceSeconds || timestamp > now+webhookToleranceSeconds {
		return nil, &ErrWebhookVerification{Reason: "timestamp outside tolerance window"}
	}

	toSign := msgID + "." + msgTimestamp + "." + string(payload)
	mac := hmac.New(sha256.New, w.key)
	mac.Write([]byte(toSign))
	expected := "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	for _, candidate := range strings.Split(msgSignature, " ") {
		if hmac.Equal([]byte(candidate), []byte(expected)) {
			var result any
			if err := json.Unmarshal(payload, &result); err != nil {
				return nil, &ErrWebhookVerification{Reason: "invalid JSON payload"}
			}
			return result, nil
		}
	}

	return nil, &ErrWebhookVerification{Reason: "signature mismatch"}
}

func headerValue(headers http.Header, name string) string {
	for key, values := range headers {
		if strings.EqualFold(key, name) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func asciiDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// IsWebhookVerificationError reports whether err is a webhook verification failure.
func IsWebhookVerificationError(err error) bool {
	var target *ErrWebhookVerification
	return errors.As(err, &target)
}
