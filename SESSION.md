# SESSION.md

## Current session

ID: `032-systemd-user-service`

Status: completed

## Objective

Add the initial optional `systemd --user` service template for scheduled `histkit scan` runs and align the documentation with the scan-only automation contract.

## Scope

Implement:

- a contributed `systemd --user` oneshot service unit under `contrib/`
- deterministic tests that validate the shipped unit content without requiring a live `systemd` runtime
- documentation updates that reference scan-oriented automation artifacts consistently

## Out of scope

- timer unit implementation
- `doctor` systemd checks
- automation for `clean --apply`
- installer or enablement commands for user units

## Relevant skills

- `SKILLS/systemd-user.md`
- `SKILLS/testing.md`

## Acceptance criteria

- the repository ships a `systemd --user` service template for `histkit scan`
- the unit content is covered by automated tests that do not require `systemd`
- docs and examples describe scan-oriented automation rather than cleanup-by-default
- `go test ./...` passes

## Current repo state

Milestone 4 is complete and the roadmap recommends `032-systemd-user-service` next. The repository now ships a contributed scan service template and the automation docs consistently describe scheduled scans rather than cleanup-by-default behavior.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target
- Default automation runs `scan`, not destructive apply

## Risks to watch

- Shipping a service template that diverges from the documented automation contract will confuse users.
- Overreaching into timer or doctor behavior in this slice would blur the milestone boundaries.

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

- Added `contrib/histkit-scan.service` as the initial optional `systemd --user` oneshot template for scheduled `histkit scan` runs.
- Added a deterministic repository test that locks the shipped unit content without requiring a live `systemd` runtime.
- Updated automation docs and example filenames to consistently describe scan-oriented scheduling rather than cleanup-by-default behavior.

Files changed:

- README.md
- SESSION.md
- SESSIONS/032-systemd-user-service.md
- contrib/histkit-scan.service
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
- DECISIONS.md
- SESSIONS/031-failure-recovery-tests.md
- contrib/wrappers_test.go
- contrib/histkit.bash
- contrib/histkit.zsh
- internal/doctor/checks.go
- internal/config/config.go

Tests added:

- exact-content coverage for `contrib/histkit-scan.service`

Tests run:

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./contrib`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- Keep the initial contributed unit as a scan-only `systemd --user` service template under `contrib/`.
- Keep `doctor` and timer work deferred to roadmap slices `034` and `033` respectively.
- Keep the explicit config path in the service template to match the existing documented contract for automation examples.

Commands run:

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

Assumptions made:

- The existing automation examples establish a sufficient contract for the initial service template, so no human decision is required before implementing slice `032`.

Risks introduced or reduced:

- Reduced: automation docs and artifact names now consistently describe scheduled scans instead of cleanup-by-default behavior.
- Reduced: the repository now ships a tested scan service template instead of only prose examples.
- Remaining: the explicit `--config` path in the service template assumes users install automation after creating a config file at the documented location.

Next recommended session:

- `033-systemd-user-timer`
