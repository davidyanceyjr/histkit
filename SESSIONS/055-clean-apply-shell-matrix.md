# 055-clean-apply-shell-matrix

## Summary

- Added command-level `clean --apply --shell` tests for mixed bash/zsh source presence, selected-shell absence, and selected-shell no-match behavior.
- Verified that rewrite, backup, and audit side effects stay scoped to the selected shell only.
- Refreshed `SESSION.md` because the previous carry-forward state still described the already-merged PR `#50` follow-up.
- Pushed branch `055-clean-apply-shell-matrix` and opened draft PR `#51`.

## Objective completed or not completed

- Completed.

## Files read

- `AGENTS.md`: session workflow and closeout requirements
- `SESSION.md`: previous session state and stale carry-forward context
- `ROADMAP.md`: roadmap boundaries for the test slice
- `SKILLS/testing.md`: verification expectations
- `SKILLS/backup-restore.md`: backup/audit constraints for apply-path assertions
- `internal/cli/clean_test.go`: existing clean command coverage
- `internal/cli/clean.go`: `--shell` apply flow and output behavior
- `internal/history/detect.go`: current detector contract and shell filtering limits
- `internal/config/config.go`: default path layout
- `internal/audit/log.go`: audit append behavior
- `internal/audit/model.go`: audit line rendering and rule ordering
- `internal/backup/create.go`: backup layout
- `internal/sanitize/apply.go`: rewrite behavior used by the command path

## Files changed

- `internal/cli/clean_test.go`: added `TestCleanApplyShellMixedSources`, `TestCleanApplyShellNoMatchingSources`, `TestCleanApplyShellBackupScope`, and small local helpers
- `SESSION.md`: replaced stale prior-session state with the new working state
- `SESSIONS/055-clean-apply-shell-matrix.md`: recorded this session

## Tests added

- `TestCleanApplyShellMixedSources`
- `TestCleanApplyShellNoMatchingSources`
- `TestCleanApplyShellBackupScope`

## Tests run

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
- `git add SESSION.md internal/cli/clean_test.go SESSIONS/055-clean-apply-shell-matrix.md && git commit -m "Add clean apply shell matrix tests"`
- `git push -u origin 055-clean-apply-shell-matrix`
- `gh pr create --base main --head 055-clean-apply-shell-matrix --title "Add clean apply shell matrix tests" ...`: failed with a GitHub GraphQL no-commits/head-ref error despite the pushed branch
- GitHub connector `_create_pull_request`: opened draft PR `#51`

## Known failures

- None.

## Commands run

- `sed -n '1,240p' AGENTS.md`
- `sed -n '1,260p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '1,220p' SKILLS/backup-restore.md`
- `git status --short --branch`
- `git checkout -b 055-clean-apply-shell-matrix`
- `sed -n '1,360p' internal/cli/clean_test.go`
- `sed -n '1,320p' internal/cli/clean.go`
- `sed -n '360,520p' internal/cli/clean_test.go`
- `sed -n '1,260p' internal/history/detect.go`
- `sed -n '1,260p' internal/config/config.go`
- `sed -n '1,260p' internal/audit/log.go`
- `sed -n '1,260p' internal/audit/model.go`
- `sed -n '1,260p' internal/backup/create.go`
- `rg -n "backupMatches|AuditLog|clean apply:" internal/cli/*_test.go`
- `sed -n '1,260p' internal/sanitize/apply.go`
- `gofmt -w internal/cli/clean_test.go`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

## Decisions made

- Keep the matrix aligned with the current detector contract of one candidate source per shell.
- Assert backup and audit scoping at the command layer because that is the user-facing contract being protected.
- Use small local helper functions in `clean_test.go` to keep the new command tests concise.

## Assumptions made

- `NON-BLOCKING`: one source per shell is the intended current contract, so the new tests should not fabricate unsupported same-shell multi-source cases. This is safe because it preserves current behavior and has low reversal cost if source detection expands later.

## Unresolved questions

- None.

## Risks introduced or reduced

- Reduced: shell-filter regressions that would rewrite, back up, or audit the wrong shell source are now covered at the command level.
- Remaining: coverage still reflects the current detector limit of one source per shell; broader same-shell multi-source behavior is not testable until the detector contract changes.

## Next slice recommendation

- Review draft PR `#51`.
- After approval, merge and clean up the branch.
- If future work needs broader shell-source matrices, decide whether `history.DetectSources` should support multiple paths per shell before expanding tests.
