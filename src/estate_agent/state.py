from __future__ import annotations

import json
import sqlite3
from pathlib import Path

from .models import WorkItem


class StateStore:
    def __init__(self, db_path: Path) -> None:
        self.db_path = db_path
        self.db_path.parent.mkdir(parents=True, exist_ok=True)
        self._init()

    def upsert(self, work: WorkItem) -> None:
        payload = json.dumps(work.to_json_dict(), sort_keys=True)
        with self._connect() as conn:
            conn.execute(
                """
                insert into work_items (id, kind, source, repo_slug, issue_number, title, payload)
                values (?, ?, ?, ?, ?, ?, ?)
                on conflict(id) do update set
                  kind=excluded.kind,
                  source=excluded.source,
                  repo_slug=excluded.repo_slug,
                  issue_number=excluded.issue_number,
                  title=excluded.title,
                  payload=excluded.payload,
                  updated_at=datetime('now')
                """,
                (
                    work.id,
                    work.kind.value,
                    work.source,
                    work.repo_slug,
                    work.issue_number,
                    work.title,
                    payload,
                ),
            )

    def get(self, work_id: str) -> dict | None:
        with self._connect() as conn:
            row = conn.execute(
                "select payload from work_items where id = ?",
                (work_id,),
            ).fetchone()
        if row is None:
            return None
        return json.loads(row[0])

    def _connect(self) -> sqlite3.Connection:
        return sqlite3.connect(self.db_path)

    def _init(self) -> None:
        with self._connect() as conn:
            conn.execute(
                """
                create table if not exists work_items (
                  id text primary key,
                  kind text not null,
                  source text not null,
                  repo_slug text,
                  issue_number integer,
                  title text not null,
                  payload text not null,
                  created_at text not null default (datetime('now')),
                  updated_at text not null default (datetime('now'))
                )
                """
            )

