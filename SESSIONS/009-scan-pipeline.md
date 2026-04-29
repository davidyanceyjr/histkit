# Session 009: Scan Pipeline

## Objective

Implement the initial `scan` pipeline from source detection through SQLite indexing.

## Completed

- Replaced the `scan` placeholder with an end-to-end read-only pipeline in `internal/cli/scan.go`.
- Added support for `--shell` and `--config` flags.
- Detected supported history sources, parsed them, initialized the SQLite schema, and wrote normalized entries through `internal/index`.
- Added compact scan summary output with source, parsed-entry, insert, skip, and warning counts.
- Added deterministic CLI tests that run against temporary home directories and temporary SQLite databases.

## Files changed

- SESSION.md
- internal/cli/root.go
- internal/cli/root_test.go
- internal/cli/scan.go
- internal/cli/scan_test.go
- SESSIONS/009-scan-pipeline.md

## Files read

- AGENT.md
- ROADMAP.md
- SESSION.md
- SKILLS/sqlite.md
- SKILLS/testing.md
- SKILLS/history-parsing.md
- SKILLS/go-cli.md
- docs/histkit-implementation-plan.md
- internal/config/config.go
- internal/config/config_test.go
- internal/cli/root.go
- internal/cli/root_test.go
- internal/cli/doctor.go
- internal/cli/stats.go
- internal/cli/scan.go
- internal/history/detect.go
- internal/history/bash.go
- internal/history/zsh.go
- internal/index/writer.go

## Tests added

- TestExecuteScanIndexesBashHistory
- TestExecuteScanShellFlagFiltersSources
- TestExecuteScanRejectsUnsupportedShell

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git fetch origin main
git merge --ff-only origin/main
git checkout -b 009-scan-pipeline
sed -n '1,240p' AGENT.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SESSION.md
sed -n '1,220p' SKILLS/sqlite.md
sed -n '1,220p' SKILLS/testing.md
sed -n '1,220p' SKILLS/history-parsing.md
sed -n '1,220p' SKILLS/go-cli.md
sed -n '1,220p' internal/cli/doctor.go
sed -n '1,220p' internal/cli/stats.go
rg -n "runScan|writeScanUsage|scan" internal/cli/*test.go internal/cli
sed -n '1,220p' internal/config/config_test.go
sed -n '1,220p' internal/config/config.go
sed -n '1,220p' internal/cli/root.go
sed -n '1,240p' internal/cli/root_test.go
sed -n '1,260p' internal/cli/scan.go
sed -n '1,220p' internal/history/detect.go
sed -n '1,220p' internal/history/bash.go
sed -n '1,240p' internal/history/zsh.go
sed -n '1,220p' internal/index/writer.go
sed -n '1,220p' docs/histkit-implementation-plan.md
gofmt -w internal/cli/scan.go internal/cli/scan_test.go internal/cli/root.go internal/cli/root_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Load an explicit config file only when `--config` is provided; otherwise scan uses default paths directly.
- Detect all supported history sources by default, and apply shell filtering only when `--shell` is provided.
- Initialize the SQLite schema on every scan before writing entries so first-run indexing stays self-contained.
- Report a compact scan summary with sources, parsed entries, inserted rows, skipped rows, and warnings.

## Assumptions

- A compact scan summary is sufficient for this slice; per-source or per-warning reporting can wait for later reporting work.

## Known issues

- `stats` is still a placeholder and does not expose indexed counts yet.
- Duplicate behavior remains defined by the current `(source_file, hash)` uniqueness boundary in the schema.

## Risks reduced

- `histkit scan` now exercises the real parser-to-index path under test instead of relying on a placeholder command.

## Next recommended session

`010-stats-command`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
