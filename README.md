# Estate Agent

Estate Agent is an open source service for autonomous software-estate management.

The intended shape is a long-running service that owns a website or application estate. It listens to operational signals and GitHub issues, uses a BYO-key LLM provider such as OpenRouter, evaluates whether work should proceed, and then drives accepted work through an auditable software delivery loop.

This repo is starting fresh in Go.

## Product Direction

Estate Agent should support:

- GitHub issues as the control surface for bugs, feature requests, and maintenance work
- production error ingestion from cloud logs and alerting systems
- OpenRouter-compatible model execution with bring-your-own API keys
- cloud-neutral deployment across AWS, Azure, and GCP
- test-driven development for core routing, safety, and model-request behavior
- deterministic eval fixtures before the agent can mutate code
- GitHub issues as the project ledger while the system is built

## Current Status

The Go foundation includes:

- typed runtime config
- HTTP health endpoint
- GitHub webhook signature verification
- issue-event routing into normalized work items
- error-payload routing into normalized work items
- OpenRouter chat-completion request shaping
- deterministic eval fixture runner

## Local Development

```bash
go test ./...
go run ./cmd/estate-agent serve
```

Health check:

```bash
curl http://127.0.0.1:8080/health
```

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `ESTATE_AGENT_ADDR` | `127.0.0.1:8080` | HTTP listen address. |
| `ESTATE_AGENT_CLOUD` | `local` | Runtime profile: `local`, `aws`, `azure`, or `gcp`. |
| `ESTATE_AGENT_DRY_RUN` | `true` | Keeps the agent from mutating GitHub or repos. |
| `ESTATE_AGENT_GITHUB_WEBHOOK_SECRET` | unset | Secret for GitHub webhook signature checks. |
| `ESTATE_AGENT_OPENROUTER_API_KEY` | unset | OpenRouter API key supplied by the operator. |
| `ESTATE_AGENT_OPENROUTER_MODEL` | `openai/gpt-4.1-mini` | Default model for planning/eval calls. |
| `ESTATE_AGENT_OPENROUTER_REFERER` | unset | Optional OpenRouter app attribution URL. |
| `ESTATE_AGENT_OPENROUTER_TITLE` | `Estate Agent` | Optional OpenRouter app title. |

## Active Work

The project is tracked in GitHub issues:

- #2 evaluation harness
- #3 OpenRouter BYO-key model execution
- #4 cloud-neutral AWS/Azure/GCP deployment
- #5 Go service foundation

