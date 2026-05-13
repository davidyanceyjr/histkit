# SESSION.md

## Current session

ID: `058-ci-binary-smoke-checks`

Status: in progress

## Objective

Add basic executable smoke checks to histkit's GitHub Actions workflow after building the CLI binary.

## Scope

Implement:

- build the CLI binary as a named artifact in CI
- run `./histkit --help` in CI
- run `./histkit doctor --help` in CI

## Out of scope

- adding SARIF or other external security reporting integrations
- any production code changes unless the smoke check uncovers a defect
- introducing a workflow matrix or separate security workflow

## Relevant skills

- `SKILLS/testing.md`
- `SKILLS/go-cli.md`

## Acceptance criteria

- CI builds a named CLI binary artifact
- CI runs `./histkit --help` successfully
- CI runs `./histkit doctor --help` successfully
- `go test ./...` and `go build ./cmd/histkit` pass

## Current repo state

Branch `058-ci-binary-smoke-checks` contains the workflow update for binary smoke checks only. No PR is open yet.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target
- Default automation runs `scan`, not destructive apply
- Wrapper logic stays outside the Go binary under `contrib/`
- README-promised `--config` support should fail early and consistently across `scan`, `clean`, and `restore`
- bare `clean` and `clean --dry-run` are the same planning mode and should stay equivalent
- `--shell` filtering during `clean --apply` must restrict mutation, backup creation, and audit logging to the selected shell source
- shell-filter follow-up coverage should stay at the command layer because the contract spans detection, rewrite, backup, and audit together
- CI `gosec` coverage should stay scoped to histkit packages and avoid turning repo-local caches or known path-contract findings into recurring noise
- dependency vulnerability checks should run in CI with `govulncheck` against the full module graph
- CI should smoke-test the built executable rather than just compile it

## Risks to watch

- smoke tests should stay limited to help output so they do not require mutable history fixtures or a real shell session
- command-level assertions intentionally depend on current user-visible output fragments; any future wording changes should be updated deliberately in tests

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

No non-blocking questions currently recorded.

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

No answered questions were recorded during this session.

## Working state

- intent: add binary smoke checks to the existing CI workflow without restructuring it
- scope: `.github/workflows/ci.yml`, `SESSION.md`, and the final session note
- constraints: keep the smoke stage to help output only, leave the repository buildable, and use writable temp Go cache paths for local verification in this shell environment
- files read:
  - `SESSION.md`: previous closed session state and carry-forward structure
  - `ROADMAP.md`: roadmap boundary confirmation for the slice
  - `SKILLS/testing.md`: verification expectations
  - `SKILLS/go-cli.md`: implementation constraints relevant to the CLI repository
  - `.github/workflows/ci.yml`: existing CI job shape and insertion point for the smoke stage
- files changed:
  - `.github/workflows/ci.yml`: updated build output and added binary help smoke checks
  - `SESSION.md`: replaced the previous closed session carry-forward with this session's working state
  - `SESSIONS/058-ci-binary-smoke-checks.md`: records this session
- commands run:
  - `git status --short --branch`: inspected repository state on `main`
  - `sed -n '1,240p' SESSION.md`: reviewed the previous closed session state
  - `sed -n '1,220p' ROADMAP.md`: confirmed roadmap boundaries
  - `sed -n '1,220p' SKILLS/testing.md`: reviewed verification expectations
  - `sed -n '1,220p' .github/workflows/ci.yml`: reviewed the current CI workflow
  - `ls -1 SESSIONS | sort | tail -n 4`: identified the next session number
  - `git checkout -b 058-ci-binary-smoke-checks`: created the session branch
  - `sed -n '1,220p' SKILLS/go-cli.md`: reviewed relevant implementation constraints
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`: passed
  - `mkdir -p /tmp/histkit-smoke && env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build -o /tmp/histkit-smoke/histkit ./cmd/histkit && /tmp/histkit-smoke/histkit --help && /tmp/histkit-smoke/histkit doctor --help`: passed
- tests:
  - added: none
  - changed: none
  - run:
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
    - `mkdir -p /tmp/histkit-smoke && env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build -o /tmp/histkit-smoke/histkit ./cmd/histkit && /tmp/histkit-smoke/histkit --help && /tmp/histkit-smoke/histkit doctor --help`
  - skipped: none
  - failing: none
- decisions:
  - keep the smoke-check change in the existing single Ubuntu workflow
  - build a named binary and smoke-test the same artifact in CI
- assumptions:
  - none currently recorded
- unresolved questions:
  - none currently recorded
- next step: open a PR for the workflow change and wait for review

## End-of-session notes

Summary:

- Adding built-binary smoke checks to CI after the CLI build step.
- Smoke testing `./histkit --help` and `./histkit doctor --help` against a named artifact to validate the executable rather than just compilation.

Tests run:

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
- `mkdir -p /tmp/histkit-smoke && env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build -o /tmp/histkit-smoke/histkit ./cmd/histkit && /tmp/histkit-smoke/histkit --help && /tmp/histkit-smoke/histkit doctor --help`

Known failures:

- None currently recorded.

Next recommended session:

- Add a small CI step or job if the project later wants to split smoke checks from the main build.
