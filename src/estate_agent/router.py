from __future__ import annotations

import hashlib
import json
from typing import Any

from .config import Settings
from .models import WorkItem, WorkKind


def work_item_from_github_issue_event(
    event_name: str | None,
    payload: dict[str, Any],
    settings: Settings,
) -> WorkItem | None:
    if event_name != "issues":
        return None

    action = payload.get("action")
    if action not in {"opened", "edited", "labeled", "reopened"}:
        return None

    issue = payload.get("issue") or {}
    repo = payload.get("repository") or {}
    owner = (repo.get("owner") or {}).get("login")
    repo_name = repo.get("name")
    labels = tuple(_label_names(issue.get("labels") or []))
    kind = _kind_from_labels(labels, settings)

    if kind is None:
        return None

    issue_number = issue.get("number")
    return WorkItem(
        id=f"github:{owner}/{repo_name}#{issue_number}",
        kind=kind,
        source="github_issue",
        title=issue.get("title") or "Untitled issue",
        body=issue.get("body") or "",
        repo_owner=owner,
        repo_name=repo_name,
        issue_number=issue_number,
        labels=labels,
        raw=payload,
    )


def work_item_from_error_payload(
    payload: dict[str, Any],
    settings: Settings,
) -> WorkItem:
    service = str(payload.get("service") or payload.get("app") or "unknown-service")
    environment = str(payload.get("environment") or payload.get("env") or "unknown-env")
    severity = str(payload.get("severity") or "unknown")
    error = str(payload.get("error") or payload.get("message") or "Unknown error")
    stack = str(payload.get("stack") or payload.get("trace") or "")

    fingerprint_source = "\n".join([service, environment, error, stack])
    fingerprint = hashlib.sha256(fingerprint_source.encode("utf-8")).hexdigest()[:16]

    title = f"[{environment}] {service}: {error.splitlines()[0][:120]}"
    body = f"""## Production Error

**Service:** `{service}`
**Environment:** `{environment}`
**Severity:** `{severity}`
**Fingerprint:** `{fingerprint}`

### Error

```text
{error}
```

### Stack / Trace

```text
{stack or "No stack trace provided."}
```

### Raw Payload

```json
{json.dumps(payload, indent=2, sort_keys=True)}
```
"""

    return WorkItem(
        id=f"error:{fingerprint}",
        kind=WorkKind.BUG,
        source="error_ingest",
        title=title,
        body=body,
        repo_owner=settings.default_owner,
        repo_name=settings.default_repo,
        labels=("bug", "estate:bug", f"severity:{severity}"),
        raw=payload,
    )


def _label_names(labels: list[dict[str, Any]]) -> list[str]:
    names: list[str] = []
    for label in labels:
        name = label.get("name")
        if isinstance(name, str):
            names.append(name)
    return names


def _kind_from_labels(labels: tuple[str, ...], settings: Settings) -> WorkKind | None:
    normalized = {label.lower() for label in labels}
    bug_labels = {label.lower() for label in settings.bug_labels}
    feature_labels = {label.lower() for label in settings.feature_labels}

    if normalized & bug_labels:
        return WorkKind.BUG
    if normalized & feature_labels:
        return WorkKind.FEATURE
    return None
