# SESSION.md

## Current session

ID: `009-scan-pipeline`

Status: completed

## Objective

Implement the initial `scan` pipeline from source detection through SQLite indexing.

## Scope

Implement:

- `internal/cli/scan.go`
- shell source detection, parsing, and SQLite indexing flow
- scan summary output
- deterministic CLI tests using temporary home directories

## Out of scope

- stats queries
- run metadata writer logic
- sanitizer metadata tables
- snippet tables
- backup tables
- destructive cleanup

## Relevant skills

- `SKILLS/sqlite.md`
- `SKILLS/testing.md`
- `SKILLS/history-parsing.md`
- `SKILLS/go-cli.md`

## Acceptance criteria

- `go test ./...` passes
- `histkit scan` detects supported history sources and writes parsed entries to SQLite
- `histkit scan --shell <shell>` limits ingest to one shell
- the scan path remains non-destructive and only reads history files
- deterministic CLI tests cover indexing through a temporary home directory

## Current repo state

The CLI bootstrap, config/path package, history model, Bash/Zsh parsers, source detection, SQLite schema initialization, history-entry writer, and the initial scan pipeline now exist.

The stats command is still a placeholder and no aggregated reporting exists yet.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Keep migrations explicit from the first schema version.
- Keep scan behavior conservative and read-only.
- Preserve safe nullability for metadata that parsers may not always populate.

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

- Replaced the scan placeholder with an end-to-end source-detection, parse, schema-init, and index-write pipeline.
- Added `scan` support for `--shell` and `--config` flags with conservative, read-only behavior.
- Added temporary-home CLI tests that verify indexing, shell filtering, and unsupported-shell errors.

Files changed:

- SESSION.md
- internal/cli/root.go
- internal/cli/root_test.go
- internal/cli/scan.go
- internal/cli/scan_test.go
- SESSIONS/009-scan-pipeline.md

Files read:

- AGENT.md
- ROADMAP.md
- SESSION.md
- SKILLS/sqlite.md
- SKILLS/testing.md
- SKILLS/history-parsing.md
- SKILLS/go-cli.md
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
- docs/histkit-implementation-plan.md

Tests added:

- TestExecuteScanIndexesBashHistory
- TestExecuteScanShellFlagFiltersSources
- TestExecuteScanRejectsUnsupportedShell

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Load an explicit config file only when `--config` is provided; otherwise scan uses default paths directly.
- Detect all supported history sources by default, and apply shell filtering only when `--shell` is provided.
- Initialize the SQLite schema on every scan before writing entries so first-run indexing stays self-contained.
- Report a compact scan summary with sources, parsed entries, inserted rows, skipped rows, and warnings.

Commands run:

- `git status --short --branch`
- `git fetch origin main`
- `git merge --ff-only origin/main`
- `git checkout -b 009-scan-pipeline`
- `sed -n '1,240p' AGENT.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' SKILLS/sqlite.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '1,220p' SKILLS/history-parsing.md`
- `sed -n '1,220p' SKILLS/go-cli.md`
- `sed -n '1,220p' internal/cli/doctor.go`
- `sed -n '1,220p' internal/cli/stats.go`
- `sed -n '1,220p' internal/config/config_test.go`
- `sed -n '1,220p' internal/cli/root.go`
- `sed -n '1,240p' internal/cli/root_test.go`
- `rg -n "runScan|writeScanUsage|scan" internal/cli/*test.go internal/cli`
- `sed -n '1,260p' internal/cli/scan.go`
- `sed -n '1,260p' internal/config/config.go`
- `sed -n '1,220p' internal/history/detect.go`
- `sed -n '1,220p' internal/history/bash.go`
- `sed -n '1,240p' internal/history/zsh.go`
- `sed -n '1,220p' internal/index/writer.go`
- `sed -n '1,220p' docs/histkit-implementation-plan.md`
- `gofmt -w internal/cli/scan.go internal/cli/scan_test.go internal/cli/root.go internal/cli/root_test.go`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- A compact scan summary is sufficient for this slice; per-source or per-warning reporting can wait for later reporting work.

Risks introduced or reduced:

- Reduced: `histkit scan` now exercises the real parser-to-index path under test instead of relying on a placeholder command.
- Ongoing: duplicate handling is still defined by the current `(source_file, hash)` schema and may collapse repeated identical commands within one source file.

Next recommended session:

- `010-stats-command`
