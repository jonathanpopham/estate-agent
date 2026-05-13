# Architecture

Estate Agent is a cloud-neutral Go service for autonomous software estate management.

## Control Plane

The control plane accepts:

- GitHub issues and webhooks
- production error payloads
- scheduled maintenance triggers
- future cloud log events from AWS, Azure, and GCP

It normalizes each signal into a work item, evaluates whether the request should proceed, and records decisions before any mutation.

## Model Provider

OpenRouter is the first model-provider target. Operators bring their own OpenRouter API key. Provider-specific keys can be managed inside OpenRouter BYOK settings.

The core service should keep model execution behind an interface so future providers can be added without changing issue intake, evals, or cloud deployment code.

## Runtime Profiles

`ESTATE_AGENT_CLOUD` selects a deployment profile:

- `local`
- `aws`
- `azure`
- `gcp`

The application code should stay cloud-neutral. Terraform and runtime environment should handle cloud-specific concerns such as networking, secrets, queues, and logs.

## Safety Model

The service starts in dry-run mode. Mutating actions should be introduced in layers:

1. comment on issues
2. create issues from production errors
3. create branches
4. open draft PRs
5. update PRs after review feedback

Auto-merge is intentionally out of scope until evals and safety gates are strong.

