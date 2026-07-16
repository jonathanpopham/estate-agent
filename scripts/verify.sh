#!/usr/bin/env bash
# The one command that says whether this actually works. CI is a thin wrapper
# around it, and it is green on a fresh clone or the change does not ship. The
# service that stewards other people's estates holds itself to the discipline
# it enforces on them.
set -euo pipefail
cd "$(dirname "$0")/.."

step() { printf '\n== %s ==\n' "$1"; }
ok() { printf '  ok\n'; }

step "gofmt"
unformatted="$(gofmt -l . | { grep -v '^vendor/' || true; })"
if [ -n "$unformatted" ]; then
  echo "gofmt: these files are not formatted:" >&2
  echo "$unformatted" >&2
  exit 1
fi
ok

step "go vet"
go vet ./...
ok

step "go build"
go build ./...
ok

step "go test (race detector)"
go test -race ./...

printf '\nVERIFY OK: gofmt clean, vet clean, build clean, tests pass under the race detector\n'
