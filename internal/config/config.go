package config

import (
	"os"
	"strings"
)

type Cloud string

const (
	CloudLocal Cloud = "local"
	CloudAWS   Cloud = "aws"
	CloudAzure Cloud = "azure"
	CloudGCP   Cloud = "gcp"
)

type Config struct {
	Addr                  string
	Cloud                 Cloud
	DryRun                bool
	GitHubWebhookSecret   string
	OpenRouterAPIKey      string
	OpenRouterModel       string
	OpenRouterReferer     string
	OpenRouterTitle       string
	OpenRouterAllowBYOKFB bool
}

func FromEnv() Config {
	return Config{
		Addr:                  env("ESTATE_AGENT_ADDR", "127.0.0.1:8080"),
		Cloud:                 parseCloud(env("ESTATE_AGENT_CLOUD", "local")),
		DryRun:                envBool("ESTATE_AGENT_DRY_RUN", true),
		GitHubWebhookSecret:   os.Getenv("ESTATE_AGENT_GITHUB_WEBHOOK_SECRET"),
		OpenRouterAPIKey:      os.Getenv("ESTATE_AGENT_OPENROUTER_API_KEY"),
		OpenRouterModel:       env("ESTATE_AGENT_OPENROUTER_MODEL", "openai/gpt-4.1-mini"),
		OpenRouterReferer:     os.Getenv("ESTATE_AGENT_OPENROUTER_REFERER"),
		OpenRouterTitle:       env("ESTATE_AGENT_OPENROUTER_TITLE", "Estate Agent"),
		OpenRouterAllowBYOKFB: envBool("ESTATE_AGENT_OPENROUTER_ALLOW_BYOK_FALLBACKS", false),
	}
}

func parseCloud(value string) Cloud {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "aws":
		return CloudAWS
	case "azure":
		return CloudAzure
	case "gcp":
		return CloudGCP
	default:
		return CloudLocal
	}
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return fallback
	}
	return value == "1" || value == "true" || value == "yes" || value == "on"
}
