# SESSION.md

## Current session

ID: `030-restore-command`

Status: completed

## Objective

Add the first `histkit restore` slice for listing backup records and restoring a selected backup safely.

## Scope

Implement:

- backup-record persistence alongside created backup files
- a restore path that loads backup metadata, validates backup integrity, and rewrites the original history file atomically
- `restore` CLI wiring for backup listing and restore-by-id
- focused tests for metadata persistence and restore behavior

## Out of scope

- audit-list CLI wiring
- backup pruning or retention policies
- failure-recovery matrix testing across partial apply or restore failures

## Relevant skills

- `SKILLS/backup-restore.md`
- `SKILLS/testing.md`

## Acceptance criteria

- `histkit restore` lists available backups when no backup ID is given
- `histkit restore <backup-id>` validates backup integrity and restores the original history file atomically
- successful restore operations append an audit record
- `go test ./...` passes

## Current repo state

The repository now has backup creation, atomic rewrite, audit logging, cleanup apply, and an initial restore command backed by persisted backup metadata files.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Older backups created before metadata persistence will not be restorable through this first metadata-driven restore flow.
- Later failure-recovery testing still needs to exercise interrupted apply and restore scenarios.

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

No non-blocking questions currently recorded.

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

No questions answered this session.

## End-of-session notes

Summary:

- Added persisted `record.toml` metadata for every new backup so restore can recover source path, checksum, and creation time safely.
- Added backup listing and restore helpers that load metadata, verify checksum integrity, and rewrite the original history file atomically.
- Added the `histkit restore` command to list available backups or restore one by ID, with audit-log entries for successful restores.

Files changed:

- SESSION.md
- SESSIONS/030-restore-command.md
- internal/backup/create.go
- internal/backup/store.go
- internal/backup/store_test.go
- internal/cli/restore.go
- internal/cli/restore_test.go
- internal/cli/root.go
- internal/cli/root_test.go

Files read:

- SESSION.md
- ROADMAP.md
- README.md
- docs/histkit-implementation-plan.md
- internal/backup/create.go
- internal/backup/create_test.go
- internal/backup/atomic_test.go
- internal/cli/root_test.go
- internal/audit/model_test.go
- internal/audit/log_test.go
- internal/cli/clean_test.go

Tests added:

- `TestWriteRecordAndLoadRecord`
- `TestListRecordsSortsNewestFirst`
- `TestRestoreRewritesTargetFromBackup`
- `TestRestoreRejectsChecksumMismatch`
- `TestExecuteRestoreListsBackupsWhenNoIDProvided`
- `TestExecuteRestoreRestoresSpecificBackupAndAudits`
- `TestExecuteRestoreNoBackups`

Tests run:

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- Persist backup metadata as `record.toml` under each backup ID directory.
- Use metadata-driven restore instead of inferring source paths from backup filenames alone.
- Log restore operations through the existing audit log using a narrow `restore` rule marker and `apply=false`.

Commands run:

- `git status --short --branch`
- `sed -n '1,240p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `rg -n "restore|backup|audit|clean apply|RewriteAtomic|BackupID|backups" internal cmd README.md docs -g '!**/*_test.go'`
- `git checkout -b 030-restore-command`
- `sed -n '150,175p' README.md`
- `sed -n '236,245p' docs/histkit-implementation-plan.md`
- `sed -n '689,706p' docs/histkit-implementation-plan.md`
- `sed -n '1,220p' internal/backup/create_test.go`
- `sed -n '1,220p' internal/backup/atomic_test.go`
- `sed -n '1,240p' internal/cli/root_test.go`
- `sed -n '1,240p' internal/audit/model_test.go`
- `sed -n '1,220p' internal/audit/log_test.go`
- `sed -n '1,220p' internal/cli/clean_test.go`
- `gofmt -w internal/backup/store.go internal/backup/store_test.go internal/backup/create.go internal/cli/restore.go internal/cli/restore_test.go internal/cli/root.go internal/cli/root_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup ./internal/cli`
- `gofmt -w internal/backup/store.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup ./internal/cli`
- `gofmt -w internal/cli/restore.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`
- `git status --short`
- `date -u +%Y-%m-%d`

Assumptions made:

- Backup metadata persisted from this slice forward is the supported restore source of truth; older metadata-less backups are outside this slice.
- A restore audit entry can reuse the existing audit-record shape with `RuleNames=["restore"]` and zero counts.

Risks introduced or reduced:

- Reduced: restore now validates backup checksum before replacing the target file.
- Reduced: restore now uses persisted source-path metadata instead of guessing from backup filenames.
- Ongoing: legacy backups without metadata are not yet surfaced through the restore command.

Next recommended session:

- `031-failure-recovery-tests`
