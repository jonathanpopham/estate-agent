from __future__ import annotations

import hmac
from hashlib import sha256


def verify_github_signature(secret: str | None, body: bytes, signature: str | None) -> bool:
    """Verify GitHub's X-Hub-Signature-256 header.

    A missing secret means verification is intentionally disabled. This supports
    local development, but production deployments should always configure one.
    """

    if not secret:
        return True
    if not signature or not signature.startswith("sha256="):
        return False

    expected = "sha256=" + hmac.new(secret.encode("utf-8"), body, sha256).hexdigest()
    return hmac.compare_digest(expected, signature)

