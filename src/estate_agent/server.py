from __future__ import annotations

import json
from typing import Any

try:
    from fastapi import FastAPI, HTTPException, Request
except ImportError as exc:  # pragma: no cover
    raise RuntimeError("Install server dependencies with: pip install -e '.[server]'") from exc

from .agent import DryRunPlanner
from .config import Settings
from .github_client import GitHubClient
from .router import work_item_from_error_payload, work_item_from_github_issue_event
from .security import verify_github_signature
from .state import StateStore

settings = Settings.from_env()
store = StateStore(settings.db_path)
planner = DryRunPlanner()
app = FastAPI(title="Estate Agent", version="0.1.0")


@app.get("/health")
def health() -> dict[str, Any]:
    return {
        "ok": True,
        "dry_run": settings.dry_run,
        "db_path": str(settings.db_path),
    }


@app.post("/webhooks/github")
async def github_webhook(request: Request) -> dict[str, Any]:
    body = await request.body()
    signature = request.headers.get("X-Hub-Signature-256")
    if not verify_github_signature(settings.github_webhook_secret, body, signature):
        raise HTTPException(status_code=401, detail="Invalid GitHub signature")

    event_name = request.headers.get("X-GitHub-Event")
    payload = json.loads(body.decode("utf-8"))
    work = work_item_from_github_issue_event(event_name, payload, settings)

    if work is None:
        return {"ignored": True, "reason": "event is not actionable"}

    store.upsert(work)
    plan = planner.plan(work)
    comment_posted = False

    if not settings.dry_run and settings.github_token and work.repo_owner and work.repo_name and work.issue_number:
        GitHubClient(settings.github_token).create_issue_comment(
            work.repo_owner,
            work.repo_name,
            work.issue_number,
            plan,
        )
        comment_posted = True

    return {
        "ignored": False,
        "work_item": work.to_json_dict(),
        "comment_posted": comment_posted,
        "plan": plan,
    }


@app.post("/ingest/error")
async def ingest_error(request: Request) -> dict[str, Any]:
    payload = await request.json()
    work = work_item_from_error_payload(payload, settings)
    plan = planner.plan(work)
    store.upsert(work)

    issue_created = False
    issue_url = None

    if not settings.dry_run and settings.github_token and settings.default_owner and settings.default_repo:
        issue = GitHubClient(settings.github_token).create_issue(
            settings.default_owner,
            settings.default_repo,
            work.title,
            work.body + "\n\n" + plan,
            labels=list(work.labels),
        )
        issue_created = True
        issue_url = issue.get("html_url")

    return {
        "work_item": work.to_json_dict(),
        "issue_created": issue_created,
        "issue_url": issue_url,
        "plan": plan,
    }

