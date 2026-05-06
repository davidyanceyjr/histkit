# SESSION.md

## Current session

ID: `043-command-help-detail`

Status: completed

## Objective

Refresh `histkit <command> --help` output so each current command explains its purpose, safe operating mode, and supported flags without changing command behavior.

## Scope

Implement:

- expand per-command help text for `scan`, `clean`, `pick`, `doctor`, `stats`, and `restore`
- keep command help aligned with implemented behavior rather than aspirational roadmap language
- add or tighten automated tests for each command help path

## Out of scope

- new flags or command behavior changes
- root help rewrites beyond keeping terminology consistent with the prior slice
- README rewrites or broader docs refreshes
- picker behavior, scan pipeline, cleanup logic, or restore semantics changes

## Relevant skills

- `SKILLS/go-cli.md`
- `SKILLS/testing.md`

## Acceptance criteria

- each current command has help text that explains what it does in concrete terms
- `clean` help distinguishes preview behavior from apply behavior
- `restore` help distinguishes list mode from restore-by-ID mode
- automated tests verify the revised per-command help output

## Current repo state

Branch `043-command-help-detail` contains the completed per-command help refresh and matching CLI tests. The worktree is ready for staging and publish steps.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target
- Default automation runs `scan`, not destructive apply
- Wrapper logic stays outside the Go binary under `contrib/`

## Risks to watch

- Help text must stay tied to implemented behavior and avoid promising future command capabilities.
- Command help should stay compact enough for terminal use even while adding flag and mode detail.

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

## Working state

- intent: refresh the command-specific help output for the existing CLI commands
- scope: `internal/cli/{scan,clean,pick,doctor,stats,restore}.go`, matching `*_test.go` files, `SESSION.md`, and `SESSIONS/043-command-help-detail.md`
- constraints: keep command handlers thin, avoid behavior changes, preserve non-destructive defaults, keep wording aligned to current implementation
- files read:
  - `AGENTS.md`: required implementation workflow and session-record rules
  - `SESSION.md`: previous session state and working-state format
  - `ROADMAP.md`: roadmap boundaries and slice naming
  - `SKILLS/go-cli.md`: CLI constraints and handler expectations
  - `SKILLS/testing.md`: verification expectations for this slice
  - `README.md`: current command wording and safe-workflow descriptions used as a consistency source
  - `internal/cli/scan.go`: current scan help implementation
  - `internal/cli/clean.go`: current clean help implementation
  - `internal/cli/pick.go`: current pick help implementation
  - `internal/cli/doctor.go`: current doctor help implementation
  - `internal/cli/stats.go`: current stats help implementation
  - `internal/cli/restore.go`: current restore help implementation
  - `internal/cli/scan_test.go`: existing scan coverage
  - `internal/cli/clean_test.go`: existing clean coverage
  - `internal/cli/pick_test.go`: existing pick coverage
  - `internal/cli/doctor_test.go`: existing doctor coverage
  - `internal/cli/stats_test.go`: existing stats coverage
  - `internal/cli/restore_test.go`: existing restore coverage
  - `SESSIONS/000-template.md`: session note template
- files changed:
  - `internal/cli/scan.go`: expanded `scan --help` text with non-destructive behavior and flag details
  - `internal/cli/clean.go`: expanded `clean --help` text with preview vs apply guidance and flag details
  - `internal/cli/pick.go`: expanded `pick --help` text with data-source and output behavior guidance
  - `internal/cli/doctor.go`: expanded `doctor --help` text with check categories and config flag guidance
  - `internal/cli/stats.go`: expanded `stats --help` text with index-only behavior and flag guidance
  - `internal/cli/restore.go`: expanded `restore --help` text with list vs restore mode guidance and flag details
  - `internal/cli/scan_test.go`: added direct assertions for `scan --help`
  - `internal/cli/clean_test.go`: added direct assertions for `clean --help`
  - `internal/cli/pick_test.go`: added direct assertions for `pick --help`
  - `internal/cli/doctor_test.go`: added direct assertions for `doctor --help`
  - `internal/cli/stats_test.go`: added direct assertions for `stats --help`
  - `internal/cli/restore_test.go`: added direct assertions for `restore --help`
  - `SESSION.md`: recorded the active and completed session state for slice 043
  - `SESSIONS/043-command-help-detail.md`: recorded the completed session
- commands run:
  - `sed -n '1,260p' SESSION.md`: read prior session context
  - `sed -n '1,260p' ROADMAP.md`: confirmed roadmap boundaries and next slice naming
  - `sed -n '1,240p' SKILLS/go-cli.md`: loaded CLI implementation constraints
  - `sed -n '1,240p' SKILLS/testing.md`: loaded test expectations
  - `git status --short --branch`: confirmed clean starting state on `main`
  - `git checkout -b 043-command-help-detail`: created the session branch
  - `sed -n '1,240p' internal/cli/scan.go`: read current scan help
  - `sed -n '1,260p' internal/cli/clean.go`: read current clean help
  - `sed -n '1,260p' internal/cli/pick.go`: read current pick help
  - `sed -n '1,260p' internal/cli/doctor.go`: read current doctor help
  - `sed -n '1,220p' internal/cli/stats.go`: read current stats help
  - `sed -n '1,260p' internal/cli/restore.go`: read current restore help
  - `sed -n '1,260p' internal/cli/scan_test.go`: read current scan tests
  - `sed -n '1,320p' internal/cli/clean_test.go`: read current clean tests
  - `sed -n '1,260p' internal/cli/pick_test.go`: read current pick tests
  - `sed -n '1,260p' internal/cli/doctor_test.go`: read current doctor tests
  - `sed -n '1,260p' internal/cli/stats_test.go`: read current stats tests
  - `sed -n '1,320p' internal/cli/restore_test.go`: read current restore tests
  - `gofmt -w internal/cli/scan.go internal/cli/clean.go internal/cli/pick.go internal/cli/doctor.go internal/cli/stats.go internal/cli/restore.go internal/cli/scan_test.go internal/cli/clean_test.go internal/cli/pick_test.go internal/cli/doctor_test.go internal/cli/stats_test.go internal/cli/restore_test.go`: formatted the touched Go files
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`: verified CLI help and behavior tests
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`: verified the full repository test suite
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go run ./cmd/histkit clean --help`: smoke-checked one of the revised help screens through the built binary entrypoint
  - `sed -n '1,240p' SESSIONS/000-template.md`: loaded the session note template
  - `git status --short --branch`: captured the modified file set after implementation
- tests:
  - added:
    - `TestExecuteScanHelp`
    - `TestExecuteCleanHelp`
    - `TestExecutePickHelp`
    - `TestExecuteDoctorHelp`
    - `TestExecuteStatsHelp`
    - `TestExecuteRestoreHelp`
  - changed: none beyond those added test functions
  - run:
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
  - skipped: none
  - failing: none
- decisions:
  - keep per-command help terminal-oriented by limiting each screen to purpose, safe-mode distinctions, and supported flags
  - use README wording only as a consistency source for implemented behavior, not as a reason to widen the slice into broader docs work
  - test each `--help` path directly so future copy drift is caught close to the command it affects
- assumptions:
  - `NON-BLOCKING`: command help may describe current behavior in short prose without adding examples for every command; safe because this changes wording only and can be revised at low cost if later usability testing shows examples are needed
- unresolved questions:
  - none currently recorded
- next step: publish the branch and then take a small follow-up slice for shared help formatting cleanup if repetition across usage functions starts to drift

## End-of-session notes

Summary:

- Each current command now has a fuller `--help` screen that explains what it does, what mode it runs in, and which flags matter.
- `clean --help` now makes the preview-versus-apply distinction explicit, and `restore --help` now makes list-versus-restore mode explicit.
- Command help coverage was expanded with one direct `--help` assertion per command so future text changes are easier to verify.

Tests run:

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go run ./cmd/histkit clean --help`

Known failures:

- No repository test failures.

Next recommended session:

- `044-help-format-consolidation`
