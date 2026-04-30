# Session 015: Pick Candidates

## Objective

Implement the initial merged picker candidate layer for history and snippets.

## Completed

- Added a recent-history query helper for picker use in `internal/index`.
- Added an `internal/picker` package that merges recent history and snippets into labeled candidates.
- Formatted candidates with exact `[history]` and `[snippet]` labels plus command text for later `fzf` integration.
- Added selection-line parsing for the same candidate format.
- Added deterministic tests for candidate merging, display formatting, selection-line parsing, and recent-history ordering.

## Files changed

- SESSION.md
- internal/index/picker.go
- internal/index/picker_test.go
- internal/picker/candidates.go
- internal/picker/candidates_test.go
- SESSIONS/015-pick-candidates.md

## Files read

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

## Tests added

- TestQueryRecentHistoryEntriesOrdersNewestFirst
- TestQueryRecentHistoryEntriesRequiresDBAndLimit
- TestLoadCandidatesMergesHistoryAndSnippets
- TestLoadCandidatesIncludesMissingBuiltinsWithoutOverwritingUserSnippets
- TestParseSelectedLine

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git status --short --branch
git checkout -b 015-pick-candidates
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SKILLS/fzf-picker.md
sed -n '630,670p' docs/histkit-implementation-plan.md
sed -n '1,220p' internal/index/stats.go
sed -n '1,260p' internal/index/writer.go
rg --files internal | rg 'picker|pick'
gofmt -w internal/index/picker.go internal/index/picker_test.go internal/picker/candidates.go internal/picker/candidates_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Format picker candidates with exact `[history]` and `[snippet]` labels followed by two spaces and the command text.
- Keep the candidate merger package-level only for this slice; `fzf` execution and CLI wiring remain deferred.
- When builtin snippets are included during candidate loading, user-store snippets win on ID collisions.

## Assumptions

- A package-level candidate merger without `fzf` execution is sufficient for this slice before the actual picker command exists.

## Known issues

- `fzf` invocation is still out of scope.
- No CLI `pick` command surface exists yet.
- Shell-wrapper integration is still out of scope.

## Risks reduced

- History and snippets can now be merged at presentation time without collapsing their underlying storage domains.

## Next recommended session

`016-fzf-picker`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
