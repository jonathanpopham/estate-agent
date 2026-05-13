from __future__ import annotations

import tempfile
import unittest
from pathlib import Path

from estate_agent.config import Settings
from estate_agent.models import WorkKind
from estate_agent.router import work_item_from_error_payload, work_item_from_github_issue_event


def settings() -> Settings:
    return Settings(
        dry_run=True,
        db_path=Path(tempfile.gettempdir()) / "estate-agent-test.db",
        github_token=None,
        github_webhook_secret=None,
        default_owner="jonathanpopham",
        default_repo="estate-agent",
        bug_labels=("bug", "estate:bug"),
        feature_labels=("feature", "feature request", "estate:feature"),
    )


class RouterTests(unittest.TestCase):
    def test_github_issue_with_feature_label_becomes_feature_work(self) -> None:
        payload = {
            "action": "opened",
            "repository": {"name": "demo", "owner": {"login": "octo"}},
            "issue": {
                "number": 42,
                "title": "Add dark mode",
                "body": "Users asked for it.",
                "labels": [{"name": "feature request"}],
            },
        }

        work = work_item_from_github_issue_event("issues", payload, settings())

        self.assertIsNotNone(work)
        assert work is not None
        self.assertEqual(work.kind, WorkKind.FEATURE)
        self.assertEqual(work.id, "github:octo/demo#42")

    def test_unlabeled_issue_is_ignored(self) -> None:
        payload = {
            "action": "opened",
            "repository": {"name": "demo", "owner": {"login": "octo"}},
            "issue": {"number": 42, "title": "Question", "labels": []},
        }

        self.assertIsNone(work_item_from_github_issue_event("issues", payload, settings()))

    def test_error_payload_becomes_stable_bug_work(self) -> None:
        payload = {
            "service": "checkout",
            "environment": "prod",
            "error": "TypeError: x is undefined",
            "stack": "at checkout.ts:10",
            "severity": "high",
        }

        first = work_item_from_error_payload(payload, settings())
        second = work_item_from_error_payload(payload, settings())

        self.assertEqual(first.kind, WorkKind.BUG)
        self.assertEqual(first.id, second.id)
        self.assertEqual(first.repo_slug, "jonathanpopham/estate-agent")


if __name__ == "__main__":
    unittest.main()

