# SESSION.md

## Current session

ID: `027-atomic-rewrite`

Status: completed

## Objective

Add atomic file-rewrite support for later restore and cleanup-apply slices.

## Scope

Implement:

- package-level atomic rewrite helper(s)
- same-directory temp-file replacement behavior
- permission-preserving rewrite behavior
- atomic-rewrite-focused tests

## Out of scope

- backup creation changes
- restore command wiring
- audit logging
- cleanup apply behavior

## Relevant skills

- `SKILLS/backup-restore.md`
- `SKILLS/testing.md`

## Acceptance criteria

- repository can atomically replace a target file with rewritten contents
- rewrites preserve the target file's permissions when replacing an existing regular file
- rewrites use temp files in the target directory and leave no temp artifacts on success
- `go test ./...` passes

## Current repo state

The repository now has validated backup metadata and concrete backup-file creation in `internal/backup`, but it does not yet have the atomic replacement helper needed before restore and destructive apply behaviors can safely land.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Atomic rewrite touches the same filesystem boundary that later restore and apply flows will depend on.
- Temp-file placement and rename behavior need to stay in the target directory so replacement remains atomic on Linux filesystems.

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

- Added `internal/backup/atomic.go` with an atomic rewrite helper that writes a temp file in the target directory, syncs it, renames it into place, and syncs the parent directory.
- Preserved existing file permissions when replacing an existing regular file and defaulted new files to mode `0600`.
- Added focused temp-directory tests for replacement, file creation, and invalid target-path cases.

Files changed:

- SESSION.md
- SESSIONS/027-atomic-rewrite.md
- internal/backup/atomic.go
- internal/backup/atomic_test.go

Files read:

- SESSION.md
- ROADMAP.md
- SKILLS/backup-restore.md
- docs/histkit-implementation-plan.md
- AGENT.md
- internal/backup/create.go
- internal/backup/create_test.go
- internal/backup/model.go
- SESSIONS/025-backup-model.md
- SESSIONS/026-backup-creation.md

Tests added:

- `TestRewriteAtomicReplacesContentsAndPreservesMode`
- `TestRewriteAtomicCreatesMissingFile`
- `TestRewriteAtomicRejectsEmptyPath`
- `TestRewriteAtomicRejectsNonRegularTarget`
- `TestRewriteAtomicRejectsMissingParentDirectory`

Tests run:

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- Keep this slice limited to a reusable atomic rewrite helper in `internal/backup`.
- Preserve existing target permissions when replacing a regular file.
- Use same-directory temp files and `os.Rename` for local atomic replacement.

Commands run:

- `git checkout -b 027-atomic-rewrite`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/backup-restore.md`
- `sed -n '1,220p' AGENT.md`
- `sed -n '1,220p' SESSIONS/025-backup-model.md`
- `sed -n '1,220p' SESSIONS/026-backup-creation.md`
- `sed -n '1,220p' internal/backup/model.go`
- `sed -n '1,260p' internal/backup/create.go`
- `sed -n '1,240p' internal/backup/create_test.go`
- `rg -n "atomic|rewrite|backup|restore|clean --apply|audit" -S docs internal README.md ROADMAP.md SESSIONS`
- `sed -n '687,730p' docs/histkit-implementation-plan.md`
- `sed -n '918,940p' docs/histkit-implementation-plan.md`
- `sed -n '150,170p' docs/histkit-implementation-plan.md`
- `sed -n '1,240p' SKILLS/testing.md`
- `gofmt -w internal/backup/atomic.go internal/backup/atomic_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Assumptions made:

- A package-level helper that atomically rewrites a caller-provided path is a safe precursor to restore and apply wiring.
- Defaulting new files to mode `0600` is conservative when no existing target file is present.

Risks introduced or reduced:

- Reduced: later restore and apply slices will be able to reuse a tested atomic-replace primitive instead of embedding file mutation logic ad hoc.
- Ongoing: restore semantics, audit logging, and cleanup apply wiring remain pending later slices.

Next recommended session:

- `028-audit-log`
