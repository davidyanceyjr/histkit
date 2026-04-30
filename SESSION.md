# SESSION.md

## Current session

ID: `015-pick-candidates`

Status: completed

## Objective

Implement the initial merged picker candidate layer for history and snippets.

## Scope

Implement:

- `internal/picker/candidates.go`
- recent-history query helper for picker use
- merged history/snippet candidate formatting
- deterministic picker candidate tests

## Out of scope

- `fzf` invocation
- CLI `pick` command surface
- snippet CLI commands
- backup tables
- destructive cleanup

## Relevant skills

- `SKILLS/testing.md`
- `SKILLS/snippets.md`
- `SKILLS/fzf-picker.md`

## Acceptance criteria

- `go test ./...` passes
- recent history entries can be queried for picker use
- history and snippet candidates merge into one in-memory stream
- candidate formatting uses `[history]` and `[snippet]` labels
- snippets remain separate from shell history/index storage

## Current repo state

The CLI bootstrap, config/path package, history model, Bash/Zsh parsers, source detection, SQLite schema initialization, history-entry writer, initial scan pipeline, initial stats command, initial doctor command, snippet model, snippet store, builtin snippet catalog, and picker candidate layer now exist.

There is still no `fzf` invocation or shell-wrapper integration yet.

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

- Added a recent-history query helper for picker use in `internal/index`.
- Added an `internal/picker` package that merges recent history and snippets into labeled candidates.
- Added deterministic tests for candidate merging, display formatting, selection-line parsing, and recent-history ordering.

Files changed:

- SESSION.md
- internal/index/picker.go
- internal/index/picker_test.go
- internal/picker/candidates.go
- internal/picker/candidates_test.go
- SESSIONS/015-pick-candidates.md

Files read:

- AGENT.md
- ROADMAP.md
- SESSION.md
- SKILLS/snippets.md
- SKILLS/fzf-picker.md
- README.md
- docs/histkit-implementation-plan.md
- internal/index/stats.go
- internal/index/writer.go
- internal/snippets/model.go
- internal/snippets/store.go
- internal/snippets/store_test.go

Tests added:

- TestQueryRecentHistoryEntriesOrdersNewestFirst
- TestQueryRecentHistoryEntriesRequiresDBAndLimit
- TestLoadCandidatesMergesHistoryAndSnippets
- TestLoadCandidatesIncludesMissingBuiltinsWithoutOverwritingUserSnippets
- TestParseSelectedLine

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Carry forward prior decisions from session `014-builtin-snippets`.
- Format picker candidates with exact `[history]` and `[snippet]` labels followed by two spaces and the command text.
- Keep the candidate merger package-level only for this slice; `fzf` execution and CLI wiring remain deferred.
- When builtin snippets are included during candidate loading, user-store snippets win on ID collisions.

Commands run:

- `git status --short --branch`
- `git checkout -b 015-pick-candidates`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/fzf-picker.md`
- `sed -n '630,670p' docs/histkit-implementation-plan.md`
- `sed -n '1,220p' internal/index/stats.go`
- `sed -n '1,260p' internal/index/writer.go`
- `rg --files internal | rg 'picker|pick'`
- `gofmt -w internal/index/picker.go internal/index/picker_test.go internal/picker/candidates.go internal/picker/candidates_test.go`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- A package-level candidate merger without `fzf` execution is sufficient for this slice before the actual picker command exists.

Risks introduced or reduced:

- Reduced: history and snippets can now be merged at presentation time without collapsing their underlying storage domains.
- Ongoing: the actual picker command, `fzf` process handling, and shell wrapper integration are still deferred.

Next recommended session:

- `016-fzf-picker`
