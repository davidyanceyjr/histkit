# 054-readme-p1-tests

## Summary

- Added the remaining P1 README-contract tests for config-loading failures across `scan`, `clean`, and `restore`.
- Added explicit parity coverage for `histkit clean` and `histkit clean --dry-run`.
- Added `clean --apply --shell zsh` coverage to ensure mutation, backup creation, and audit logging stay restricted to the requested shell source.
- Fixed a GitHub Actions portability issue in the `contrib` bash binding tests by accepting equivalent `bind -X` output variants.

## Objective completed or not completed

- Completed.

## Files read

- `AGENTS.md`: session workflow and closeout requirements
- `SESSION.md`: previous session state and update structure
- `ROADMAP.md`: roadmap boundaries for scan, clean, and restore
- `SKILLS/testing.md`: test expectations
- `SKILLS/config.md`: config constraints
- `contrib/wrappers_test.go`: existing bash binding assertions that failed under GitHub Actions
- `contrib/histkit.bash`: wrapper bind helper output and registration format
- `internal/cli/scan.go`: config-loading path for `scan`
- `internal/cli/clean.go`: planning/apply branching and shell filtering
- `internal/cli/restore.go`: config-loading path for `restore`
- `internal/cli/scan_test.go`: baseline scan coverage
- `internal/cli/clean_test.go`: baseline clean coverage
- `internal/cli/restore_test.go`: baseline restore coverage

## Files changed

- `contrib/wrappers_test.go`: relaxed bash bind-output assertions to accept both space-separated and colon-separated `bind -X` formats
- `internal/cli/scan_test.go`: added missing-config and invalid-TOML coverage for `scan`
- `internal/cli/clean_test.go`: added missing-config, invalid-TOML, planning-mode parity, and zsh apply filtering coverage
- `internal/cli/restore_test.go`: added missing-config and invalid-TOML coverage for `restore`
- `SESSION.md`: recorded the CI portability follow-up and current closeout state
- `SESSIONS/054-readme-p1-tests.md`: updated this session note

## Tests added

- `TestExecuteScanRejectsMissingConfigPath`
- `TestExecuteScanRejectsInvalidConfigTOML`
- `TestExecuteCleanRejectsMissingConfigPath`
- `TestExecuteCleanRejectsInvalidConfigTOML`
- `TestExecuteCleanDryRunFlagMatchesDefaultPlanningMode`
- `TestExecuteCleanApplyShellFlagFiltersToZshOnly`
- `TestExecuteRestoreRejectsMissingConfigPath`
- `TestExecuteRestoreRejectsInvalidConfigTOML`

## Tests run

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./contrib`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

## Known failures

- None.
- PR `#50` still needs this follow-up commit pushed and CI re-run before human review can resume.

## Commands run

- `sed -n '1,260p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '1,220p' SKILLS/config.md`
- `git status --short --branch`
- `ls -1 SESSIONS | sort | tail -n 8`
- `git checkout -b 054-readme-p1-tests`
- `git remote -v`
- `gh pr status`
- `gh auth status`
- `python /home/opsman/.codex/plugins/cache/openai-curated/github/1141b764/skills/gh-fix-ci/scripts/inspect_pr_checks.py --repo . --pr 50`
- `env XDG_CACHE_HOME=/tmp/codex-gh-cache gh pr checks 50`
- `env XDG_CACHE_HOME=/tmp/codex-gh-cache gh run view 25761451321 --json name,workflowName,conclusion,status,url,event,headBranch,headSha,jobs`
- `env XDG_CACHE_HOME=/tmp/codex-gh-cache gh run view 25761451321 --log`
- `sed -n '1,260p' internal/cli/scan_test.go`
- `sed -n '1,320p' internal/cli/clean_test.go`
- `sed -n '1,280p' internal/cli/restore_test.go`
- `sed -n '1,220p' internal/cli/scan.go`
- `sed -n '1,260p' internal/cli/clean.go`
- `sed -n '1,240p' internal/cli/restore.go`
- `sed -n '1,220p' contrib/wrappers_test.go`
- `sed -n '1,220p' contrib/histkit.bash`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./contrib`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

## Decisions made

- Keep the remaining P1 README-contract tests at the command layer in `internal/cli/*_test.go`.
- Use exact stdout equality to protect the equivalence of bare `clean` and `clean --dry-run`.
- Verify zsh apply filtering through history mutation scope, backup scope, and audit scope together.
- Accept both space-separated and colon-separated `bind -X` output in `contrib` tests because both formats describe the same registered binding.

## Assumptions made

- `NON-BLOCKING`: exact-output parity is the correct safety check because the README presents `clean` and `clean --dry-run` as the same mode.

## Unresolved questions

- None.

## Risks introduced or reduced

- Reduced risk that broken config paths or invalid TOML behave inconsistently across documented commands.
- Reduced risk that `clean --dry-run` drifts from the default planning mode.
- Reduced risk that shell filtering during apply accidentally mutates non-targeted history files or creates incorrect backup/audit artifacts.
- Reduced risk that GitHub Actions bash version or runner output formatting causes false-negative wrapper test failures.

## Next slice recommendation

- Push the CI portability fix for PR `#50`, wait for checks to return green, then resume the normal review/merge cleanup flow. After that, broader multi-source apply coverage would be the next logical increment if needed.
