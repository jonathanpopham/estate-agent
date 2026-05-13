package config

import "testing"

func TestFromEnvDefaults(t *testing.T) {
	t.Setenv("ESTATE_AGENT_ADDR", "")
	t.Setenv("ESTATE_AGENT_CLOUD", "")
	t.Setenv("ESTATE_AGENT_DRY_RUN", "")
	t.Setenv("ESTATE_AGENT_OPENROUTER_MODEL", "")

	cfg := FromEnv()

	if cfg.Addr != "127.0.0.1:8080" {
		t.Fatalf("Addr = %q", cfg.Addr)
	}
	if cfg.Cloud != CloudLocal {
		t.Fatalf("Cloud = %q", cfg.Cloud)
	}
	if !cfg.DryRun {
		t.Fatal("DryRun = false, want true")
	}
	if cfg.OpenRouterModel == "" {
		t.Fatal("OpenRouterModel is empty")
	}
}

func TestFromEnvClouds(t *testing.T) {
	tests := map[string]Cloud{
		"aws":   CloudAWS,
		"azure": CloudAzure,
		"gcp":   CloudGCP,
		"nope":  CloudLocal,
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			t.Setenv("ESTATE_AGENT_CLOUD", input)
			if got := FromEnv().Cloud; got != want {
				t.Fatalf("Cloud = %q, want %q", got, want)
			}
		})
	}
}
