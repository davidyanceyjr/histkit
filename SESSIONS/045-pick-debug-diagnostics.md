# Session 045: Pick Debug Diagnostics

## Objective

Add opt-in diagnostics for `histkit pick` so a user can identify whether the command is stalling before database access, candidate loading, or `fzf` launch.

## Completed

- Added a `--debug` flag to `histkit pick`.
- Routed `pick` stderr through the CLI dispatcher so diagnostics can be emitted without affecting stdout.
- Logged stage-level start and completion timing for the pre-`fzf` path and selection outcome.
- Added a deterministic CLI test covering the debug diagnostics path.

## Files changed

- `internal/cli/root.go`
- `internal/cli/pick.go`
- `internal/cli/pick_test.go`
- `SESSION.md`
- `SESSIONS/045-pick-debug-diagnostics.md`

## Tests added

- `TestExecutePickDebugWritesStageDiagnosticsToStderr`

## Tests run

```bash
env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli
env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...
```

## Results

All targeted and full repository tests passed. Normal `pick` behavior remained silent on stderr, and `pick --debug` emitted the expected stage diagnostics.

## Decisions

- Keep diagnostics opt-in behind `histkit pick --debug`.
- Write diagnostics to stderr so selected commands still emit cleanly to stdout.
- Focus this slice on identification of the blocking stage, not on remediation of SQLite or terminal issues.

## Known issues

- The new diagnostics identify the stage boundary but do not yet add lock timeouts, retries, or automatic remediation.

## Next recommended session

`046-pick-stall-root-cause-follow-up`


## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
