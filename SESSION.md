# SESSION.md

## Current session

ID: `014-builtin-snippets`

Status: completed

## Objective

Implement the initial builtin snippet catalog.

## Scope

Implement:

- `internal/snippets/builtin.go`
- curated builtin snippet definitions
- builtin import helper for the user snippet store
- deterministic builtin tests

## Out of scope

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
- builtin snippets validate as normal snippet model values
- builtin import seeds missing snippets into the user store
- existing user snippets are not overwritten by builtin import
- builtin snippets remain separate from shell history/index storage

## Current repo state

The CLI bootstrap, config/path package, history model, Bash/Zsh parsers, source detection, SQLite schema initialization, history-entry writer, initial scan pipeline, initial stats command, initial doctor command, snippet model, snippet store, and builtin snippet catalog now exist.

There is still no picker integration or snippet CLI surface yet.

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

- Added a curated builtin snippet catalog in `internal/snippets`.
- Added a builtin import helper that seeds missing snippets into the user store without overwriting existing user snippets.
- Added deterministic tests for builtin validation, initial import, non-overwriting import behavior, and idempotent imports.

Files changed:

- SESSION.md
- internal/snippets/builtin.go
- internal/snippets/builtin_test.go
- SESSIONS/014-builtin-snippets.md

Files read:

- AGENT.md
- ROADMAP.md
- SESSION.md
- SKILLS/testing.md
- SKILLS/snippets.md
- README.md
- docs/histkit-implementation-plan.md
- internal/snippets/model.go
- internal/snippets/model_test.go
- internal/snippets/store.go
- internal/snippets/store_test.go

Tests added:

- TestBuiltinsValidate
- TestImportBuiltinsSeedsMissingSnippets
- TestImportBuiltinsDoesNotOverwriteExistingSnippet
- TestImportBuiltinsIdempotent

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Carry forward prior decisions from session `013-snippet-store`.
- Keep builtin snippets as ordinary `Snippet` values validated through the same model rules as user snippets.
- Import builtins only when their IDs are missing from the user store; existing user snippets always win.
- Keep the builtin import helper idempotent so repeated imports do not reorder or duplicate entries.

Commands run:

- `git status --short --branch`
- `git checkout -b 014-builtin-snippets`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/snippets.md`
- `rg -n "builtin snippets|builtin = true|Builtin snippets|import builtins|builtin catalog|snippets]\" README.md docs/histkit-implementation-plan.md SKILLS internal`
- `sed -n '330,345p' README.md`
- `sed -n '805,825p' docs/histkit-implementation-plan.md`
- `sed -n '1,220p' internal/snippets/model.go`
- `sed -n '1,240p' internal/snippets/store.go`
- `sed -n '1,240p' internal/snippets/store_test.go`
- `gofmt -w internal/snippets/builtin.go internal/snippets/builtin_test.go`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- A small curated builtin catalog is sufficient for the first builtin slice before any CLI import surface exists.

Risks introduced or reduced:

- Reduced: builtin snippets now have a deterministic catalog and a safe import path that preserves user-owned entries.
- Ongoing: no CLI surface exists yet for listing or importing builtins directly.

Next recommended session:

- `015-pick-candidates`
