# Session 025: Backup Model

## Objective

Add the initial backup record model and metadata helpers.

## Completed

- Added `internal/backup/model.go` with a backup record model and validation.
- Added deterministic backup ID generation.
- Added metadata-only backup record construction with derived backup paths.
- Added tests for required fields, ID generation, derived paths, and invalid inputs.

## Files changed

- SESSION.md
- SESSIONS/025-backup-model.md
- internal/backup/model.go
- internal/backup/model_test.go

## Files read

- SESSION.md
- ROADMAP.md
- SKILLS/backup-restore.md
- docs/histkit-implementation-plan.md
- README.md
- internal/sanitize/quarantine.go
- internal/sanitize/quarantine_test.go

## Tests added

- `TestRecordValidate`
- `TestRecordValidateRequiresFields`
- `TestBackupID`
- `TestBuildRecord`
- `TestBuildRecordRejectsInvalidInputs`

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 025-backup-model
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SKILLS/backup-restore.md
rg -n "backup|restore|checksum|backup_path|source_file|created_at|backup ID|audit" docs/histkit-implementation-plan.md README.md -S
sed -n '456,470p' docs/histkit-implementation-plan.md
sed -n '687,705p' docs/histkit-implementation-plan.md
sed -n '155,163p' README.md
rg --files internal | rg 'sanitize|backup|audit|restore'
sed -n '1,260p' internal/sanitize/quarantine.go
sed -n '1,260p' internal/sanitize/quarantine_test.go
gofmt -w internal/backup/model.go internal/backup/model_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Keep the first backup slice metadata-only.
- Derive backup paths beneath a per-backup directory keyed by backup ID.
- Use deterministic timestamp-plus-sequence IDs for now.

## Assumptions

- Metadata-only path derivation is sufficient before copy logic exists.
- External checksum supply is acceptable until checksum generation lands.

## Known issues

- No backup files are created yet.
- No restore behavior exists yet.

## Risks reduced

- The codebase now has a validated backup metadata contract for later creation and restore slices.

## Next recommended session

`026-backup-creation`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
