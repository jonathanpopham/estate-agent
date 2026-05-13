from __future__ import annotations

from .models import WorkItem, WorkKind


class DryRunPlanner:
    """Produces a deterministic SDLC plan without modifying a repository."""

    def plan(self, work: WorkItem) -> str:
        if work.kind is WorkKind.BUG:
            decision = "Proceed if the error is reproducible or connected to a clear failing test."
            first_step = "Reproduce the failure from the issue context or production stack trace."
            tests = "Add or update a regression test that fails before the fix and passes after it."
        else:
            decision = "Proceed only after clarifying scope, acceptance criteria, and user value."
            first_step = "Convert the request into concrete acceptance criteria and non-goals."
            tests = "Add tests or checks that prove the requested behavior exists."

        repo = work.repo_slug or "unbound repository"
        labels = ", ".join(work.labels) if work.labels else "none"

        return f"""## Estate Agent Plan

**Mode:** dry run
**Work item:** `{work.id}`
**Kind:** `{work.kind.value}`
**Source:** `{work.source}`
**Repository:** `{repo}`
**Labels:** {labels}

### Decision

{decision}

### Proposed SDLC

1. **Triage:** Confirm severity, ownership, and whether the issue is actionable.
2. **Context gathering:** Read the issue, related code paths, recent commits, and relevant docs.
3. **Plan:** Produce a small implementation plan with files likely to change.
4. **Implement:** Make the smallest coherent change on a dedicated branch.
5. **Verify:** {tests}
6. **Review:** Check the diff for regressions, security risk, and product fit.
7. **PR:** Open a draft PR with summary, tests, and rollback notes.

### First Action

{first_step}

### Escalation Triggers

- missing reproduction or ambiguous acceptance criteria
- secrets, auth, billing, data deletion, or deployment configuration changes
- diff expected to exceed the configured size limit
- tests cannot be run in the agent environment

### Original Request

{work.body.strip() or "_No body provided._"}
"""

