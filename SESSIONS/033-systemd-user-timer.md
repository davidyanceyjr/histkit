# 033-systemd-user-timer

Status: completed

## Summary

Added the initial `systemd --user` scan timer template under `contrib/`, covered it with deterministic repository tests, and aligned the automation docs with the shipped timer-and-service artifact pair.

## Objective completed or not completed

Completed.

## Files read

- `AGENTS.md` - session workflow and safety constraints.
- `SESSION.md` - prior session handoff and required working-state structure.
- `ROADMAP.md` - confirmed `033-systemd-user-timer` as the next slice.
- `SKILLS/systemd-user.md` - timer contract and automation constraints.
- `SKILLS/testing.md` - test discipline for non-runtime unit verification.
- `README.md` - current automation documentation and timer example wording.
- `docs/histkit-implementation-plan.md` - implementation-plan timer and service examples.
- `contrib/systemd_units_test.go` - existing `contrib` test conventions for unit-file content locks.
- `contrib/histkit-scan.service` - existing service template paired with the new timer.
- `SESSIONS/032-systemd-user-service.md` - previous session handoff and remaining risks.

## Files changed

- `contrib/histkit-scan.timer` - added the initial `systemd --user` timer template for scheduled scans.
- `contrib/systemd_units_test.go` - added exact-content regression coverage for the timer template.
- `README.md` - updated automation docs to describe the shipped timer-and-service pair and aligned the service description text.
- `docs/histkit-implementation-plan.md` - aligned the service example text with the shipped contributed template.
- `SESSION.md` - updated the active working state and recorded the completed session details.
- `SESSIONS/033-systemd-user-timer.md` - recorded this session artifact.

## Tests added

- `TestSystemdScanTimerTemplate` in `contrib/systemd_units_test.go`.

## Tests run

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./contrib`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

## Known failures

- No test failures.
- `systemd-analyze verify contrib/histkit-scan.service contrib/histkit-scan.timer` reported that `/home/opsman/.local/bin/histkit` is not installed in this workspace. The timer and service syntax were otherwise accepted, and the missing executable is expected for an uninstalled template artifact.

## Commands run

- `git branch --all --list '*033*'`
- `git remote -v`
- `git log --oneline -5`
- `git checkout -b 033-systemd-user-timer`
- `sed -n '1,220p' contrib/systemd_units_test.go`
- `sed -n '520,575p' README.md`
- `sed -n '728,766p' docs/histkit-implementation-plan.md`
- `sed -n '1,260p' SESSION.md`
- `gofmt -w contrib/systemd_units_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./contrib`
- `command -v systemd-analyze >/dev/null 2>&1 && systemd-analyze verify contrib/histkit-scan.service contrib/histkit-scan.timer || true`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

## Decisions made

- Ship the first automation timer artifact as a contributed `systemd --user` template under `contrib/`.
- Keep the timer bound to the existing scan-only service template rather than introducing any apply-mode automation.
- Align the service description text in docs with the shipped contributed service template to reduce drift.

## Assumptions made

- `NON-BLOCKING`: The timer interval contract from `SKILLS/systemd-user.md` is safe to ship as-is because it already matches the project documentation and can be revised later at low cost if scheduling defaults change.

## Unresolved questions

- No active blocking questions.
- No active non-blocking questions.

## Risks introduced or reduced

- Reduced: the repository now ships the timer artifact already described in docs, tightening the automation contract.
- Reduced: exact-content tests now guard both contributed `systemd --user` unit files against accidental drift.
- Remaining: the contributed units still assume the documented install path and config path exist before the user enables automation.

## Next slice recommendation

- `034-doctor-systemd-checks`
