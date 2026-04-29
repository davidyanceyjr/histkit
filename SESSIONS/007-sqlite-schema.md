# Session 007: SQLite Schema

## Objective

Implement the initial SQLite schema for normalized history entries and run metadata.

## Completed

- Added `internal/index/schema.go` with explicit schema initialization.
- Added the initial `history_entries` and `runs` tables.
- Added a partial unique index to dedupe `history_entries` by `source_file` and `hash`.
- Added schema versioning through `PRAGMA user_version`.
- Added deterministic schema tests using temporary databases.
- Added the SQLite driver dependency.

## Files changed

- go.mod
- go.sum
- internal/index/schema.go
- internal/index/schema_test.go
- SESSION.md
- SESSIONS/007-sqlite-schema.md

## Tests added

- TestOpenRequiresPath
- TestInitSchemaCreatesTables
- TestHistoryEntriesDedupeBySourceAndHash
- TestRunsInsertAndReopenDatabase
- TestInitSchemaRequiresDB

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go mod tidy
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Decisions

- Use `modernc.org/sqlite` for the initial SQLite schema layer.
- Start schema versioning immediately with `PRAGMA user_version = 1`.
- Enforce dedupe through a partial unique index on `(source_file, hash)` only when `hash` is present.

## Known issues

- No history-entry writer exists yet.
- No stats query layer exists yet.
- Deferred tables such as `history_actions`, `snippets`, and `backups` remain out of scope.

## Next recommended session

`008-index-writer`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
