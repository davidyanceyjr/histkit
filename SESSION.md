# SESSION.md

## Current session

ID: `026-backup-creation`

Status: completed

## Objective

Add concrete on-disk backup creation on top of the backup record model.

## Scope

Implement:

- package-level backup file creation helpers
- source-file checksum computation
- backup directory creation and file copy behavior
- backup-creation-focused tests

## Out of scope

- restore behavior
- audit logging
- cleanup apply behavior

## Relevant skills

- `SKILLS/backup-restore.md`
- `SKILLS/testing.md`

## Acceptance criteria

- repository can create a backup file from a source history file
- created backup records include a computed checksum and on-disk backup path
- backup creation verifies copied bytes by checksum before returning
- `go test ./...` passes

## Current repo state

The repository has a validated backup record model in `internal/backup`, but it still needs a concrete creation path that copies source files into the derived backup location and computes checksums internally.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Backup creation touches the filesystem, so tests must stay isolated to temporary files and directories.
- Later restore slices will depend on the backup path and checksum contract established here.

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

- Added `internal/backup/create.go` with package-level backup creation and checksum helpers.
- Added checksum-verified file-copy behavior that creates the derived backup directory and refuses to overwrite an existing backup artifact.
- Added focused temp-file tests for successful creation, directory creation, checksum generation, and invalid-input handling.

Files changed:

- internal/backup/create.go
- internal/backup/create_test.go
- SESSION.md
- SESSIONS/026-backup-creation.md

Files read:

- SESSION.md
- ROADMAP.md
- SKILLS/backup-restore.md
- docs/histkit-implementation-plan.md
- internal/backup/model.go
- internal/backup/model_test.go
- internal/index/writer.go

Tests added:

- `TestCreate`
- `TestCreateCreatesBackupDirectory`
- `TestCreateRejectsInvalidInputs`
- `TestCreateRejectsExistingBackupPath`
- `TestChecksumFile`
- `TestChecksumFileRejectsEmptyPath`

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- Keep this slice limited to backup file creation and checksum verification.
- Reuse the existing backup record model and derived backup-path format.
- Prefix file checksums as `sha256:` to match the model contract.

Commands run:

- `git checkout -b 026-backup-creation`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/backup-restore.md`
- `sed -n '1,220p' internal/backup/model.go`
- `sed -n '1,220p' internal/backup/model_test.go`
- `rg -n "backup" ROADMAP.md SESSION.md docs internal -g '!**/*.sum'`
- `sed -n '680,740p' docs/histkit-implementation-plan.md`
- `sed -n '840,875p' docs/histkit-implementation-plan.md`
- `rg -n "sha256|checksum" internal docs README.md`
- `sed -n '1,180p' internal/index/writer.go`
- `gofmt -w internal/backup/create.go internal/backup/create_test.go internal/backup/model.go internal/backup/model_test.go`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- Backup creation should fail rather than overwrite an existing derived backup path.
- Computing the checksum from the source file before copy and revalidating the copied file is sufficient for this slice.

Risks introduced or reduced:

- Reduced: later apply and restore work will have a concrete, checksum-verified backup artifact to depend on.
- Ongoing: restore, audit logging, and atomic history rewrite are still pending later slices.

Next recommended session:

- `027-atomic-rewrite`
