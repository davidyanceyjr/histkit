# Session 043: Command help detail

## Objective

Refresh `histkit <command> --help` output so each current command explains purpose, safe mode, and supported flags without changing command behavior.

## Completed

- Expanded help text for `scan`, `clean`, `pick`, `doctor`, `stats`, and `restore`.
- Kept the copy aligned with current behavior, especially for non-destructive defaults and backup/audit semantics.
- Added direct `--help` assertions for every current command.

## Files changed

- `internal/cli/scan.go`
- `internal/cli/clean.go`
- `internal/cli/pick.go`
- `internal/cli/doctor.go`
- `internal/cli/stats.go`
- `internal/cli/restore.go`
- `internal/cli/scan_test.go`
- `internal/cli/clean_test.go`
- `internal/cli/pick_test.go`
- `internal/cli/doctor_test.go`
- `internal/cli/stats_test.go`
- `internal/cli/restore_test.go`
- `SESSION.md`
- `SESSIONS/043-command-help-detail.md`

## Tests added

- `TestExecuteScanHelp`
- `TestExecuteCleanHelp`
- `TestExecutePickHelp`
- `TestExecuteDoctorHelp`
- `TestExecuteStatsHelp`
- `TestExecuteRestoreHelp`

## Tests run

```bash
env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli
env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...
env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go run ./cmd/histkit clean --help
```

## Results

Pass. The CLI tests and full repository suite passed, and the `clean --help` smoke check printed the revised usage text through the main entrypoint.

## Decisions

- Keep command help focused on purpose, safe behavior, and flags rather than full worked examples.
- State destructive boundaries explicitly where they matter most, especially in `clean` and `restore`.
- Assert help output per command so copy regressions are localized.

## Known issues

- Help text is duplicated across usage writers and tests, so future wording edits may justify a small formatting consolidation slice.

## Next recommended session

`044-help-format-consolidation`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
