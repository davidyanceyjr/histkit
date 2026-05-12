# 048-module-path-github

## Summary

- Updated the repository module path from `histkit` to `github.com/davidyanceyjr/histkit`.
- Rewrote all in-repo Go self-imports to use the GitHub-qualified module path.
- Verified the full repository test suite passes after the module identity change.

## Objective completed or not completed

- Completed.

## Files read

- `AGENTS.md`: session workflow and recording requirements
- `SESSION.md`: active session context and carry-forward structure
- `ROADMAP.md`: scope boundaries for the slice
- `SKILLS/go-cli.md`: Go CLI implementation constraints
- `SKILLS/testing.md`: verification expectations
- `go.mod`: existing module declaration

## Files changed

- `go.mod`: changed the module declaration to `github.com/davidyanceyjr/histkit`
- `cmd/histkit/main.go`: updated the CLI package import to the GitHub-qualified module path
- `internal/cli/pick.go`, `internal/cli/restore.go`, `internal/cli/restore_test.go`, `internal/cli/clean.go`, `internal/cli/clean_test.go`, `internal/cli/stats.go`, `internal/cli/scan.go`, `internal/cli/scan_test.go`, `internal/cli/doctor.go`, `internal/cli/pick_test.go`: rewrote self-imports to the new module path
- `internal/picker/candidates.go`, `internal/picker/candidates_test.go`: rewrote self-imports to the new module path
- `internal/doctor/checks.go`: rewrote self-imports to the new module path
- `internal/index/picker.go`, `internal/index/picker_test.go`, `internal/index/writer.go`, `internal/index/writer_test.go`: rewrote self-imports to the new module path
- `internal/audit/model.go`, `internal/audit/log_test.go`, `internal/audit/model_test.go`: rewrote self-imports to the new module path
- `internal/sanitize/apply.go`, `internal/sanitize/apply_test.go`, `internal/sanitize/matcher.go`, `internal/sanitize/matcher_test.go`, `internal/sanitize/preview.go`, `internal/sanitize/preview_test.go`, `internal/sanitize/quarantine.go`, `internal/sanitize/quarantine_test.go`, `internal/sanitize/secrets.go`, `internal/sanitize/secrets_test.go`, `internal/sanitize/trivial.go`, `internal/sanitize/trivial_test.go`: rewrote self-imports to the new module path
- `SESSION.md`: recorded the completed session state

## Tests added

- None.

## Tests run

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

## Known failures

- None.

## Commands run

- `git status --short --branch`: confirmed the starting worktree state
- `git branch --show-current`: confirmed the starting branch
- `git checkout -b 048-module-path-github`: created the implementation branch
- `sed -n '1,220p' SESSION.md`: reviewed current session context
- `sed -n '1,220p' ROADMAP.md`: reviewed roadmap boundaries
- `sed -n '1,220p' SKILLS/go-cli.md`: reviewed Go CLI constraints
- `sed -n '1,220p' SKILLS/testing.md`: reviewed test expectations
- `sed -n '1,120p' go.mod`: reviewed the existing module declaration
- `rg -n '"histkit(/|"$)' -g'*.go' -g'go.mod' -g'*.md'`: identified self-import usage
- `python - <<'PY' ... PY`: performed the mechanical module-path replacement across `go.mod`, `cmd/`, and `internal/`
- `gofmt -w $(rg -l 'github.com/davidyanceyjr/histkit/' cmd internal)`: formatted the touched Go files
- `rg -n 'module histkit|"histkit/' go.mod cmd internal README.md docs contrib SESSION.md SESSIONS DECISIONS.md RISKS.md`: verified no stale self-import or module declaration remained in code
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`: verified the full repository test suite

## Decisions made

- Keep the slice narrowly focused on module identity and import resolution.
- Preserve the `histkit` binary name and user-facing CLI references; only the Go module path changes in this slice.

## Assumptions made

- `NON-BLOCKING`: documentation references to `histkit` as the executable do not need changes because the deliverable is the module path, not a binary rename.

## Unresolved questions

- None.

## Risks introduced or reduced

- Reduced: the module now uses a GitHub-compatible path that matches the repository remote and supports canonical Go imports.
- Reduced: the repo no longer depends on the placeholder `histkit/...` self-import path.

## Next slice recommendation

- Optional documentation follow-up if you want explicit `go install github.com/davidyanceyjr/histkit/cmd/histkit@latest` guidance in `README.md`.
