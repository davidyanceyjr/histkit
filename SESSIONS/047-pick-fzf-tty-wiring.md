# Session 047: Pick fzf TTY Wiring

## Objective

Make `histkit pick` mirror `fzf` interactive output to the controlling terminal while preserving stdout capture for the selected command.

## Completed

- Mirrored `fzf` stderr to `/dev/tty` when available.
- Preserved stdout buffering so the selected command still returns cleanly from `histkit pick`.
- Added picker tests for TTY mirroring and no-TTY fallback behavior.
- Carried the uncommitted `046` planning artifacts forward in this slice.

## Files changed

- `internal/picker/fzf.go`
- `internal/picker/fzf_test.go`
- `SESSION.md`
- `SESSIONS/046-pick-fzf-tty-wiring-plan.md`
- `SESSIONS/047-pick-fzf-tty-wiring.md`

## Tests added

- `TestSelectMirrorsFZFStderrToTTYWhenAvailable`
- `TestSelectReturnsCapturedErrorWhenTTYUnavailable`

## Tests run

```bash
env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/picker
env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli
env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...
```

## Results

All targeted and full repository tests passed. The picker now mirrors `fzf` stderr to the controlling terminal when one is available and still captures stdout for the selected command.

## Decisions

- Use `/dev/tty` explicitly for the interactive `fzf` stderr channel.
- Keep the change inside the picker layer and avoid shell wrapper changes.
- Preserve internal stderr buffering so returned errors still include `fzf` output.

## Known issues

- The fix mirrors `fzf` stderr to the TTY, but if a specific environment also requires broader stdio inheritance, an additional follow-up slice may still be needed.

## Next recommended session

`048-pick-tty-follow-up`


## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

- Should the implementation prefer inheriting parent stdio directly or opening `/dev/tty` explicitly for `fzf`?
  - Answered by choosing `/dev/tty` for interactive stderr mirroring in `internal/picker/fzf.go`.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
