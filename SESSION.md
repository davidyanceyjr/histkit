# SESSION.md

## Current session

ID: `003-history-model`

Status: completed

## Objective

Create the normalized history data model for histkit.

## Scope

Implement:

- `internal/history/model.go`
- normalized `HistoryEntry` type
- shell/source metadata fields
- parse-warning metadata types needed by later parser slices
- tests for model semantics

## Out of scope

- real history parsing
- source detection
- SQLite schema
- index writing
- sanitization
- fzf integration
- backups
- systemd units
- destructive cleanup

## Relevant skills

- `SKILLS/history-parsing.md`
- `SKILLS/testing.md`

## Acceptance criteria

- `go test ./...` passes
- `HistoryEntry` model exists with the planned normalized fields
- optional timestamp/exit metadata are represented safely
- model helpers are covered by deterministic unit tests

## Current repo state

The CLI bootstrap and initial config/path packages exist.

The normalized history model does not exist yet.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Avoid embedding parser behavior into the model slice.
- Preserve raw history data without mutating or normalizing away source context.
- Keep the model extensible for Bash and Zsh parser slices.

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

- Added the normalized history model in `internal/history`.
- Added shell constants plus a parse-warning type for later parser slices.
- Added validation and optional-metadata helpers with deterministic unit tests.

Files changed:

- internal/history/model.go
- internal/history/model_test.go
- SESSIONS/003-history-model.md

Tests added:

- TestHistoryEntryValidate
- TestHistoryEntryValidateRequiresFields
- TestHistoryEntryOptionalMetadata
- TestParseWarningValidate
- TestParseWarningValidateRequiresFields

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Keep the history model in its own `internal/history` package before parser implementations arrive.
- Include a `ParseWarning` type in this slice so later parsers can report malformed lines without silently dropping them.

Next recommended session:

- `004-bash-parser`
