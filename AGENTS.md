# Estate Agent Principles

This file is for coding agents and maintainers working on this repository. It captures the intended product shape and the architectural principles that should guide implementation.

## Product Identity

Estate Agent is a self-hosted website steward.

It should run the website, observe the website, receive user reports, turn runtime failures into issues, reproduce problems locally, repair them, test them, deploy them through staging, and verify that production is healthy again.

The long-term goal is not a generic AI coding bot. The goal is a website that comes with its own maintainer.

## Primary Deployment Model

The first-class deployment target is an Estate Node: one trusted server that runs the website and the agent side by side.

The node should own:

- the public website containers
- a staging environment
- the private agent UI and API
- a durable queue and worker system
- a local ledger database
- git mirrors and disposable workspaces
- containerized dev and test environments
- logs, metrics, traces, and run artifacts
- secrets needed for GitHub, model providers, email, and deployment

Humans should access the private agent surface over Tailscale. Public ingress should be limited to the website itself and any intentionally exposed intake endpoints.

Cloud support should initially mean running the same Estate Node on a VM, not rewriting the core around cloud-managed services. AWS, Azure, GCP, Hetzner, or a local machine should all be able to host the same node.

## Core Loop

The core product loop is:

```text
host site
observe site
receive user messages and runtime errors
create or update issues
triage issues
reproduce in a local container
patch the code
run tests
deploy to staging
verify staging
merge or promote
deploy to production
observe production
close the issue or ask for more information
```

Runtime failure should automatically become software work. A server throw, failed health check, repeated error log, or user email should all land in the same issue system.

## Issue Sources

Estate Agent should ingest issues from multiple sources:

- production errors and structured logs
- user emails to a configured support/reporting address
- uptime checks and health probes
- deployment failures
- GitHub issues, when enabled
- manual entries from the private UI

GitHub is useful for code review and PRs, but it should not be the only product surface. The Estate Node needs its own ledger and issue view.

## Triage States

Initial issue triage should classify each issue as one of:

- `already_solved`: the reported problem no longer reproduces or was fixed by a recent deploy
- `actionable`: there is enough information to reproduce, test, or safely attempt a fix
- `needs_info`: the issue may be real, but the reporter or operator needs to provide missing details
- `close`: spam, unsupported request, duplicate, by-design behavior, or unsafe/destructive request

Do not collapse `needs_info` into `close`. User reports are often incomplete, and the agent should ask precise follow-up questions when that would move the issue forward.

## Runtime Error Handling

Application logs should be shaped for repair, not only for debugging.

Error events should preserve:

- error message
- stack trace
- route, handler, job, or component name
- request ID and trace ID
- release commit SHA
- container image version
- environment
- sanitized request metadata
- hashed user or session reference when useful
- first seen, last seen, and occurrence count
- reproduction hints when available

Errors should be fingerprinted and deduplicated. One bug happening many times should become one issue with severity and frequency data, not hundreds of noisy issues.

Secrets and PII must be redacted before logs are sent to a model or displayed broadly.

## Reproduction Before Repair

The agent should try to prove the bug exists before editing code.

For runtime errors, the expected flow is:

```text
identify failing release commit
checkout the repo in a disposable workspace
build or start the local dev container
attempt to reproduce the failure
write or identify a failing test when practical
patch the code
run the configured test suite
deploy to staging
verify the issue against staging
only then proceed toward production
```

Guessing from a stack trace is allowed only as a last resort and should be recorded as such in the ledger.

## Staging Is Required

Production should not be the first place a fix is exercised.

Every code-changing repair should have a staging step unless the repository policy explicitly marks the change class as safe to skip staging. Staging should be close enough to production to catch configuration, routing, container, and integration failures.

The agent should be able to compare production and staging by:

- deployed commit
- image version
- environment config class
- health checks
- smoke tests
- relevant logs and traces

## Safety Modes

The system should gain power in layers:

- `dry_run`: create issues, classify, plan, and record decisions only
- `assist`: create branches or draft PRs, but require human merge/deploy
- `autonomous_low_risk`: auto-merge narrow low-risk fixes after tests and staging verification
- `autonomous_full`: broader self-healing with canaries, rollback, and strong policy gates

The first implementation should default to `dry_run`.

Auto-merge to production must be gated by policy, tests, staging verification, rollback availability, and observability.

## Architecture Principles

Keep the center of the system around durable work, not HTTP handlers.

The preferred shape is:

```text
signal
-> normalized work item
-> durable queue
-> worker run
-> policy and eval gates
-> plan
-> action attempts
-> tests and verification
-> ledger
```

Important boundaries:

- HTTP handlers validate, normalize, enqueue, and return.
- Workers perform slow work such as model calls, git operations, test runs, and GitHub actions.
- Policy decides what is allowed.
- Action adapters execute allowed operations.
- The ledger records every meaningful decision and action.
- Core logic should not know which cloud or VM provider is hosting the node.
- Cloud, container, GitHub, email, model, and secret systems should be adapters behind interfaces.

## Local Execution Model

The node should use isolated per-run execution.

Each run should have:

- a unique run ID
- a disposable workspace
- a known repo commit
- a bounded container or process environment
- CPU, memory, wall-time, and disk limits
- explicit network mode
- an allowlist of secrets
- captured stdout, stderr, logs, artifacts, and test results
- automatic cleanup rules

Long-lived repo mirrors and dependency caches are useful, but workspaces should be disposable.

## Observability

Observability is core product behavior, not an add-on.

Estate Agent should provide three kinds of observability:

1. Operational telemetry for the node itself.
2. Product auditability for every issue and run.
3. Debuggability for failed repairs and unsafe decisions.

Every run should carry correlation fields:

- `run_id`
- `issue_id`
- `work_item_id`
- `repo`
- `commit`
- `deployment_id`
- `trace_id`
- `request_id`
- `model_request_id`
- `action_attempt_id`

Required metrics include:

- queue depth
- issue ingestion rate
- triage outcomes
- run duration
- run success and failure counts
- reproduction success rate
- test pass and fail counts
- staging verification outcomes
- production rollback count
- model latency
- token usage
- cost estimate
- policy denials
- eval failures
- GitHub API failures

The private UI should show a readable timeline for each issue and run. Operators should not need to read raw logs to understand what happened.

## Ledger Requirements

The ledger is the trust backbone.

It should record:

- source signal
- normalized issue or work item
- deduplication fingerprint
- triage decision and reason
- policy decisions
- model provider, model, and configuration
- prompts and responses when safe to store
- commands run
- files changed
- tests run
- staging deployment result
- production deployment result
- verification result
- user or operator replies
- final close reason

The ledger should support idempotency. Replayed webhooks, repeated logs, repeated emails, or retried jobs must not create inconsistent state.

## Configuration Principles

Configuration should exist at two levels: node-level and repo-level.

Node-level configuration controls the server:

```text
ESTATE_AGENT_HOME
ESTATE_AGENT_ADDR
ESTATE_AGENT_MODE
ESTATE_AGENT_PUBLIC_SITE_URL
ESTATE_AGENT_STAGING_SITE_URL
ESTATE_AGENT_CONTAINER_RUNTIME
ESTATE_AGENT_WORKSPACE_ROOT
ESTATE_AGENT_REPO_CACHE_ROOT
ESTATE_AGENT_MAX_PARALLEL_RUNS
ESTATE_AGENT_RUN_TIMEOUT
ESTATE_AGENT_DEFAULT_NETWORK_MODE
ESTATE_AGENT_STORAGE_URL
ESTATE_AGENT_QUEUE_URL
```

Repo-level configuration controls what the agent may do:

```yaml
repos:
  - github: owner/name
    production_url: https://example.com
    staging_url: https://staging.example.com
    default_branch: main
    test_commands:
      - go test ./...
    smoke_tests:
      - GET /
      - GET /health
    allowed_actions:
      - comment
      - label
      - branch
      - draft_pr
    protected_paths:
      - .github/workflows/**
      - infra/prod/**
    approval_required_for:
      - auth
      - billing
      - secrets
      - migrations
      - destructive_data_changes
```

Prefer explicit configuration over hidden behavior. Defaults should be conservative.

## Model Provider Principles

Model execution should be provider-neutral, with OpenRouter as the first adapter.

The model layer should support:

- bring-your-own API keys
- request-scoped secret references
- configurable model, temperature, token limits, and routing options
- clear failure when required keys are missing
- no API keys in logs, prompts, traces, or committed files
- deterministic offline eval fixtures that do not require model calls

Provider choice should not leak into intake, policy, queueing, or git execution.

## Human Control

The agent should make the software estate calmer, not less understandable.

Humans need:

- a private Tailscale-accessible UI
- current issues and runs
- pending decisions
- proposed plans
- diffs and test results
- staging links
- approve, reject, pause, retry, and rollback controls
- clear reasons for closes and escalations

When the agent replies to a user email, the reply should be specific, honest about uncertainty, and tied to the issue state.

## Definition Of Done

A feature is not done merely because code compiles.

For core behavior, done means:

- the signal becomes a durable issue or work item
- duplicate events are handled intentionally
- the run is visible in the ledger
- policy and safety gates are explicit
- logs, metrics, and correlation IDs exist
- tests cover the intended behavior
- staging or verification behavior is defined
- failure modes are observable
- secrets are not exposed
- the human can understand what happened

