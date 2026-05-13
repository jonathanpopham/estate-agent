from __future__ import annotations

from dataclasses import asdict, dataclass
from enum import Enum
from typing import Any


class WorkKind(str, Enum):
    BUG = "bug"
    FEATURE = "feature"


class WorkState(str, Enum):
    OPEN = "open"
    PLANNED = "planned"
    ESCALATED = "escalated"


@dataclass(frozen=True)
class WorkItem:
    id: str
    kind: WorkKind
    source: str
    title: str
    body: str
    repo_owner: str | None = None
    repo_name: str | None = None
    issue_number: int | None = None
    labels: tuple[str, ...] = ()
    raw: dict[str, Any] | None = None

    @property
    def repo_slug(self) -> str | None:
        if not self.repo_owner or not self.repo_name:
            return None
        return f"{self.repo_owner}/{self.repo_name}"

    def to_json_dict(self) -> dict[str, Any]:
        data = asdict(self)
        data["kind"] = self.kind.value
        data["labels"] = list(self.labels)
        return data

