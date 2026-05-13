package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestVerifyGitHubSignature(t *testing.T) {
	secret := "secret"
	body := []byte(`{"hello":"world"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !VerifyGitHubSignature(secret, body, signature) {
		t.Fatal("valid signature rejected")
	}
	if VerifyGitHubSignature(secret, body, "sha256=bad") {
		t.Fatal("invalid signature accepted")
	}
	if !VerifyGitHubSignature("", body, "") {
		t.Fatal("empty secret should allow local development")
	}
}
