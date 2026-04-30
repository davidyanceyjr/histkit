# Session 022: Trivial Command Rules

## Objective

Add the initial curated built-in trivial-command and noise-rule catalog for history hygiene.

## Completed

- Added `internal/sanitize/trivial.go` with a curated built-in trivial/noise rule catalog.
- Added `MatchTrivialRules` for evaluating normalized history entries against that catalog.
- Covered exact trivial commands `clear`, `pwd`, `ls`, and `ll`.
- Reused the large-paste heuristic as hygiene noise and added false-positive guard coverage.

## Files changed

- SESSION.md
- SESSIONS/022-trivial-command-rules.md
- internal/sanitize/trivial.go
- internal/sanitize/trivial_test.go

## Files read

- SESSION.md
- ROADMAP.md
- SKILLS/sanitizer.md
- docs/histkit-implementation-plan.md
- README.md
- internal/sanitize/secrets.go
- internal/sanitize/secrets_test.go

## Tests added

- `TestBuiltinTrivialRulesValidate`
- `TestMatchTrivialRulesTruePositives`
- `TestMatchTrivialRulesFalsePositiveGuards`

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 022-trivial-command-rules
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SKILLS/sanitizer.md
rg -n "trivial|noise|cleanup|dedupe|drop_trivial|commands =|clear|pwd|ls|ll|large paste" docs/histkit-implementation-plan.md README.md -S
sed -n '300,360p' docs/histkit-implementation-plan.md
sed -n '40,48p' README.md
sed -n '396,408p' README.md
sed -n '1,240p' internal/sanitize/secrets.go
sed -n '1,260p' internal/sanitize/secrets_test.go
gofmt -w internal/sanitize/trivial.go internal/sanitize/trivial_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Keep the first trivial/noise catalog built-in and curated.
- Use delete for narrow exact-match trivial commands.
- Use quarantine for large accidental paste blobs as hygiene noise.

## Assumptions

- Exact-match-only trivial handling is the safest first cleanup slice.
- Large paste blobs can reasonably count as hygiene noise here.

## Known issues

- The trivial/noise catalog is not configurable yet.
- Dedupe and stale-entry logic are still out of scope.

## Risks reduced

- The sanitizer now has a distinct built-in hygiene rule set with false-positive guard coverage.

## Next recommended session

`023-dry-run-preview`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
