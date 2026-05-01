# Session 027: Atomic Rewrite

## Objective

Add atomic file-rewrite support for later restore and cleanup-apply work.

## Completed

- Added `internal/backup/atomic.go` with `RewriteAtomic`, a helper that writes to a temp file in the target directory, syncs it, renames it into place, and syncs the containing directory.
- Preserved existing file permissions when replacing an existing regular file.
- Defaulted newly created targets to mode `0600`.
- Added focused tests for successful replacement, new-file creation, and invalid target-path behavior.

## Files changed

- SESSION.md
- SESSIONS/027-atomic-rewrite.md
- internal/backup/atomic.go
- internal/backup/atomic_test.go

## Files read

- AGENT.md
- SESSION.md
- ROADMAP.md
- SKILLS/backup-restore.md
- SKILLS/testing.md
- docs/histkit-implementation-plan.md
- internal/backup/create.go
- internal/backup/create_test.go
- internal/backup/model.go
- SESSIONS/025-backup-model.md
- SESSIONS/026-backup-creation.md

## Tests added

- `TestRewriteAtomicReplacesContentsAndPreservesMode`
- `TestRewriteAtomicCreatesMissingFile`
- `TestRewriteAtomicRejectsEmptyPath`
- `TestRewriteAtomicRejectsNonRegularTarget`
- `TestRewriteAtomicRejectsMissingParentDirectory`

## Tests run

```bash
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 027-atomic-rewrite
sed -n '1,260p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,240p' SKILLS/backup-restore.md
sed -n '1,220p' SKILLS/testing.md
sed -n '1,220p' AGENT.md
sed -n '1,220p' SESSIONS/025-backup-model.md
sed -n '1,220p' SESSIONS/026-backup-creation.md
sed -n '1,260p' internal/backup/create.go
sed -n '1,240p' internal/backup/create_test.go
sed -n '1,220p' internal/backup/model.go
rg -n "atomic|rewrite|backup|restore|clean --apply|audit" -S docs internal README.md ROADMAP.md SESSIONS
sed -n '150,170p' docs/histkit-implementation-plan.md
sed -n '687,730p' docs/histkit-implementation-plan.md
sed -n '918,940p' docs/histkit-implementation-plan.md
gofmt -w internal/backup/atomic.go internal/backup/atomic_test.go
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...
```

## Decisions

- Keep this slice limited to a reusable atomic rewrite primitive without adding restore or apply wiring.
- Use a temp file in the target directory so `os.Rename` stays on the same filesystem.
- Preserve existing file permissions when replacing a regular file, and use `0600` for newly created targets.

## Assumptions

- A package-level helper in `internal/backup` is the right intermediate home before restore wiring exists.
- Failing when the parent directory is missing is acceptable at this slice because directory-creation policy belongs to higher-level restore or apply flows.

## Known issues

- Failure paths around directory sync or rename interruption are not directly unit tested.
- Restore semantics and forced-overwrite policy remain out of scope.

## Risks reduced

- The repo now has a tested atomic-replace primitive for safety-sensitive history mutation work.

## Next recommended session

`028-audit-log`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
