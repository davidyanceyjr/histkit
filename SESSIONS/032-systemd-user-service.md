# 032-systemd-user-service

Status: completed

## Summary

Added the initial `systemd --user` scan service template under `contrib/`, covered it with a deterministic repository test, and aligned the automation documentation with the project decision that scheduled automation defaults to `histkit scan` rather than cleanup.

## Objective completed or not completed

Completed.

## Files read

- `AGENTS.md` - session workflow and safety constraints.
- `SESSION.md` - active work contract and prior session handoff.
- `ROADMAP.md` - confirmed `032-systemd-user-service` as the next slice.
- `SKILLS/systemd-user.md` - unit placement and automation constraints.
- `SKILLS/testing.md` - test discipline for non-runtime unit verification.
- `README.md` - existing automation and file-path documentation.
- `docs/histkit-implementation-plan.md` - implementation-plan filenames and service examples.
- `DECISIONS.md` - confirmed the durable `scan`-first automation decision.
- `SESSIONS/031-failure-recovery-tests.md` - prior session handoff.
- `contrib/wrappers_test.go` - existing `contrib` test package conventions.
- `contrib/histkit.bash` - existing contributed artifact layout.
- `contrib/histkit.zsh` - existing contributed artifact layout.
- `internal/doctor/checks.go` - confirmed systemd checks are not yet part of `doctor`.
- `internal/config/config.go` - checked the default config path used by examples.

## Files changed

- `contrib/histkit-scan.service` - added the initial oneshot `systemd --user` service template for scheduled scans.
- `contrib/systemd_units_test.go` - added an exact-content regression test for the service template.
- `README.md` - updated automation filenames and examples to consistently describe scheduled scans.
- `docs/histkit-implementation-plan.md` - aligned contributed filenames and timer description with scan-oriented automation.
- `SESSION.md` - updated the active working state and recorded the completed session details.
- `SESSIONS/032-systemd-user-service.md` - recorded this session artifact.

## Tests added

- `TestSystemdScanServiceTemplate` in `contrib/systemd_units_test.go`.

## Tests run

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./contrib`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

## Known failures

- No test failures.
- `systemd-analyze verify contrib/histkit-scan.service` reported that `/home/opsman/.local/bin/histkit` is not installed in this workspace. The unit syntax was otherwise accepted, and the missing executable is expected for an uninstalled template artifact.

## Commands run

- `pwd`
- `rg --files -g 'SESSION.md' -g 'ROADMAP.md' -g 'SKILLS/**' -g 'AGENTS.md'`
- `git status --short --branch`
- `sed -n '1,240p' SESSION.md`
- `sed -n '1,260p' ROADMAP.md`
- `sed -n '1,240p' SKILLS/systemd-user.md`
- `rg -n "systemd|timer|service|doctor|scan" .`
- `rg --files SESSIONS docs cmd internal | rg 'systemd|doctor|scan|config|paths'`
- `sed -n '1,240p' internal/doctor/checks.go`
- `sed -n '520,575p' README.md`
- `sed -n '728,766p' docs/histkit-implementation-plan.md`
- `sed -n '1,220p' internal/config/config.go`
- `ls -1`
- `find . -maxdepth 2 -type d | sort`
- `sed -n '1,220p' SESSIONS/031-failure-recovery-tests.md`
- `git branch --all --list '*032*'`
- `find contrib -maxdepth 3 -type f | sort`
- `git log --oneline -5`
- `git checkout -b 032-systemd-user-service`
- `sed -n '1,220p' contrib/wrappers_test.go`
- `sed -n '1,220p' contrib/histkit.bash`
- `sed -n '1,220p' contrib/histkit.zsh`
- `sed -n '1,120p' DECISIONS.md`
- `rg -n "histkit-clean|cleanup periodically|Run histkit cleanup periodically|histkit scan --config|systemd integration|histkit.*service|histkit.*timer" README.md docs/histkit-implementation-plan.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '430,570p' README.md`
- `sed -n '140,180p' docs/histkit-implementation-plan.md`
- `gofmt -w contrib/systemd_units_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./contrib`
- `command -v systemd-analyze >/dev/null 2>&1 && systemd-analyze verify contrib/histkit-scan.service || true`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`
- `git status --short`
- `git diff -- README.md docs/histkit-implementation-plan.md contrib/histkit-scan.service contrib/systemd_units_test.go SESSION.md`
- `git remote -v`

## Decisions made

- Ship the first automation artifact as a contributed `systemd --user` service template only, leaving timer and `doctor` integration to their roadmap slices.
- Keep the service template’s `ExecStart` command scan-only and aligned with the explicit config path already used in project docs.
- Normalize automation naming and prose around `histkit-scan.*` artifacts so the docs match decision `006`.

## Assumptions made

- `NON-BLOCKING`: Using the existing explicit config path in the service template is safe for this slice because it matches the current documentation contract and can be revised later with low cost if the project prefers config-optional automation defaults.

## Unresolved questions

- No active blocking questions.
- No active non-blocking questions.

## Risks introduced or reduced

- Reduced: automation examples and artifact names now consistently describe scheduled scans instead of cleanup-by-default behavior.
- Reduced: the repository now ships a tested service template instead of only prose examples.
- Remaining: the template still assumes `histkit` is installed at `%h/.local/bin/histkit` and that the documented config file exists before the user enables automation.

## Next slice recommendation

- `033-systemd-user-timer`
