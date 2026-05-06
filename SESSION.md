# SESSION.md

## Current session

ID: `033-systemd-user-timer`

Status: completed

## Objective

Add the initial optional `systemd --user` timer template for scheduled `histkit scan` runs and align the documentation with the shipped timer-and-service automation contract.

## Scope

Implement:

- a contributed `systemd --user` timer unit under `contrib/`
- deterministic tests that validate the shipped timer content without requiring a live `systemd` runtime
- documentation updates that reference the shipped timer-and-service template pair consistently

## Out of scope

- `doctor` systemd checks
- automation for `clean --apply`
- installer or enablement commands for user units
- changes to the scan service command itself

## Relevant skills

- `SKILLS/systemd-user.md`
- `SKILLS/testing.md`

## Acceptance criteria

- the repository ships a `systemd --user` timer template for scheduled `histkit scan` runs
- the timer content is covered by automated tests that do not require `systemd`
- docs and examples describe the shipped timer-and-service template pair consistently
- `go test ./...` passes

## Current repo state

Milestone 5 remains in progress. The repository now ships both contributed scan automation templates under `contrib/`, and the docs describe the timer and service pair consistently. The next roadmap slice is `034-doctor-systemd-checks`.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target
- Default automation runs `scan`, not destructive apply

## Risks to watch

- Shipping timer and service examples that diverge from the contributed artifacts will confuse users.
- Overreaching into `doctor` checks in this slice would blur the milestone boundary with `034`.

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

- Added `contrib/histkit-scan.timer` as the initial optional `systemd --user` timer template for scheduled `histkit scan` runs.
- Extended deterministic repository tests to lock the shipped timer content alongside the existing service template.
- Updated automation docs so they describe the shipped timer-and-service pair consistently.

Files changed:

- README.md
- SESSION.md
- SESSIONS/033-systemd-user-timer.md
- contrib/histkit-scan.timer
- contrib/systemd_units_test.go
- docs/histkit-implementation-plan.md

Files read:

- AGENTS.md
- SESSION.md
- ROADMAP.md
- SKILLS/systemd-user.md
- SKILLS/testing.md
- README.md
- docs/histkit-implementation-plan.md
- contrib/systemd_units_test.go
- contrib/histkit-scan.service
- SESSIONS/032-systemd-user-service.md

Tests added:

- exact-content coverage for `contrib/histkit-scan.timer`

Tests run:

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./contrib`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- Ship the timer as a contributed `systemd --user` template under `contrib/` with the interval contract from `SKILLS/systemd-user.md`.
- Keep timer automation scan-only by pairing it with the existing `histkit-scan.service` template.
- Keep `doctor` integration deferred to roadmap slice `034`.

Commands run:

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

Assumptions made:

- The timer interval contract already documented in `SKILLS/systemd-user.md` is sufficient for this slice, so no human decision is required before shipping the contributed timer template.

Risks introduced or reduced:

- Reduced: the repository now ships the timer artifact that the docs already referenced, eliminating drift between prose and contributed automation files.
- Reduced: timer content is now covered by deterministic tests instead of only inline documentation.
- Remaining: both templates still assume `histkit` is installed at `%h/.local/bin/histkit` and that the documented config file exists before the user enables automation.

Next recommended session:

- `034-doctor-systemd-checks`
