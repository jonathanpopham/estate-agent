# MVP Plan

## Milestone 0: Public Scaffold

- [x] public repo
- [x] dry-run server
- [x] GitHub issue event normalization
- [x] generic error payload normalization
- [x] SQLite work-item state
- [x] deterministic SDLC plan generation

## Milestone 1: GitHub Commenter

- [ ] create GitHub App or fine-grained PAT setup guide
- [ ] post dry-run plans back to labeled issues
- [ ] create issues from error payloads
- [ ] add idempotency so duplicate errors update/comment instead of spamming
- [ ] add webhook replay tests

## Milestone 2: Agent Sandbox

- [ ] clone target repo into isolated workspace
- [ ] create branch per work item
- [ ] run configurable agent command
- [ ] enforce timeout and max-iteration limits
- [ ] collect changed files and test output

## Milestone 3: Draft PR Builder

- [ ] push branch
- [ ] open draft PR
- [ ] include test output, risk notes, rollback notes
- [ ] fail closed on large diffs or risky paths

## Milestone 4: Cloud Deployment

- [ ] AWS Terraform target
- [ ] Azure Terraform target
- [ ] secret management
- [ ] queue-backed worker
- [ ] structured logs and metrics

## Candidate Names

The repo is named `estate-agent` because the concept is a software estate steward. Other possible product names:

- Estate Agent
- Code Steward
- App Steward
- Maintainer
- Groundskeeper
- Software Estate Manager

