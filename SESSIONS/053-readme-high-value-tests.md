# 053-readme-high-value-tests

## Summary

- Split the README-derived high-value test plan into three parts and implemented them in sequence.
- Added command-level safety and config-path tests for `scan`, `clean`, and `restore`.
- Added sparse-config default-preservation coverage and sanitizer action-precedence coverage.

## Objective completed or not completed

- Completed.

## Files read

- `AGENTS.md`: session workflow and closeout requirements
- `SESSION.md`: previous session state and update structure
- `ROADMAP.md`: roadmap boundaries for scan, clean, and restore
- `SKILLS/testing.md`: test expectations
- `SKILLS/config.md`: config constraints and required tests
- `SKILLS/backup-restore.md`: restore safety boundaries and required tests
- `README.md`: behavior contract used to derive missing test coverage
- `internal/cli/scan.go`: `scan` config-loading behavior
- `internal/cli/clean.go`: `clean` preview/apply behavior and no-match handling
- `internal/cli/restore.go`: restore listing and backup-ID behavior
- `internal/config/config.go`: default loading and `~` expansion behavior
- `internal/sanitize/apply.go`: action-precedence implementation
- `internal/sanitize/secrets.go`: secret-rule actions
- `internal/sanitize/trivial.go`: trivial-rule actions
- existing tests in `internal/cli`, `internal/config`, and `internal/sanitize`: baseline coverage

## Files changed

- `internal/cli/clean_test.go`: added mutual-exclusion, no-match apply, and `~` config-path tests
- `internal/cli/scan_test.go`: added `~` config-path coverage for `scan`
- `internal/cli/restore_test.go`: added `~` config-path and missing-backup non-mutation tests
- `internal/config/config_test.go`: added sparse-config default-preservation coverage
- `internal/sanitize/apply_test.go`: added security-over-delete precedence coverage
- `SESSION.md`: recorded the completed session state
- `SESSIONS/053-readme-high-value-tests.md`: added this session note

## Tests added

- `TestExecuteCleanRejectsApplyAndDryRunTogether`
- `TestExecuteCleanApplyNoMatchesDoesNotCreateBackupOrAudit`
- `TestExecuteCleanConfigPathExpandsTilde`
- `TestExecuteScanConfigPathExpandsTilde`
- `TestExecuteRestoreConfigPathExpandsTilde`
- `TestExecuteRestoreMissingBackupIDLeavesHistoryUntouched`
- `TestLoadFromPathPreservesDefaultsForOmittedFields`
- `TestFinalActionPrefersSaferSecurityOutcomeOverDelete`

## Tests run

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/config`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/sanitize`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

## Known failures

- None.

## Commands run

- `sed -n '1,240p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `git status --short --branch`
- `ls -1 SESSIONS | sort | tail -n 8`
- `git checkout -b 053-readme-high-value-tests`
- `sed -n '1,220p' SKILLS/config.md`
- `sed -n '1,220p' SKILLS/backup-restore.md`
- `sed -n '1,240p' internal/backup/store.go`
- `sed -n '1,240p' internal/cli/root.go`
- `sed -n '1,240p' internal/sanitize/trivial.go`
- `sed -n '1,260p' internal/sanitize/secrets.go`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/config`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/sanitize`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

## Decisions made

- Split the work into three parts: CLI behavior, config defaults, and sanitizer precedence.
- Cover README `--config` behavior at the command layer where the user contract lives.
- Test safety precedence through `finalAction` directly because the current built-in rules do not naturally overlap in a single apply input.

## Assumptions made

- `NON-BLOCKING`: direct `finalAction` tests are sufficient coverage for the documented safety ordering until overlapping built-in rules exist.

## Unresolved questions

- None.

## Risks introduced or reduced

- Reduced risk of regressions in `clean --apply` safety boundaries and no-op apply behavior.
- Reduced risk of `--config` path-expansion regressions across documented commands.
- Reduced risk that future sanitizer changes accidentally prefer deletion over safer security actions.
- Residual risk remains that end-to-end precedence is not exercised by current built-in rule overlap.

## Next slice recommendation

- Add the remaining P1 README-contract tests for invalid config-file failures and explicit `clean --dry-run` parity with bare `clean`.
