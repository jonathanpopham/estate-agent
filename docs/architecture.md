# Architecture

Estate Agent is a resident software-maintenance agent. It should eventually be deployable as a small service on AWS or Azure, with narrow credentials to one or more repositories and observability sources.

## Core Loop

```text
signals -> work item -> triage -> plan -> implement -> verify -> review -> PR -> deploy gate
```

Signals can come from:

- GitHub issues labeled `estate:bug` or `estate:feature`
- production error logs
- Azure Application Insights alerts
- AWS CloudWatch alarms
- Sentry-style webhook payloads
- scheduled health checks

## Work Item Types

### Bug

A bug usually begins as either a production error or a GitHub issue. The agent should require a reproduction path or a failing test before implementation.

### Feature

A feature begins as a GitHub issue. The agent should decide whether the issue is actionable, comment with clarifying questions when needed, then propose acceptance criteria before implementation.

## Agent Modes

### Dry Run

Default. The agent writes a plan only. No comments, branches, issues, or PRs are created.

### Commenter

The agent can comment on issues with plans and clarification questions.

### Builder

The agent can create branches, commit changes, run tests, and open draft PRs.

### Maintainer

The agent can manage recurring upkeep tasks such as dependency updates, dead-code cleanup, docs refreshes, and low-risk refactors.

Auto-merge is intentionally outside the MVP.

## State

The MVP uses SQLite for local state:

- work item id
- source
- kind
- repository
- issue number
- latest payload
- timestamps

Later deployments can swap this for Postgres, DynamoDB, Azure Table Storage, or Durable Functions state.

## Safety Gates

Before opening or marking a PR ready, the agent should check:

- diff size
- changed path allowlist/denylist
- CI status
- test output
- dependency/security impact
- whether secrets or credentials were touched
- whether the issue has enough acceptance criteria

Default high-risk paths:

- `.github/`
- infrastructure directories
- auth and billing modules
- migrations
- deployment manifests
- secret/config files

## Optional Code Intelligence

Supermodel-style graph context can be one tool in the context layer:

- blast radius for changed symbols
- callers/callees
- dependency edges
- dead-code candidates
- architecture/domain ownership

Estate Agent should keep that as an adapter, not a hard dependency.

## Deployment Model

The expected production shape is:

```text
GitHub Webhooks ─┐
Cloud Logs ──────┼─> HTTPS endpoint -> queue -> worker -> isolated checkout -> PR
Scheduler ───────┘
```

Use a queue between intake and execution so webhook handling stays fast and agent work can be retried, rate-limited, or cancelled.

