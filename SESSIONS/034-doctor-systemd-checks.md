# 034-doctor-systemd-checks

Status: completed

## Summary

Extended `histkit doctor` with an optional `systemd --user` automation check that reports unconfigured systems as healthy, warns on partial installs of the shipped timer-and-service pair, and stays deterministic without a live `systemd` runtime.

## Objective completed or not completed

Completed.

## Files read

- `AGENTS.md` - session workflow and safety constraints.
- `SESSION.md` - prior session handoff and required working-state structure.
- `ROADMAP.md` - confirmed `034-doctor-systemd-checks` as the next slice.
- `SKILLS/systemd-user.md` - systemd automation constraints and doctor expectations.
- `SKILLS/testing.md` - deterministic test guidance.
- `internal/doctor/checks.go` - existing doctor checks and report flow.
- `internal/cli/doctor.go` - doctor output formatting and invocation path.
- `internal/cli/doctor_test.go` - current CLI-level doctor expectations.
- `internal/config/config.go` - confirmed there is no separate automation config surface in this slice.
- `README.md` - current doctor documentation.
- `docs/histkit-implementation-plan.md` - roadmap-level doctor contract.

## Files changed

- `internal/doctor/checks.go` - added `systemd_user_units` checking for the default user-unit directory.
- `internal/doctor/checks_test.go` - added deterministic unit tests for not-configured, partial, and complete automation states.
- `internal/cli/doctor_test.go` - updated doctor CLI expectations and added coverage for partial installs.
- `README.md` - clarified that the doctor systemd check applies when histkit automation is installed.
- `SESSION.md` - updated the active working state and recorded the completed session details.
- `SESSIONS/034-doctor-systemd-checks.md` - recorded this session artifact.

## Tests added

- `TestCheckSystemdUserUnitsNotConfigured`
- `TestCheckSystemdUserUnitsWarnsForPartialInstall`
- `TestCheckSystemdUserUnitsReportsPresentPair`
- `TestExecuteDoctorWarnsForPartialSystemdInstall`

## Tests run

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/doctor ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

## Known failures

- No test failures.

## Commands run

- `git status --short --branch`
- `sed -n '1,260p' SESSION.md`
- `sed -n '1,260p' ROADMAP.md`
- `sed -n '1,240p' SKILLS/systemd-user.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '1,260p' internal/doctor/checks.go`
- `sed -n '1,260p' internal/config/config.go`
- `rg -n "doctor|systemd|timer|service|automation" README.md docs internal cmd`
- `rg --files internal cmd testdata | rg 'doctor|config'`
- `sed -n '1,260p' internal/cli/doctor_test.go`
- `sed -n '210,235p' docs/histkit-implementation-plan.md`
- `sed -n '100,140p' README.md`
- `sed -n '170,190p' docs/histkit-implementation-plan.md`
- `sed -n '1,220p' internal/cli/doctor.go`
- `sed -n '1,140p' internal/cli/root_test.go`
- `rg -n "systemd --user unit visibility|backup path availability|doctor --json" README.md docs/histkit-implementation-plan.md internal/cli`
- `git checkout -b 034-doctor-systemd-checks`
- `gofmt -w internal/doctor/checks.go internal/doctor/checks_test.go internal/cli/doctor_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/doctor ./internal/cli`
- `gofmt -w internal/doctor/checks_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/doctor ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`
- `git status --short`
- `git diff -- README.md internal/doctor/checks.go internal/doctor/checks_test.go internal/cli/doctor_test.go`

## Decisions made

- Treat absent histkit user units as `OK` because automation remains optional.
- Warn on partial installs of the expected `histkit-scan.service` and `histkit-scan.timer` pair.
- Keep the check limited to installed file visibility rather than `systemctl` runtime inspection.

## Assumptions made

- `NON-BLOCKING`: Filesystem presence in `~/.config/systemd/user` is a safe proxy for “configured” in this slice because the project does not yet expose a dedicated automation configuration surface, and revisiting that trigger later is low-cost.

## Unresolved questions

- No active blocking questions.
- No active non-blocking questions.

## Risks introduced or reduced

- Reduced: doctor now identifies incomplete histkit automation installs that would otherwise be easy to miss.
- Reduced: fresh systems without automation are no longer at risk of noisy warnings from an optional feature check.
- Remaining: the doctor check does not inspect enablement, timers.target wiring, or runtime service health.

## Next slice recommendation

- `035-shell-wrapper-polish`
