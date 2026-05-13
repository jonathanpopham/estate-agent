from __future__ import annotations

import os
from dataclasses import dataclass
from pathlib import Path


def _bool_from_env(value: str | None, default: bool) -> bool:
    if value is None:
        return default
    return value.strip().lower() in {"1", "true", "yes", "on"}


def _csv(value: str | None, default: tuple[str, ...]) -> tuple[str, ...]:
    if not value:
        return default
    return tuple(part.strip() for part in value.split(",") if part.strip())


@dataclass(frozen=True)
class Settings:
    dry_run: bool
    db_path: Path
    github_token: str | None
    github_webhook_secret: str | None
    default_owner: str | None
    default_repo: str | None
    bug_labels: tuple[str, ...]
    feature_labels: tuple[str, ...]

    @classmethod
    def from_env(cls) -> "Settings":
        return cls(
            dry_run=_bool_from_env(os.getenv("ESTATE_AGENT_DRY_RUN"), True),
            db_path=Path(os.getenv("ESTATE_AGENT_DB_PATH", ".estate-agent/state.db")),
            github_token=os.getenv("ESTATE_AGENT_GITHUB_TOKEN") or None,
            github_webhook_secret=os.getenv("ESTATE_AGENT_GITHUB_WEBHOOK_SECRET") or None,
            default_owner=os.getenv("ESTATE_AGENT_DEFAULT_OWNER") or None,
            default_repo=os.getenv("ESTATE_AGENT_DEFAULT_REPO") or None,
            bug_labels=_csv(
                os.getenv("ESTATE_AGENT_BUG_LABELS"),
                ("bug", "estate:bug", "bug:autofix"),
            ),
            feature_labels=_csv(
                os.getenv("ESTATE_AGENT_FEATURE_LABELS"),
                ("feature", "feature request", "estate:feature"),
            ),
        )

