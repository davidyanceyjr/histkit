# SESSION.md

## Current session

ID: `038-inline-password-flag-hardening`

Status: completed

## Objective

Harden inline password detection and reduce false positives from high-entropy token matching.

## Scope

Implement:

- narrower inline password matching that keeps `--password`/`--passwd` coverage while avoiding broad short-flag matches
- redaction behavior that preserves the password-bearing flag shape while masking only the secret value
- high-entropy token guards that avoid quarantining common package, path, and device-label commands
- focused sanitizer regression tests for both true positives and false positives

## Out of scope

- new user-facing commands
- config-driven custom secret rules
- shell-aware command parsing beyond current tokenizer/field handling
- broader sanitizer config or scoring changes outside the built-in secret rules

## Relevant skills

- `SKILLS/sanitizer.md`
- `SKILLS/testing.md`

## Acceptance criteria

- commands like `mysql --password hunter2` and `mysql -phunter2` still match and redact safely
- commands like `ssh -p 2222 prod-box` do not match the inline password rule
- high-entropy matching still catches sensitive key/value forms such as exported tokens
- package, path, and device-label commands covered by tests do not trigger high-entropy quarantine
- `go test ./...` passes

## Current repo state

Milestone 3 secret rules already exist, but inline password detection is overly broad and the high-entropy heuristic still trusts generic mixed-case tokens too readily.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Over-narrowing the secret heuristics could miss real credentials that the current catalog catches.
- Over-broad matching will damage trust by flagging routine admin commands as secrets.

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

No non-blocking questions currently recorded.

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

No questions answered this session.

## End-of-session notes

Summary:

- Replaced the broad inline password regex with a structured heuristic that still catches `--password`, `--passwd`, and MySQL-style `-psecret` forms without flagging unrelated `-p` options.
- Kept password-bearing commands reviewable by redacting only the secret value instead of collapsing the whole argument.
- Narrowed high-entropy matching to sensitive key/value contexts so package names, filesystem paths, and device labels are not quarantined by the generic entropy rule.

Files changed:

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

Files read:

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

Tests added:

- inline password true-positive coverage for MySQL short `-psecret`
- false-positive guards for `ssh -p`, `grep -p`, package names, and device-label/path commands
- redaction coverage for preserving password flags while masking only values
- high-entropy flag-value redaction coverage for sensitive `--session-token` style arguments

Tests run:

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/sanitize`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- Treat inline password detection as a structured heuristic instead of a broad regex so short `-p` matching can stay command-aware.
- Preserve flag structure in redacted output for password-bearing commands to keep dry-run and rewritten history more reviewable.
- Restrict the built-in high-entropy heuristic to sensitive key/value contexts rather than arbitrary mixed-case tokens.
- Fold the false-positive guard work from slice `039-high-entropy-token-false-positive-guards` into this session because it shares the same sanitizer boundary.

Commands run:

- `git status --short --branch`
- `sed -n '1,220p' SKILLS/sanitizer.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '1,240p' ROADMAP.md`
- `rg -n "password|passwd|entropy|token|secret|credential|bearer|apikey|api_key|inline" internal/sanitize internal -g '!**/*_test.go'`
- `sed -n '1,260p' internal/sanitize/secrets.go`
- `sed -n '1,260p' internal/sanitize/matcher.go`
- `sed -n '1,260p' internal/sanitize/redact.go`
- `sed -n '1,260p' internal/sanitize/secrets_test.go`
- `sed -n '1,260p' internal/sanitize/redact_test.go`
- `sed -n '1,260p' internal/sanitize/preview_test.go`
- `sed -n '1,140p' internal/sanitize/apply_test.go`
- `sed -n '1,240p' internal/cli/clean_test.go`
- `sed -n '1,220p' SESSIONS/021-secret-rules.md`
- `sed -n '1,220p' docs/histkit-implementation-plan.md`
- `sed -n '1,220p' RISKS.md`
- `git checkout -b 038-inline-password-flag-hardening`
- `gofmt -w internal/sanitize/secrets.go internal/sanitize/matcher.go internal/sanitize/redact.go internal/sanitize/secrets_test.go internal/sanitize/redact_test.go internal/sanitize/preview_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/sanitize`
- `gofmt -w internal/sanitize/matcher.go internal/sanitize/secrets_test.go internal/sanitize/apply_test.go`
- `gofmt -w internal/sanitize/secrets_test.go internal/sanitize/redact_test.go`
- `gofmt -w internal/cli/clean_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`
- `git status --short`

Assumptions made:

- Field-based tokenization with `strings.Fields` is sufficient for the current hardening slice even though it is not a full shell parser.

Risks introduced or reduced:

- Reduced: routine `-p` flags and mixed-case path or package tokens are much less likely to trigger false secret matches.
- Reduced: preview and rewritten output now preserve password flag context instead of hiding the entire argument shape.
- Remaining: quoted or shell-escaped password values are still handled by the current field-based heuristic, not shell-specific parsing.

Next recommended session:

- `032-systemd-user-service`
