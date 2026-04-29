# SESSION.md

## Current session

ID: `007-sqlite-schema`

Status: completed

## Objective

Implement the initial SQLite schema for normalized history entries and run metadata.

## Scope

Implement:

- `internal/index/schema.go`
- explicit schema initialization
- initial `history_entries` and `runs` tables
- deterministic schema tests using temporary databases

## Out of scope

- scan/index writer logic
- stats queries beyond what schema tests need
- sanitizer metadata tables
- snippet tables
- backup tables
- destructive cleanup

## Relevant skills

- `SKILLS/sqlite.md`
- `SKILLS/testing.md`

## Acceptance criteria

- `go test ./...` passes
- SQLite schema initializes successfully in a temporary database
- `history_entries` and `runs` tables exist
- dedupe by `source_file` and `hash` is enforced
- database can be reopened after initialization

## Current repo state

The CLI bootstrap, config/path package, history model, Bash/Zsh parsers, and source detection exist.

SQLite schema initialization does not exist yet.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Keep migrations explicit from the first schema version.
- Do not introduce index-writing behavior early.
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

- Added the initial SQLite schema layer in `internal/index`.
- Added `history_entries` and `runs` tables plus a dedupe index for source/hash.
- Added temporary-database schema tests and the SQLite driver dependency.

Files changed:

- go.mod
- go.sum
- internal/index/schema.go
- internal/index/schema_test.go
- SESSIONS/007-sqlite-schema.md

Tests added:

- TestOpenRequiresPath
- TestInitSchemaCreatesTables
- TestHistoryEntriesDedupeBySourceAndHash
- TestRunsInsertAndReopenDatabase
- TestInitSchemaRequiresDB

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go mod tidy`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Use `modernc.org/sqlite` for the initial SQLite schema layer.
- Start schema versioning immediately with `PRAGMA user_version = 1`.
- Enforce dedupe through a partial unique index on `(source_file, hash)` only when `hash` is present.

Next recommended session:

- `008-index-writer`
