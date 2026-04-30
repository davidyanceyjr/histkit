# SESSION.md

## Current session

ID: `013-snippet-store`

Status: completed

## Objective

Implement the initial snippet store over a TOML user file.

## Scope

Implement:

- `internal/snippets/store.go`
- TOML-backed snippet load/save/list helpers
- add/remove snippet operations
- deterministic store tests

## Out of scope

- builtin snippet catalog
- picker integration
- snippet CLI commands
- backup tables
- destructive cleanup

## Relevant skills

- `SKILLS/testing.md`
- `SKILLS/snippets.md`
- `SKILLS/config.md`

## Acceptance criteria

- `go test ./...` passes
- snippet TOML files load into validated model values
- snippet store save/list/add/remove operations work in a temporary directory
- command templates are preserved exactly through round-trip storage
- snippet store remains separate from shell history/index storage

## Current repo state

The CLI bootstrap, config/path package, history model, Bash/Zsh parsers, source detection, SQLite schema initialization, history-entry writer, initial scan pipeline, initial stats command, initial doctor command, snippet model, and snippet store now exist.

There is still no builtin snippet catalog or picker integration yet.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Keep migrations explicit from the first schema version.
- Keep command behavior read-only outside explicit cleanup/apply milestones.
- Preserve separation between raw history, indexed history, and snippets.

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

- Added a TOML-backed snippet store with load/save/list/add/remove operations.
- Added snippet config defaults for the user snippet file path and default path tracking in `config.Paths`.
- Added deterministic tests for round-trip storage, duplicate-ID rejection, and snippet add/remove behavior.

Files changed:

- SESSION.md
- internal/config/config.go
- internal/config/config_test.go
- internal/snippets/store.go
- internal/snippets/store_test.go
- SESSIONS/013-snippet-store.md

Files read:

- AGENT.md
- ROADMAP.md
- SESSION.md
- SKILLS/testing.md
- SKILLS/snippets.md
- SKILLS/config.md
- README.md
- docs/histkit-implementation-plan.md
- internal/config/config.go
- internal/config/config_test.go
- internal/snippets/model.go
- internal/snippets/model_test.go

Tests added:

- TestStoreListMissingFileReturnsEmpty
- TestStoreSaveAndListRoundTrip
- TestStoreListRejectsDuplicateIDs
- TestStoreAddAndRemove
- TestStoreRemoveMissingIDFails

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Carry forward prior decisions from session `012-snippet-model`.
- Treat a missing user snippet file as an empty store rather than an error.
- Persist snippets as `[[snippets]]` TOML records using the existing model validation rules.
- Keep add/remove operations file-backed and in-memory ordered; do not introduce search or sorting semantics yet.

Commands run:

- `git status --short --branch`
- `git checkout -b 013-snippet-store`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/snippets.md`
- `sed -n '635,715p' docs/histkit-implementation-plan.md`
- `sed -n '390,430p' README.md`
- `sed -n '1,220p' internal/config/config.go`
- `sed -n '1,220p' internal/snippets/model.go`
- `sed -n '1,220p' internal/snippets/model_test.go`
- `rg -n "snippet store|snippets.toml|snippets list|snippets add|snippets remove|user_file|parse snippets TOML|search/list snippets" docs/histkit-implementation-plan.md README.md SKILLS`
- `sed -n '618,680p' docs/histkit-implementation-plan.md`
- `sed -n '190,225p' README.md`
- `sed -n '1,220p' internal/config/config_test.go`
- `sed -n '1,220p' SKILLS/config.md`
- `gofmt -w internal/config/config.go internal/config/config_test.go internal/snippets/store.go internal/snippets/store_test.go`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- A single user TOML file is sufficient for the first snippet-store slice before builtin catalogs are introduced.

Risks introduced or reduced:

- Reduced: snippet persistence now exists in a store separate from shell history and SQLite index data.
- Ongoing: builtin snippet catalogs and snippet CLI surfaces are still deferred.

Next recommended session:

- `014-builtin-snippets`
