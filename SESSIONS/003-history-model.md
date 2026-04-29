# Session 003: History Model

## Objective

Create the normalized history data model for histkit.

## Completed

- Added `internal/history/model.go` with the normalized `HistoryEntry` shape from the implementation plan.
- Added shell constants for the first supported shells: Bash and Zsh.
- Added a minimal `ParseWarning` type for upcoming parser slices.
- Added lightweight validation and optional-metadata helpers for deterministic model behavior.
- Added unit tests covering the model semantics.

## Files changed

- internal/history/model.go
- internal/history/model_test.go
- SESSION.md
- SESSIONS/003-history-model.md

## Tests added

- TestHistoryEntryValidate
- TestHistoryEntryValidateRequiresFields
- TestHistoryEntryOptionalMetadata
- TestParseWarningValidate
- TestParseWarningValidateRequiresFields

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Decisions

- Keep the history model in its own `internal/history` package before parser implementations arrive.
- Include a `ParseWarning` type in this slice so later parsers can report malformed lines without silently dropping them.

## Known issues

- No parser implementations exist yet.
- Hash generation and source detection remain deferred to later slices.

## Next recommended session

`004-bash-parser`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
