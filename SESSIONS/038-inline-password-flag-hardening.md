# Session 038: Inline Password Flag Hardening

## Objective

Harden inline password detection and reduce false positives from high-entropy token matching.

## Completed

- Replaced the built-in inline password regex with a heuristic detector that matches `--password`, `--passwd`, and MySQL-family `-psecret` forms without treating arbitrary `-p` flags as secrets.
- Updated secret redaction so password-bearing commands keep their flag structure while masking only the sensitive value.
- Narrowed the high-entropy heuristic to sensitive key/value contexts, which prevents routine package names, filesystem paths, and device labels from being quarantined.
- Updated sanitizer and CLI tests to lock in the narrower matching and the more reviewable redaction output.
- Folded the planned `039-high-entropy-token-false-positive-guards` work into the same session because it touched the same built-in sanitizer boundary.

## Files changed

- SESSION.md
- SESSIONS/038-inline-password-flag-hardening.md
- internal/cli/clean_test.go
- internal/sanitize/apply_test.go
- internal/sanitize/matcher.go
- internal/sanitize/preview_test.go
- internal/sanitize/redact.go
- internal/sanitize/redact_test.go
- internal/sanitize/secrets.go
- internal/sanitize/secrets_test.go

## Files read

- SESSION.md
- ROADMAP.md
- SKILLS/sanitizer.md
- SKILLS/testing.md
- docs/histkit-implementation-plan.md
- RISKS.md
- SESSIONS/021-secret-rules.md
- internal/sanitize/secrets.go
- internal/sanitize/matcher.go
- internal/sanitize/redact.go
- internal/sanitize/secrets_test.go
- internal/sanitize/redact_test.go
- internal/sanitize/preview_test.go
- internal/sanitize/apply_test.go
- internal/cli/clean_test.go

## Tests added

- MySQL short `-psecret` coverage for the inline password detector
- False-positive guards for `ssh -p`, `grep -p`, package names, and device-label/path commands
- Password redaction assertions that preserve the flag while masking only the secret value
- High-entropy sensitive-flag redaction coverage for `--session-token` style arguments

## Tests run

```bash
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/sanitize
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git status --short --branch
sed -n '1,220p' SKILLS/sanitizer.md
sed -n '1,220p' SKILLS/testing.md
sed -n '1,240p' ROADMAP.md
rg -n "password|passwd|entropy|token|secret|credential|bearer|apikey|api_key|inline" internal/sanitize internal -g '!**/*_test.go'
sed -n '1,260p' internal/sanitize/secrets.go
sed -n '1,260p' internal/sanitize/matcher.go
sed -n '1,260p' internal/sanitize/redact.go
sed -n '1,260p' internal/sanitize/secrets_test.go
sed -n '1,260p' internal/sanitize/redact_test.go
sed -n '1,260p' internal/sanitize/preview_test.go
sed -n '1,140p' internal/sanitize/apply_test.go
sed -n '1,240p' internal/cli/clean_test.go
sed -n '1,220p' SESSIONS/021-secret-rules.md
sed -n '1,220p' docs/histkit-implementation-plan.md
sed -n '1,220p' RISKS.md
git checkout -b 038-inline-password-flag-hardening
gofmt -w internal/sanitize/secrets.go internal/sanitize/matcher.go internal/sanitize/redact.go internal/sanitize/secrets_test.go internal/sanitize/redact_test.go internal/sanitize/preview_test.go
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/sanitize
gofmt -w internal/sanitize/matcher.go internal/sanitize/secrets_test.go internal/sanitize/apply_test.go
gofmt -w internal/sanitize/secrets_test.go internal/sanitize/redact_test.go
gofmt -w internal/cli/clean_test.go
env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...
git status --short
```

## Decisions

- Treat inline password detection as structured heuristic matching rather than a broad regex.
- Keep password-bearing flag names visible in redacted output.
- Limit high-entropy quarantine to sensitive key/value contexts instead of arbitrary mixed-case tokens.
- Complete the `039-high-entropy-token-false-positive-guards` follow-up inside this same slice because the work stayed local to the sanitizer heuristics.

## Assumptions

- Field-based tokenization is sufficient for the current hardening slice and can be revisited later if shell-escaped edge cases become important.

## Known issues

- The current detector still relies on simple field splitting rather than shell-specific parsing.

## Risks reduced

- Routine admin commands using `-p` no longer look like inline password leaks by default.
- Package-name, path, and device-label commands are much less likely to be quarantined as high-entropy secrets.
- Redacted output is easier to review because the command shape remains intact.

## Next recommended session

`032-systemd-user-service`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
