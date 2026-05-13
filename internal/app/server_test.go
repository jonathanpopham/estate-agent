package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonathanpopham/estate-agent/internal/config"
)

func TestHealth(t *testing.T) {
	handler := NewServer(config.Config{
		Addr:   "127.0.0.1:0",
		Cloud:  config.CloudLocal,
		DryRun: true,
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["ok"] != true {
		t.Fatalf("ok = %v, want true", body["ok"])
	}
	if body["cloud"] != string(config.CloudLocal) {
		t.Fatalf("cloud = %v, want %s", body["cloud"], config.CloudLocal)
	}
}
