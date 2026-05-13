package intake

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/jonathanpopham/estate-agent/internal/work"
)

type ErrorPayload struct {
	Service     string `json:"service"`
	Environment string `json:"environment"`
	Severity    string `json:"severity"`
	Error       string `json:"error"`
	Stack       string `json:"stack"`
}

func ErrorToWorkItem(payload ErrorPayload) work.Item {
	service := valueOr(payload.Service, "unknown-service")
	environment := valueOr(payload.Environment, "unknown-env")
	severity := valueOr(payload.Severity, "unknown")
	message := valueOr(payload.Error, "unknown error")
	stack := payload.Stack
	fingerprint := fingerprint(service, environment, message, stack)

	return work.Item{
		ID:       "error:" + fingerprint,
		Kind:     work.KindBug,
		Source:   "error",
		Title:    fmt.Sprintf("[%s] %s: %s", environment, service, message),
		Body:     stack,
		Severity: severity,
	}
}

func fingerprint(parts ...string) string {
	hash := sha256.New()
	for _, part := range parts {
		_, _ = hash.Write([]byte(part))
		_, _ = hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))[:16]
}

func valueOr(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
