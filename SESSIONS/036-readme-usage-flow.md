# 036-readme-usage-flow

Status: completed

## Summary

Tightened the README so the documented usage flow is conservative and matches the current binary: `doctor`, `scan`, optional `stats` or `pick`, `clean --dry-run`, explicit `clean --apply`, and `restore`.

## Objective completed or not completed

Completed.

## Files read

- `AGENTS.md` - session workflow and safety constraints.
- `SESSION.md` - active session contract and required handoff structure.
- `ROADMAP.md` - confirmed `036-readme-usage-flow` as the next slice.
- `SKILLS/testing.md` - required verification workflow.
- `README.md` - target documentation and current drift.
- `internal/cli/root.go` - current top-level command surface.
- `internal/cli/scan.go` - implemented scan flags and behavior.
- `internal/cli/clean.go` - implemented clean flags, dry-run/apply behavior, and backup requirement.
- `internal/cli/pick.go` - picker contract and config support.
- `internal/cli/doctor.go` - doctor flags and output contract.
- `internal/cli/stats.go` - actual stats output scope.
- `internal/cli/restore.go` - restore listing and apply behavior.
- `internal/config/config.go` - actual config schema and default paths.
- `configs/config.example.toml` - example config currently shipped in the repo.
- `internal/doctor/checks.go` - current doctor checks.
- `internal/sanitize/apply.go` - actual quarantine/apply behavior.
- `internal/sanitize/quarantine.go` - quarantine record model for terminology validation.
- `internal/audit/log.go` - audit file behavior.
- `internal/backup/model.go` - backup ID format.
- `cmd/histkit/main.go` - actual exit behavior.
- `contrib/histkit-scan.service` - shipped automation command.
- `contrib/histkit-scan.timer` - shipped automation schedule.
- `contrib/histkit.bash` - wrapper behavior for picker docs.
- `contrib/histkit.zsh` - wrapper behavior for picker docs.
- `docs/histkit-implementation-plan.md` - intended safe workflow and automation posture.
- `DECISIONS.md` - durable project decisions around conservative behavior and automation.
- `RISKS.md` - existing risk framing for trust and history mutation.
- `SESSIONS/035-shell-wrapper-polish.md` - prior session artifact format and latest merged context.
- `SESSIONS/034-doctor-systemd-checks.md` - prior closeout pattern for recent docs-related work.
- `SESSION_PROMPT.md` - closeout requirements.

## Files changed

- `README.md` - aligned the documented workflow, command surface, flags, config example, examples, wrapper notes, and `systemd --user` automation section with the current implementation.
- `SESSION.md` - updated the working state for review/merge handoff after this documentation slice.
- `SESSIONS/036-readme-usage-flow.md` - recorded this session artifact.

## Tests added

- None. This was a documentation-only slice.

## Tests run

- `env GOCACHE=/home/opsman/project_git/histkit/.gocache GOMODCACHE=/home/opsman/project_git/histkit/.gomodcache GOPATH=/home/opsman/project_git/histkit/.gopath go test ./...`

## Known failures

- No repository test failures.
- The first `go test ./...` attempt failed before execution because the default Go build and module cache locations under `/home/opsman` were read-only in this environment; rerunning with repo-local caches succeeded.

## Commands run

- `git status --short --branch` - inspected repo state; branch started on `main` and worktree was clean.
- `sed -n '1,220p' SESSION.md` - loaded the active session contract.
- `sed -n '1,260p' ROADMAP.md` - confirmed the roadmap slice.
- `sed -n '1,220p' SKILLS/testing.md` - loaded the only required local skill.
- `git checkout -b 036-readme-usage-flow` - created the session branch.
- `sed -n ... README.md` - inspected the current documentation sections.
- `sed -n ... internal/cli/*.go` - validated current command surface and flags.
- `sed -n ... internal/config/config.go` - validated config schema and default paths.
- `sed -n ... internal/doctor/checks.go` - validated doctor checks.
- `sed -n ... internal/sanitize/apply.go` - validated apply and quarantine behavior.
- `sed -n ... cmd/histkit/main.go` - validated exit-code behavior.
- `sed -n ... contrib/histkit-scan.service contrib/histkit-scan.timer contrib/histkit.bash contrib/histkit.zsh` - validated automation and wrapper docs.
- `go run ./cmd/histkit help` - failed because the default Go cache directory was read-only.
- `go run ./cmd/histkit help clean` - failed for the same cache reason.
- `rg -n -e ... README.md` - verified removal of unsupported command and flag references.
- `git diff -- README.md` - reviewed the documentation changes.
- `mkdir -p .gocache .gomodcache && env GOCACHE=... GOMODCACHE=... GOPATH=... go test ./...` - ran the full test suite successfully with repo-local caches.
- `chmod -R u+w .gocache .gomodcache .gopath && rm -rf .gocache .gomodcache .gopath` - cleaned the temporary Go cache directories after verification.

## Decisions made

- The README should document only the commands and flags that exist in the current binary.
- The recommended README workflow should start with `doctor`, keep `scan` and `clean` separated, and present `clean --apply` as an explicit reviewed step.
- Automation documentation should stay limited to the shipped scan-only `systemd --user` units and must not imply unattended destructive cleanup.
- Snippet and quarantine behavior should be described conservatively without claiming unimplemented management commands.

## Assumptions made

- `NON-BLOCKING`: README examples may use a representative backup ID format from the current backup model because that format is already implemented and changing the example later is low-cost.

## Unresolved questions

- No active blocking questions.
- No active non-blocking questions.

## Risks introduced or reduced

- Reduced: the README no longer advertises unsupported commands or flags.
- Reduced: the end-to-end workflow now better reflects the project’s conservative apply and restore posture.
- Reduced: automation guidance now clearly states that the shipped timer runs `scan`, not `clean --apply`.

## Next slice recommendation

- `037-release-readiness-pass` after this README slice is reviewed and merged.
