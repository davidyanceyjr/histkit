# Session 018: Rule Model

## Objective

Add the initial sanitization rule and rule-match model types with validation.

## Completed

- Added `internal/sanitize/model.go` with rule-type, action, and confidence enums.
- Added `Rule` and `RuleMatch` structs for sanitizer model data.
- Added validation for exact, contains, regex, keyword-group, and heuristic rules.
- Added regex compilation checks and duplicate-rule detection.
- Added rule-match validation, including required redaction output for redact actions.

## Files changed

- SESSION.md
- SESSIONS/018-rule-model.md
- internal/sanitize/model.go
- internal/sanitize/model_test.go

## Files read

- SESSION.md
- ROADMAP.md
- SKILLS/sanitizer.md
- docs/histkit-implementation-plan.md
- README.md
- internal/snippets/model.go
- internal/snippets/model_test.go

## Tests added

- `TestRuleValidateAcceptsSupportedRuleTypes`
- `TestRuleValidateRejectsInvalidRules`
- `TestValidateRulesRejectsDuplicateNames`
- `TestValidateRulesAcceptsDistinctRules`
- `TestRuleMatchValidate`
- `TestRuleMatchValidateRejectsInvalidMatches`

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 018-rule-model
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SKILLS/sanitizer.md
rg -n "rule model|rules|sanitize|quarantine|redact|confidence|action|heuristic" docs/histkit-implementation-plan.md README.md internal -S
sed -n '430,620p' docs/histkit-implementation-plan.md
rg --files internal | rg 'sanitize|rule|clean'
sed -n '268,390p' docs/histkit-implementation-plan.md
sed -n '260,320p' README.md
sed -n '396,432p' README.md
sed -n '1,240p' internal/snippets/model.go
sed -n '1,260p' internal/snippets/model_test.go
gofmt -w internal/sanitize/model.go internal/sanitize/model_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Limit the first sanitizer slice to model and validation behavior only.
- Use `low`, `medium`, and `high` as the confidence set.
- Require redact matches to carry a transformed `After` value.

## Assumptions

- A `Detector` field is enough to identify heuristic rules before matcher logic exists.
- Model validation should reject invalid regex patterns immediately.

## Known issues

- No matcher implementation exists yet.
- Rules are not loaded from config yet.

## Risks reduced

- Later matching and preview slices now have a validated rule-model contract to build on.

## Next recommended session

`019-rule-matching`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
