# SESSION.md

## Current session

ID: `004-bash-parser`

Status: completed

## Objective

Implement the initial Bash history parser for plain line-oriented history files.

## Scope

Implement:

- `internal/history/bash.go`
- plain Bash history parsing from an `io.Reader`
- fixture-driven tests for Bash history parsing
- warning behavior for invalid Bash lines

## Out of scope

- `HISTTIMEFORMAT` parsing
- multiline reconstruction
- source detection
- Zsh parsing
- SQLite schema
- index writing
- sanitization
- destructive cleanup

## Relevant skills

- `SKILLS/history-parsing.md`
- `SKILLS/testing.md`

## Acceptance criteria

- `go test ./...` passes
- plain Bash history lines parse into normalized `HistoryEntry` values
- empty lines are handled deterministically
- parser preserves raw lines and command text
- no history files are modified

## Current repo state

The CLI bootstrap, config/path package, and normalized history model exist.

No shell history parser exists yet.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Do not infer Bash timestamps from ambiguous formats in this slice.
- Preserve raw lines without mutating command text.
- Avoid overreaching into source detection or serialization.

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

- Added the initial plain-line Bash parser in `internal/history`.
- Preserved raw lines and command text for valid entries.
- Added warning behavior for whitespace-only lines and fixture-driven parser tests.

Files changed:

- internal/history/bash.go
- internal/history/bash_test.go
- testdata/history/bash/plain.hist
- SESSIONS/004-bash-parser.md

Tests added:

- TestParseBashFixture
- TestParseBashEmptyInput
- TestParseBashRequiresSourceFile
- TestParseBashRequiresReader

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Treat empty Bash history lines as ignorable input.
- Treat whitespace-only Bash history lines as parse warnings so they are not silently converted into invalid commands.

Next recommended session:

- `005-zsh-parser`
