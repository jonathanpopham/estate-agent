# Evals

Estate Agent needs evals before it earns write access to repositories.

Initial deterministic eval classes:

- issue classification: bug vs feature vs ignore
- escalation detection: auth, billing, secrets, destructive data changes
- plan quality: actionable first step, verification path, rollback notes
- provider request shape: no missing model, no missing key, fallback behavior explicit

Metrics to track:

- actionable-plan rate
- false-accept rate
- escalation correctness
- cost per accepted work item
- token use per planning pass
- issue-to-PR success rate once builder mode exists

