# SESSION.md

## Current session

ID: `022-trivial-command-rules`

Status: completed

## Objective

Add the initial curated built-in trivial-command and noise-rule catalog for history hygiene.

## Scope

Implement:

- built-in trivial/noise rule definitions
- helper to evaluate normalized history entries against that built-in hygiene rule set
- tests for true positives and false-positive guards

## Out of scope

- user-configurable cleanup rule loading
- dry-run preview output
- quarantine persistence
- cleanup apply behavior
- dedupe or stale-entry logic

## Relevant skills

- `SKILLS/sanitizer.md`
- `SKILLS/testing.md`

## Acceptance criteria

- repository contains a curated initial trivial/noise rule catalog
- initial detections cover trivial commands like `clear`, `pwd`, `ls`, and `ll`
- large accidental paste blobs remain classed as hygiene noise
- broad false positives remain guarded against
- `go test ./...` passes

## Current repo state

The repository now includes a curated built-in trivial-command and noise rule catalog in `internal/sanitize` plus a `MatchTrivialRules` helper that applies those rules to normalized history entries.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- The built-in trivial catalog is still intentionally narrow and not yet configurable from TOML.
- Exact-match trivial rules deliberately avoid matching compound commands like `pwd && ls`.

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

- Added `internal/sanitize/trivial.go` with the first curated built-in trivial/noise rule catalog.
- Added `MatchTrivialRules` to run normalized history entries through that catalog.
- Added tests for expected trivial detections and false-positive guard cases.

Files changed:

- internal/sanitize/trivial.go
- internal/sanitize/trivial_test.go
- SESSION.md
- SESSIONS/022-trivial-command-rules.md

Files read:

- SESSION.md
- ROADMAP.md
- SKILLS/sanitizer.md
- docs/histkit-implementation-plan.md
- README.md
- internal/sanitize/secrets.go
- internal/sanitize/secrets_test.go

Tests added:

- `TestBuiltinTrivialRulesValidate`
- `TestMatchTrivialRulesTruePositives`
- `TestMatchTrivialRulesFalsePositiveGuards`

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Keep the initial trivial/noise catalog built-in and curated rather than loading from config yet.
- Use delete for narrow exact-match trivial commands.
- Use quarantine for large accidental paste blobs as hygiene noise.

Commands run:

- `git checkout -b 022-trivial-command-rules`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/sanitizer.md`
- `rg -n "trivial|noise|cleanup|dedupe|drop_trivial|commands =|clear|pwd|ls|ll|large paste" docs/histkit-implementation-plan.md README.md -S`
- `sed -n '300,360p' docs/histkit-implementation-plan.md`
- `sed -n '40,48p' README.md`
- `sed -n '396,408p' README.md`
- `sed -n '1,240p' internal/sanitize/secrets.go`
- `sed -n '1,260p' internal/sanitize/secrets_test.go`
- `gofmt -w internal/sanitize/trivial.go internal/sanitize/trivial_test.go`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- Exact-match-only trivial command handling is the safest first slice for noise cleanup.
- Treating large paste blobs as hygiene noise is acceptable even though they can also overlap with security concerns.

Risks introduced or reduced:

- Reduced: the sanitizer now has a separate built-in hygiene rule set rather than overloading secret detection for cleanup noise.
- Ongoing: broader hygiene logic like dedupe or stale-entry handling still needs later slices.

Next recommended session:

- `023-dry-run-preview`
