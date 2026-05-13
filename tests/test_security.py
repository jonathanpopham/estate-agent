from __future__ import annotations

import hmac
import unittest
from hashlib import sha256

from estate_agent.security import verify_github_signature


class SignatureTests(unittest.TestCase):
    def test_signature_verifies(self) -> None:
        secret = "top-secret"
        body = b'{"hello":"world"}'
        digest = hmac.new(secret.encode("utf-8"), body, sha256).hexdigest()

        self.assertTrue(verify_github_signature(secret, body, f"sha256={digest}"))

    def test_bad_signature_fails(self) -> None:
        self.assertFalse(verify_github_signature("secret", b"body", "sha256=bad"))

    def test_missing_secret_disables_verification_for_local_dev(self) -> None:
        self.assertTrue(verify_github_signature(None, b"body", None))


if __name__ == "__main__":
    unittest.main()

