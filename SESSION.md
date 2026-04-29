# SESSION.md

## Current session

ID: `006-source-detection`

Status: completed

## Objective

Implement initial shell history source detection for Bash and Zsh.

## Scope

Implement:

- `internal/history/detect.go`
- canonical Bash and Zsh history source candidates
- filtering by supported shell name
- parser selection for detected sources
- deterministic tests for detection behavior

## Out of scope

- custom shell history path discovery
- CLI wiring for `--shell`
- SQLite schema
- index writing
- sanitization
- destructive cleanup

## Relevant skills

- `SKILLS/history-parsing.md`
- `SKILLS/testing.md`

## Acceptance criteria

- `go test ./...` passes
- supported history sources are detected from canonical paths
- missing source files are ignored deterministically
- unsupported shell names are rejected clearly
- parser selection is available for supported shells

## Current repo state

The CLI bootstrap, config/path package, history model, and Bash/Zsh parsers exist.

History source detection does not exist yet.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Keep source detection limited to canonical paths in this slice.
- Do not overreach into CLI integration or parser execution.
- Reject unsupported shell filters clearly.

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

#### Q001: Should initial source detection support only canonical `~/.bash_history` and `~/.zsh_history` paths?

- Area: source detection defaults
- Temporary assumption: yes, detect only canonical Bash and Zsh history paths in slice `006`
- Why non-blocking: configurable/custom history path support can be added later without changing parser correctness
- Why assumption is safe: it narrows behavior conservatively and does not mutate data
- Reversal cost: low
- Status: assumed-non-blocking

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

No questions answered yet.

## End-of-session notes

Summary:

- Added initial history source detection in `internal/history`.
- Added canonical Bash/Zsh source candidates, shell filtering, and parser lookup.
- Added deterministic detection tests and recorded the canonical-path assumption.

Files changed:

- internal/history/detect.go
- internal/history/detect_test.go
- docs/OPEN_QUESTIONS.md
- SESSIONS/006-source-detection.md

Tests added:

- TestCandidateSources
- TestCandidateSourcesRequiresHome
- TestDetectSourcesFindsExistingFiles
- TestDetectSourcesFiltersByShell
- TestDetectSourcesIgnoresMissingFiles
- TestDetectSourcesRejectsUnsupportedShell
- TestParserForShell
- TestParserForShellRejectsUnsupportedShell

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Detect only canonical `~/.bash_history` and `~/.zsh_history` files in the initial source-detection slice.
- Keep parser lookup alongside source detection so later scan wiring can resolve supported parsers cleanly.

Next recommended session:

- `007-sqlite-schema`
