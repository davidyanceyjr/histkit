# SESSION.md

## Current session

ID: `028-audit-log`

Status: completed

## Objective

Add the first audit-log slice for recording cleanup/apply run summaries.

## Scope

Implement:

- an `internal/audit` package for audit run records
- append-only audit-log writing to the configured state directory
- deterministic human-readable line rendering
- audit-log-focused tests

## Out of scope

- restore command wiring
- cleanup apply behavior
- audit-list CLI wiring

## Relevant skills

- `SKILLS/backup-restore.md`
- `SKILLS/testing.md`

## Acceptance criteria

- repository can validate and render an audit record for a cleanup/apply run summary
- repository can append audit records to a stable `audit.log` file under the state directory
- rendered audit entries are deterministic and include run id, timestamps, shell, rules, counts, backup id, and apply mode
- `go test ./...` passes

## Current repo state

The repository now has backup creation and atomic rewrite primitives, but it still lacks the audit-log package needed before destructive apply and restore workflows are treated as production-ready.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Audit format choices in this slice need to stay narrow so later apply and restore wiring can reuse them without locking the CLI too early.
- Log writing must stay append-only and deterministic because later audit review and tests will depend on stable output.

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

No non-blocking questions currently recorded.

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

No questions answered yet.

## End-of-session notes

Summary:

- Added `internal/audit` with a validated run-summary record model and deterministic line rendering for human-readable audit entries.
- Added append-only `audit.log` writing that creates the state directory as needed and writes private log files with `0600` permissions.
- Exposed the default `audit.log` path through `internal/config.Paths` and added focused audit/config tests.

Files changed:

- SESSION.md
- SESSIONS/028-audit-log.md
- internal/audit/log.go
- internal/audit/log_test.go
- internal/audit/model.go
- internal/audit/model_test.go
- internal/config/config.go
- internal/config/config_test.go

Files read:

- SESSION.md
- ROADMAP.md
- SKILLS/backup-restore.md
- SKILLS/testing.md
- README.md
- docs/histkit-implementation-plan.md
- internal/config/config.go
- internal/config/config_test.go
- internal/sanitize/model.go
- internal/sanitize/preview.go
- SESSIONS/023-dry-run-preview.md

Tests added:

- `TestRecordValidate`
- `TestRecordValidateRejectsInvalidInputs`
- `TestRenderLineIsDeterministic`
- `TestAppendCreatesLogFileAndDirectory`
- `TestAppendAppendsSecondRecord`
- `TestAppendRejectsInvalidRecord`
- `TestAppendRejectsEmptyPath`

Tests run:

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/audit ./internal/config`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- Keep this slice limited to audit-record modeling and append-only log writing.
- Use a human-readable deterministic line format for the first audit-log slice.
- Defer audit-list command behavior until later wiring exists.

Commands run:

- `git checkout -b 028-audit-log`
- `sed -n '1,260p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,240p' SKILLS/backup-restore.md`
- `sed -n '1,240p' SKILLS/testing.md`
- `rg -n "audit" README.md docs internal SESSIONS -S`
- `sed -n '709,730p' docs/histkit-implementation-plan.md`
- `sed -n '250,265p' docs/histkit-implementation-plan.md`
- `sed -n '430,450p' README.md`
- `sed -n '1,260p' internal/config/config.go`
- `sed -n '1,240p' internal/config/config_test.go`
- `sed -n '1,260p' internal/sanitize/model.go`
- `sed -n '1,260p' internal/sanitize/preview.go`
- `sed -n '1,220p' SESSIONS/023-dry-run-preview.md`
- `sed -n '140,160p' docs/histkit-implementation-plan.md`
- `rg -n "audit.log|audit list|run identifier|counts by action|apply mode" docs README.md SESSIONS -S`
- `ls internal`
- `gofmt -w internal/audit/log.go internal/audit/log_test.go internal/audit/model.go internal/audit/model_test.go internal/config/config.go internal/config/config_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/audit ./internal/config`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Assumptions made:

- A validated run-summary record plus append-only log writer is sufficient before CLI wiring exists.
- A deterministic key-value line format is adequately human-readable for this first audit-log slice.

Risks introduced or reduced:

- Reduced: later apply and restore flows now have a stable audit-log primitive and default path to write into.
- Ongoing: apply wiring, restore semantics, and audit listing remain pending later slices.

Next recommended session:

- `029-clean-apply`
