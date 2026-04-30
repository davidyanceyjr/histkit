# Session 023: Dry-Run Preview

## Objective

Add the initial non-destructive dry-run preview layer for cleanup matches.

## Completed

- Added `internal/sanitize/preview.go` with preview report types and preview generation helpers.
- Combined built-in secret and trivial rules into a package-level dry-run preview flow.
- Added text rendering that reports matched rule, confidence, action, original value, transformed value when redacted, and reason.
- Added tests for no-mutation behavior, counts-by-action, required rendered fields, and no-match handling.

## Files changed

- SESSION.md
- SESSIONS/023-dry-run-preview.md
- internal/sanitize/preview.go
- internal/sanitize/preview_test.go

## Files read

- SESSION.md
- ROADMAP.md
- SKILLS/sanitizer.md
- docs/histkit-implementation-plan.md
- README.md
- internal/sanitize/model.go
- internal/sanitize/secrets.go
- internal/sanitize/trivial.go

## Tests added

- `TestPreviewEntryReturnsMatchesWithoutMutatingEntry`
- `TestPreviewEntriesCountsActionsAcrossSecretAndTrivialRules`
- `TestRenderPreviewTextIncludesRequiredFields`
- `TestRenderPreviewTextHandlesNoMatches`

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 023-dry-run-preview
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SKILLS/sanitizer.md
rg -n "dry-run|preview|matched rule|reason|action|transformed value|counts by action|confidence" docs/histkit-implementation-plan.md README.md -S
sed -n '570,586p' docs/histkit-implementation-plan.md
sed -n '714,722p' docs/histkit-implementation-plan.md
sed -n '298,304p' README.md
rg --files internal/sanitize
sed -n '1,260p' internal/sanitize/model.go
sed -n '1,260p' internal/sanitize/secrets.go
sed -n '1,260p' internal/sanitize/trivial.go
gofmt -w internal/sanitize/preview.go internal/sanitize/preview_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Preview generation should aggregate both secret and trivial rule catalogs.
- Preview items should include only matched entries.
- Text rendering should include transformed output only for redact actions.

## Assumptions

- Package-level preview generation is sufficient before CLI wiring.
- Simple reviewable text output is enough for the first preview slice.

## Known issues

- No `clean --dry-run` command surface exists yet.
- Preview formatting is intentionally minimal.

## Risks reduced

- The sanitizer now has a non-destructive preview layer suitable for later CLI integration.

## Next recommended session

`024-quarantine-records`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
