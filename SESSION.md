# SESSION.md

## Current session

ID: `034-doctor-systemd-checks`

Status: completed

## Objective

Extend `histkit doctor` to verify optional `systemd --user` automation state without requiring a live `systemd` runtime or forcing automation on users who have not installed it.

## Scope

Implement:

- a `doctor` check for contributed `systemd --user` scan unit visibility in the default user-unit directory
- deterministic tests for unconfigured, partial-install, and fully present automation states
- documentation updates that describe the automation check as optional and install-sensitive

## Out of scope

- installer or enablement commands for user units
- runtime status checks against `systemctl`
- automation for `clean --apply`
- changes to the timer or service template contents

## Relevant skills

- `SKILLS/systemd-user.md`
- `SKILLS/testing.md`

## Acceptance criteria

- `histkit doctor` reports `OK` when no histkit user units are installed
- `histkit doctor` warns when only part of the expected timer-and-service pair is present
- deterministic tests cover the new systemd check without requiring a live `systemd` runtime
- `go test ./...` passes

## Current repo state

Milestone 5 remains in progress. `doctor` now inspects the default `systemd --user` unit directory for the shipped `histkit-scan.service` and `histkit-scan.timer` pair, while still treating automation as optional. The next roadmap slice is `035-shell-wrapper-polish`.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target
- Default automation runs `scan`, not destructive apply

## Risks to watch

- `doctor` must not warn on fresh systems that have not opted into automation.
- The automation check should stay filesystem-based; runtime or installer behavior belongs in later slices if needed.

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

- Added a `systemd_user_units` doctor check that treats missing histkit user units as unconfigured rather than unhealthy.
- Added deterministic tests for unconfigured, partial, and complete user-unit states.
- Updated the `doctor` documentation to describe the systemd check as relevant only when histkit automation is installed.

Files changed:

- README.md
- SESSION.md
- SESSIONS/034-doctor-systemd-checks.md
- internal/cli/doctor_test.go
- internal/doctor/checks.go
- internal/doctor/checks_test.go

Files read:

- AGENTS.md
- SESSION.md
- ROADMAP.md
- SKILLS/systemd-user.md
- SKILLS/testing.md
- internal/doctor/checks.go
- internal/cli/doctor.go
- internal/cli/doctor_test.go
- internal/config/config.go
- README.md
- docs/histkit-implementation-plan.md

Tests added:

- `TestCheckSystemdUserUnitsNotConfigured`
- `TestCheckSystemdUserUnitsWarnsForPartialInstall`
- `TestCheckSystemdUserUnitsReportsPresentPair`
- `TestExecuteDoctorWarnsForPartialSystemdInstall`

Tests run:

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/doctor ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- Treat absent histkit user units as `OK` because automation is optional.
- Warn only on partial installs of the expected timer-and-service pair.
- Keep the systemd check filesystem-based instead of adding runtime `systemctl` dependencies in this slice.

Commands run:

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

Assumptions made:

- The absence of a dedicated automation config surface makes filesystem presence in `~/.config/systemd/user` the safest proxy for whether the systemd check is relevant in this slice.

Risks introduced or reduced:

- Reduced: `doctor` now catches incomplete histkit automation installs without warning on fresh systems.
- Reduced: the systemd automation contract is now checked by tests instead of only by docs and manual review.
- Remaining: the check does not inspect unit enablement or runtime health, only the installed file pair.

Next recommended session:

- `035-shell-wrapper-polish`
