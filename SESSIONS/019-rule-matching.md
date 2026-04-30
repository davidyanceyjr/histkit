# Session 019: Rule Matching

## Objective

Add the initial rule-matching engine for sanitizer rules over normalized history entries.

## Completed

- Added `internal/sanitize/matcher.go` with `MatchRule` and `MatchEntry`.
- Implemented exact, contains, regex, keyword-group, and heuristic rule matching.
- Added a small heuristic detector registry with `high_entropy_token` and `large_paste_blob`.
- Added matcher-focused tests for true positives, false-positive guard cases, invalid detector handling, and multi-rule matching.

## Files changed

- SESSION.md
- SESSIONS/019-rule-matching.md
- internal/sanitize/matcher.go
- internal/sanitize/matcher_test.go

## Files read

- SESSION.md
- ROADMAP.md
- SKILLS/sanitizer.md
- docs/histkit-implementation-plan.md
- README.md
- internal/history/model.go
- internal/sanitize/model.go
- internal/sanitize/model_test.go

## Tests added

- `TestMatchRule`
- `TestMatchRuleRejectsUnknownHeuristicDetector`
- `TestMatchEntryReturnsAllMatchingRules`
- `TestMatchEntryRejectsInvalidEntry`
- `TestLargePasteBlobHeuristic`

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 019-rule-matching
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SKILLS/sanitizer.md
sed -n '528,586p' docs/histkit-implementation-plan.md
sed -n '1,240p' internal/history/model.go
sed -n '1,260p' internal/sanitize/model.go
sed -n '1,320p' internal/sanitize/model_test.go
rg -n "Match|matcher|heuristic|keyword group|contains|regex" internal/sanitize README.md docs/histkit-implementation-plan.md -S
rg -n "keyword group|heuristic detector|all keywords|keywords" README.md docs/histkit-implementation-plan.md -S
rg -n "high[- ]entropy|large paste|private key|bearer token|inline password|URL" docs/histkit-implementation-plan.md README.md SKILLS/sanitizer.md -S
gofmt -w internal/sanitize/matcher.go internal/sanitize/matcher_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Keyword-group rules require all listed keywords to be present.
- Match over normalized command text rather than raw history lines.
- Keep redact matches structurally valid by copying the original command into `After` until transform logic exists.

## Assumptions

- A small detector registry is enough for the first heuristic-matching slice.
- Exact matches should remain byte-for-byte exact.

## Known issues

- Redaction output is still placeholder-only.
- No config wiring exists for user-defined cleanup rules yet.

## Risks reduced

- The sanitizer engine now has an executable matcher contract for later transform and preview slices.

## Next recommended session

`020-redaction-transforms`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
