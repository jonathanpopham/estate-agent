package llm

import (
	"context"
	"encoding/json"
	"testing"
)

func TestOpenRouterChatRequestRequiresAPIKey(t *testing.T) {
	client := OpenRouterClient{Model: "openai/gpt-4.1-mini"}

	_, err := client.NewChatRequest(context.Background(), ChatRequest{
		Messages: []Message{{Role: "user", Content: "plan this issue"}},
	})

	if err == nil {
		t.Fatal("expected missing api key error")
	}
}

func TestOpenRouterChatRequestShape(t *testing.T) {
	client := OpenRouterClient{
		APIKey:  "sk-or-test",
		Model:   "anthropic/claude-sonnet-4.5",
		Referer: "https://github.com/jonathanpopham/estate-agent",
		Title:   "Estate Agent",
	}

	req, err := client.NewChatRequest(context.Background(), ChatRequest{
		Messages: []Message{{Role: "user", Content: "triage this"}},
		Provider: &Provider{AllowFallbacks: false},
	})
	if err != nil {
		t.Fatalf("NewChatRequest: %v", err)
	}

	if req.Method != "POST" {
		t.Fatalf("method = %q", req.Method)
	}
	if got := req.Header.Get("authorization"); got != "Bearer sk-or-test" {
		t.Fatalf("authorization = %q", got)
	}
	if got := req.Header.Get("x-openrouter-title"); got != "Estate Agent" {
		t.Fatalf("x-openrouter-title = %q", got)
	}

	var body ChatRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Model != "anthropic/claude-sonnet-4.5" {
		t.Fatalf("model = %q", body.Model)
	}
	if body.Provider == nil || body.Provider.AllowFallbacks {
		t.Fatalf("provider fallback config = %#v", body.Provider)
	}
}
