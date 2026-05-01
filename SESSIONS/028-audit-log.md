# Session 028: Audit Log

## Objective

Add the first human-readable audit-log slice for cleanup/apply run summaries.

## Completed

- Added `internal/audit/model.go` with a validated audit run-summary record and deterministic line rendering.
- Added `internal/audit/log.go` with append-only `audit.log` writing that creates parent directories and syncs writes.
- Added config path support for the default `audit.log` location.
- Added focused tests for validation, deterministic rendering, log creation, and append behavior.

## Files changed

- SESSION.md
- SESSIONS/028-audit-log.md
- internal/audit/log.go
- internal/audit/log_test.go
- internal/audit/model.go
- internal/audit/model_test.go
- internal/config/config.go
- internal/config/config_test.go

## Files read

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

## Tests added

- `TestRecordValidate`
- `TestRecordValidateRejectsInvalidInputs`
- `TestRenderLineIsDeterministic`
- `TestAppendCreatesLogFileAndDirectory`
- `TestAppendAppendsSecondRecord`
- `TestAppendRejectsInvalidRecord`
- `TestAppendRejectsEmptyPath`

## Tests run

```bash
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/audit ./internal/config
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 028-audit-log
sed -n '1,260p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,240p' SKILLS/backup-restore.md
sed -n '1,240p' SKILLS/testing.md
rg -n "audit" README.md docs internal SESSIONS -S
sed -n '709,730p' docs/histkit-implementation-plan.md
sed -n '250,265p' docs/histkit-implementation-plan.md
sed -n '430,450p' README.md
sed -n '1,260p' internal/config/config.go
sed -n '1,240p' internal/config/config_test.go
sed -n '1,260p' internal/sanitize/model.go
sed -n '1,260p' internal/sanitize/preview.go
sed -n '1,220p' SESSIONS/023-dry-run-preview.md
sed -n '140,160p' docs/histkit-implementation-plan.md
rg -n "audit.log|audit list|run identifier|counts by action|apply mode" docs README.md SESSIONS -S
ls internal
gofmt -w internal/audit/log.go internal/audit/log_test.go internal/audit/model.go internal/audit/model_test.go internal/config/config.go internal/config/config_test.go
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/audit ./internal/config
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...
```

## Decisions

- Keep this slice limited to audit-record modeling and append-only log writing.
- Use a deterministic key-value line format for the first audit log so later CLI listing can parse it without requiring a schema migration.
- Keep the first audit file private and append-only under the state directory.

## Assumptions

- A validated run-summary record plus append-only log writer is sufficient before CLI wiring exists.
- A deterministic key-value line format is adequately human-readable for this first audit-log slice.

## Known issues

- No `audit list` command exists yet.
- Cleanup apply and restore flows do not write audit records yet.

## Risks reduced

- The repo now has a reusable audit-log primitive and a stable default path for later destructive workflows.

## Next recommended session

`029-clean-apply`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
