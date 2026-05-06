# 035-shell-wrapper-polish

Status: completed

## Summary

Polished the contributed shell wrappers by keeping `Ctrl-R` as the default binding, allowing alternate key sequences without editing the scripts, and improving Zsh widget handoff before launching the external picker.

## Objective completed or not completed

Completed.

## Files read

- `AGENTS.md` - session workflow and safety constraints.
- `SESSION.md` - prior session handoff and required working-state structure.
- `ROADMAP.md` - confirmed `035-shell-wrapper-polish` as the next slice.
- `SKILLS/go-cli.md` - command-surface constraints and deterministic behavior expectations.
- `SKILLS/testing.md` - deterministic test guidance.
- `SESSION_PROMPT.md` - session closeout requirements.
- `docs/HUMAN_GATES.md` - open-question recording rules.
- `docs/histkit-implementation-plan.md` - wrapper behavior and release-readiness expectations.
- `README.md` - current shell integration documentation.
- `DECISIONS.md` - confirmed shell integration stays outside the Go binary.
- `RISKS.md` - checked for shell-side risk context.
- `SESSIONS/017-shell-wrappers.md` - earlier wrapper decisions and known gaps.
- `SESSIONS/032-systemd-user-service.md` - recent session artifact format and command tracking.
- `SESSIONS/034-doctor-systemd-checks.md` - latest session handoff.
- `internal/cli/pick.go` - confirmed the picker contract stays stdout-based.
- `contrib/histkit.bash` - existing Bash wrapper behavior.
- `contrib/histkit.zsh` - existing Zsh wrapper behavior.
- `contrib/wrappers_test.go` - existing wrapper test coverage.

## Files changed

- `contrib/histkit.bash` - added an optional key-sequence parameter to the Bash binding helper while preserving the default binding.
- `contrib/histkit.zsh` - added Zsh display invalidation before picker launch and an optional bind-key parameter.
- `contrib/wrappers_test.go` - added Bash tests for default and custom binding helpers and expanded Zsh script assertions.
- `README.md` - documented alternate binding examples for Bash and Zsh while preserving the default `Ctrl-R` flow.
- `SESSION.md` - updated the active working state and recorded the completed session details.
- `SESSIONS/035-shell-wrapper-polish.md` - recorded this session artifact.

## Tests added

- `TestBashBindHelperDefaultsToControlR`
- `TestBashBindHelperAcceptsCustomKeySequence`

## Tests run

- `bash -n contrib/histkit.bash`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./contrib`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

## Known failures

- No test failures.
- `zsh` is not installed in this environment, so Zsh runtime behavior remains covered by script-structure assertions rather than live widget execution.

## Commands run

- `pwd`
- `rg --files`
- `git status --short --branch`
- `sed -n '1,240p' SESSION.md`
- `sed -n '1,260p' ROADMAP.md`
- `git checkout -b 035-shell-wrapper-polish`
- `sed -n '1,220p' SKILLS/go-cli.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '1,240p' contrib/histkit.bash`
- `sed -n '1,240p' contrib/histkit.zsh`
- `sed -n '1,240p' contrib/wrappers_test.go`
- `rg -n "035-shell-wrapper-polish|shell wrapper|wrapper|histkit.bash|histkit.zsh|pick" -S README.md docs ROADMAP.md SESSION_PROMPT.md SESSIONS`
- `sed -n '665,710p' docs/histkit-implementation-plan.md`
- `sed -n '860,885p' docs/histkit-implementation-plan.md`
- `sed -n '995,1020p' docs/histkit-implementation-plan.md`
- `sed -n '330,370p' README.md`
- `sed -n '1,240p' internal/cli/pick.go`
- `sed -n '1,220p' SESSION_PROMPT.md`
- `sed -n '1,220p' SESSIONS/032-systemd-user-service.md`
- `sed -n '1,220p' docs/HUMAN_GATES.md`
- `rg -n "wrapper.*polish|shell wrapper.*polish|Ctrl-R|bind|READLINE|BUFFER|zle|replace the current shell editing buffer" -S README.md docs SESSIONS internal contrib`
- `sed -n '1,220p' DECISIONS.md`
- `sed -n '1,220p' RISKS.md`
- `sed -n '1,220p' SESSIONS/017-shell-wrappers.md`
- `bash -n contrib/histkit.bash`
- `command -v zsh >/dev/null 2>&1 && zsh -n contrib/histkit.zsh || echo 'zsh-unavailable'`
- `rg -n "TODO|FIXME|wrapper|shell integration|histkit_bind|__histkit_pick_bash|histkit_pick_zsh" -S contrib README.md internal docs`
- `sed -n '1,240p' contrib/systemd_units_test.go`
- `sed -n '1,220p' contrib/histkit-scan.service`
- `sed -n '1,220p' contrib/histkit-scan.timer`
- `bash -lc 'f(){ local keyseq="${1:-\\C-r}"; bind -x "\\\"$keyseq\\\":\\\"__histkit_pick_bash\\\""; bind -X; }; __histkit_pick_bash(){ :; }; f; f "\\C-x\\C-r"'`
- `bash -lc 'f(){ local keyseq="${1:-\\C-r}"; bind -x "\\\"$keyseq\\\":\\\"__histkit_pick_bash\\\""; bind -X; }; __histkit_pick_bash(){ :; }; f "\\C-x\\C-r"'`
- `gofmt -w contrib/wrappers_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./contrib`
- `git diff -- contrib/histkit.bash contrib/histkit.zsh contrib/wrappers_test.go README.md`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`
- `git status --short --branch`
- `git remote -v`

## Decisions made

- Keep wrapper polish within the shell scripts instead of changing the `histkit pick` binary contract.
- Preserve `Ctrl-R` as the default binding while allowing alternate key sequences through helper arguments.
- Add the Zsh display invalidation step before launching the picker to improve widget behavior with an external full-screen tool.

## Assumptions made

- `NON-BLOCKING`: Allowing an optional key-sequence argument is safe because it only changes wrapper ergonomics and is easy to reverse without affecting stored data, audit behavior, or the public CLI output contract.

## Unresolved questions

- No active blocking questions.
- No active non-blocking questions.

## Risks introduced or reduced

- Reduced: users can avoid a `Ctrl-R` conflict without editing the contributed wrapper files manually.
- Reduced: the Zsh wrapper now better prepares the terminal for the external picker flow.
- Remaining: live Zsh widget behavior is still unverified in this environment because `zsh` is unavailable.

## Next slice recommendation

- `036-readme-usage-flow`
