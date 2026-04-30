# SESSION.md

## Current session

ID: `025-backup-model`

Status: completed

## Objective

Add the initial backup record model and metadata helpers.

## Scope

Implement:

- package-level backup record model
- backup record validation
- metadata-only backup record creation helpers
- backup-model-focused tests

## Out of scope

- backup file creation
- checksum computation
- restore behavior
- audit logging
- cleanup apply behavior

## Relevant skills

- `SKILLS/backup-restore.md`
- `SKILLS/testing.md`

## Acceptance criteria

- repository contains a backup metadata model matching the documented fields
- backup records include source file, backup path, created time, and checksum
- deterministic backup identifiers exist for later creation/restore work
- `go test ./...` passes

## Current repo state

The repository now has a package-level backup record model in `internal/backup` plus metadata-only helpers to derive backup IDs and backup paths for later backup creation work.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Backup records are still metadata-only and do not create or verify files yet.
- Checksum values are required by the model but still supplied externally until the creation slice lands.

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

- Added `internal/backup/model.go` with a backup record model and validation.
- Added deterministic backup ID generation and metadata-only record creation helpers.
- Added tests for validation, ID generation, derived backup paths, and invalid input handling.

Files changed:

- internal/backup/model.go
- internal/backup/model_test.go
- SESSION.md
- SESSIONS/025-backup-model.md

Files read:

- SESSION.md
- ROADMAP.md
- SKILLS/backup-restore.md
- docs/histkit-implementation-plan.md
- README.md
- internal/sanitize/quarantine.go
- internal/sanitize/quarantine_test.go

Tests added:

- `TestRecordValidate`
- `TestRecordValidateRequiresFields`
- `TestBackupID`
- `TestBuildRecord`
- `TestBuildRecordRejectsInvalidInputs`

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Keep the first backup slice limited to metadata and validation only.
- Derive backup paths as `<backup-dir>/<backup-id>/<basename(source-file)>`.
- Use deterministic timestamp-plus-sequence backup IDs for now.

Commands run:

- `git checkout -b 025-backup-model`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/backup-restore.md`
- `rg -n "backup|restore|checksum|backup_path|source_file|created_at|backup ID|audit" docs/histkit-implementation-plan.md README.md -S`
- `sed -n '456,470p' docs/histkit-implementation-plan.md`
- `sed -n '687,705p' docs/histkit-implementation-plan.md`
- `sed -n '155,163p' README.md`
- `rg --files internal | rg 'sanitize|backup|audit|restore'`
- `sed -n '1,260p' internal/sanitize/quarantine.go`
- `sed -n '1,260p' internal/sanitize/quarantine_test.go`
- `gofmt -w internal/backup/model.go internal/backup/model_test.go`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- Deterministic timestamp-plus-sequence identifiers are sufficient for the first backup slice.
- Metadata-only path derivation is acceptable before actual file-copy work exists.

Risks introduced or reduced:

- Reduced: later backup creation and restore slices now have a concrete validated metadata contract to build on.
- Ongoing: real backup durability and checksum generation still depend on later slices.

Next recommended session:

- `026-backup-creation`
