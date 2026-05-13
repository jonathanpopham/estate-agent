package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const OpenRouterChatCompletionsURL = "https://openrouter.ai/api/v1/chat/completions"

type OpenRouterClient struct {
	APIKey  string
	Model   string
	Referer string
	Title   string
	HTTP    *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Provider  *Provider `json:"provider,omitempty"`
	Temp      *float64  `json:"temperature,omitempty"`
	MaxTokens *int      `json:"max_tokens,omitempty"`
}

type Provider struct {
	AllowFallbacks bool `json:"allow_fallbacks"`
}

func (c OpenRouterClient) NewChatRequest(ctx context.Context, input ChatRequest) (*http.Request, error) {
	if c.APIKey == "" {
		return nil, errors.New("openrouter api key is required")
	}
	if input.Model == "" {
		input.Model = c.Model
	}
	if input.Model == "" {
		return nil, errors.New("openrouter model is required")
	}
	if len(input.Messages) == 0 {
		return nil, errors.New("at least one message is required")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal chat request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, OpenRouterChatCompletionsURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", "Bearer "+c.APIKey)
	req.Header.Set("content-type", "application/json")
	if c.Referer != "" {
		req.Header.Set("http-referer", c.Referer)
	}
	if c.Title != "" {
		req.Header.Set("x-openrouter-title", c.Title)
	}
	return req, nil
}
