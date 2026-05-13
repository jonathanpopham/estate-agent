# Estate Agent

An autonomous steward for software estates.

Estate Agent is a server-side agent that listens for operational signals and GitHub issues, turns them into scoped work items, and drives them through a software delivery loop. The first MVP is intentionally conservative: it classifies incoming work, stores state, and writes an implementation plan in dry-run mode.

The long-term goal is a resident agent that can own an app or website end to end:

- production errors become bug issues
- feature requests are filed as GitHub issues
- the agent decides whether work should proceed
- accepted work gets planned, implemented, tested, reviewed, and opened as a PR
- risky work escalates to a human instead of mutating production blindly

## MVP Shape

```text
GitHub issue webhook ─┐
                      ├─> Estate Agent server ─> work item ─> SDLC plan
Error/log webhook ────┘                         └─> SQLite state
```

Current capabilities:

- Receives GitHub issue webhooks for `bug`, `estate:bug`, `feature`, and `estate:feature` labels.
- Receives generic error payloads from log pipelines.
- Normalizes both into a `WorkItem`.
- Stores work state in SQLite.
- Generates a structured implementation plan.
- Runs in dry-run mode by default.

Planned next capabilities:

- Post plans back to GitHub issues.
- Create bug issues from production errors.
- Run a coding agent in an isolated checkout.
- Push branches and open draft PRs.
- Require CI, review, diff-size, and path-safety gates before readiness.
- Deploy on AWS or Azure with Terraform.

## Quickstart

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -e ".[server]"

export ESTATE_AGENT_DRY_RUN=true
export ESTATE_AGENT_GITHUB_WEBHOOK_SECRET=local-dev-secret

estate-agent serve --reload
```

Health check:

```bash
curl http://127.0.0.1:8080/health
```

Send a local error payload:

```bash
curl -X POST http://127.0.0.1:8080/ingest/error \
  -H 'content-type: application/json' \
  -d '{
    "service": "checkout-web",
    "environment": "prod",
    "error": "TypeError: Cannot read properties of undefined",
    "stack": "at calculateTax (src/tax.ts:42:11)",
    "severity": "high"
  }'
```

## Configuration

Environment variables:

| Variable | Default | Purpose |
| --- | --- | --- |
| `ESTATE_AGENT_DRY_RUN` | `true` | Prevents writes to GitHub when true. |
| `ESTATE_AGENT_DB_PATH` | `.estate-agent/state.db` | SQLite state path. |
| `ESTATE_AGENT_GITHUB_TOKEN` | unset | GitHub API token for comments/issues. |
| `ESTATE_AGENT_GITHUB_WEBHOOK_SECRET` | unset | Secret used to verify GitHub webhook signatures. |
| `ESTATE_AGENT_DEFAULT_OWNER` | unset | Owner for issues created from log/error intake. |
| `ESTATE_AGENT_DEFAULT_REPO` | unset | Repo for issues created from log/error intake. |
| `ESTATE_AGENT_BUG_LABELS` | `bug,estate:bug,bug:autofix` | Labels treated as bug work. |
| `ESTATE_AGENT_FEATURE_LABELS` | `feature,feature request,estate:feature` | Labels treated as feature work. |

## Why This Exists

Bugloop was focused on autonomous bug fixing. Estate Agent is the broader version: an agent that owns a software estate, accepts requests through GitHub issues, listens to production signals, and runs a guarded SDLC loop.

Supermodel-style code graph context can become one tool in the loop, but the core project does not depend on any private repository or proprietary code.

## Guardrails

Estate Agent should never be trusted to mutate production by default.

The default operating posture is:

- dry-run first
- least-privilege GitHub token
- no auto-merge
- no secret exfiltration
- no CI/config/security-path edits without explicit approval
- human escalation on ambiguous, large, or risky changes

See [docs/architecture.md](docs/architecture.md) and [docs/mvp-plan.md](docs/mvp-plan.md).

