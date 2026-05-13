# SESSION.md

## Current session

ID: `056-ci-gosec-scan`

Status: in progress

## Objective

Add a low-noise `gosec` stage to histkit's GitHub Actions workflow using a Go toolchain install.

## Scope

Implement:

- install `gosec` in GitHub Actions via `go install`
- run `gosec` in CI against histkit's own packages
- keep the scan useful for a solo CLI project by avoiding repo-local cache noise and fixing any small high-signal findings uncovered by the chosen scan profile

## Out of scope

- adding SARIF or external security reporting integrations
- broad production hardening for every current `G204` or `G304` path/exec finding
- introducing a workflow matrix or separate security workflow

## Relevant skills

- `SKILLS/testing.md`
- `SKILLS/go-cli.md`

## Acceptance criteria

- CI installs `gosec` through the Go toolchain
- CI runs a dedicated Go SAST stage
- the configured scan focuses on histkit packages instead of repo-local cache directories
- local verification with the same `gosec` flags passes
- `go test ./...` and `go build ./cmd/histkit` pass

## Current repo state

Branch `056-ci-gosec-scan` contains the workflow update plus two permission-tightening fixes needed to keep the new scan actionable. No PR is open yet.

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

## Risks to watch

- repo-local cache directories can create meaningless SAST output if the workflow scans `./...` naively
- the temporary `G204` and `G304` exclusions could hide future regressions if they are left unreviewed for too long
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

- intent: add a low-noise `gosec` CI stage without restructuring the workflow
- scope: `.github/workflows/ci.yml`, any small production fixes required to keep the chosen scan profile actionable, `DECISIONS.md`, `RISKS.md`, `SESSION.md`, and the final session note
- constraints: install `gosec` via `go install`, keep the scan focused on histkit packages, avoid broad suppressions, leave the repository buildable, and use writable temp Go cache paths for local verification in this shell environment
- files read:
  - `SESSION.md`: previous session state and carry-forward structure
  - `ROADMAP.md`: roadmap boundary confirmation for the slice
  - `SKILLS/testing.md`: verification expectations
  - `SKILLS/go-cli.md`: implementation constraints relevant to the CLI repository
  - `.github/workflows/ci.yml`: existing CI job shape and insertion point for the new stage
  - `SESSIONS/049-basic-ci-workflow.md`: previous CI workflow conventions
  - `DECISIONS.md`: current durable project decisions
  - `RISKS.md`: current risk register
  - `internal/index/schema.go`: existing SQLite parent-directory permissions
  - `internal/snippets/store.go`: existing snippet store parent-directory permissions
  - `internal/picker/fzf.go`: context for the current `G204` finding
  - `internal/history/detect.go`: context for current path-based local filesystem behavior
- files changed:
  - `.github/workflows/ci.yml`: added `gosec` installation and a scoped SAST stage
  - `internal/index/schema.go`: tightened SQLite parent-directory permissions to `0o700`
  - `internal/snippets/store.go`: tightened snippet store parent-directory permissions to `0o700`
  - `DECISIONS.md`: added the CI `gosec` policy decision
  - `RISKS.md`: added the static-analysis exclusion risk entry
  - `SESSION.md`: replaced the previous session carry-forward with this session's working state
  - `SESSIONS/056-ci-gosec-scan.md`: recorded the completed session
- commands run:
  - `git status --short --branch`: inspected repository state before branching
  - `ls -1 SESSIONS | sort | tail -n 12`: identified the next session number
  - `sed -n '1,240p' .github/workflows/ci.yml`: reviewed the current CI workflow
  - `sed -n '1,220p' SKILLS/go-cli.md`: reviewed relevant implementation constraints
  - `sed -n '1,240p' SESSION.md`: reviewed the previous session state
  - `sed -n '1,220p' ROADMAP.md`: confirmed roadmap boundaries
  - `git checkout -b 056-ci-gosec-scan`: created the session branch
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath GOBIN=/tmp/histkit-gobin go install github.com/securego/gosec/v2/cmd/gosec@latest`: installed `gosec` locally with writable temp Go paths
  - `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath gosec ./...`: showed that naive `./...` scanning was too noisy because repo-local caches were treated as source
  - `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath gosec -tests=false ./cmd/... ./internal/...`: reduced the scan to histkit packages and surfaced 12 project-local findings
  - `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath gosec -tests=false -exclude=G204,G304 ./cmd/... ./internal/...`: isolated the remaining two `G301` permission findings
  - `sed -n '1,200p' internal/index/schema.go`: reviewed SQLite parent-directory permissions
  - `sed -n '1,200p' internal/snippets/store.go`: reviewed snippet store parent-directory permissions
  - `sed -n '1,220p' internal/picker/fzf.go`: reviewed the current `fzf` subprocess launch flow
  - `sed -n '1,220p' internal/history/detect.go`: reviewed current history source path handling
  - `gofmt -w internal/index/schema.go internal/snippets/store.go`: formatted the changed Go files
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`: passed
  - `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath /tmp/histkit-gobin/gosec -tests=false -exclude=G204,G304 ./cmd/... ./internal/...`: passed
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build ./cmd/histkit`: passed
  - `rm -f histkit`: removed the local build artifact after verification
- tests:
  - added: none
  - changed: none
  - run:
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
    - `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath /tmp/histkit-gobin/gosec -tests=false -exclude=G204,G304 ./cmd/... ./internal/...`
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build ./cmd/histkit`
  - skipped: none
  - failing: none
- decisions:
  - keep the CI change in the existing single Ubuntu workflow
  - scope `gosec` to `./cmd/...` and `./internal/...` so repo-local cache directories do not become CI noise
  - disable test scanning for this stage because the deliverable is a production-focused SAST check
  - exclude `G204` and `G304` for now and document that policy explicitly rather than leaving unexplained CI failures
  - fix the two `G301` permission findings in code instead of suppressing them
- assumptions:
  - `NON-BLOCKING`: excluding `G204` and `G304` is safe for this slice because the findings arise from known local `fzf` execution and user-local path access patterns that already sit inside histkit's current CLI contract
- unresolved questions:
  - none currently recorded
- next step: push branch `056-ci-gosec-scan`, open a PR, and get human review on the new CI scan policy

## End-of-session notes

Summary:

- Added a `gosec` stage to CI using a Go toolchain install and a scan profile scoped to histkit packages.
- Tightened SQLite and snippet store parent-directory permissions so the scan can stay useful without suppressing those findings.
- Recorded the temporary `G204` and `G304` policy and its residual risk explicitly.

Tests run:

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
- `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath /tmp/histkit-gobin/gosec -tests=false -exclude=G204,G304 ./cmd/... ./internal/...`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build ./cmd/histkit`

Known failures:

- None after switching local Go cache paths to writable temp directories.

Next recommended session:

- Add `govulncheck` to the workflow so dependency vulnerability reachability is covered alongside `gosec`.
