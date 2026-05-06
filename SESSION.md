# SESSION.md

## Current session

ID: `045-pick-debug-diagnostics`

Status: completed

## Objective

Add opt-in diagnostics for `histkit pick` so a user can identify whether the command is stalling before database access, candidate loading, or `fzf` launch.

## Scope

Implement:

- add a `pick` debug flag that emits progress diagnostics to stderr
- keep normal `pick` behavior unchanged when debug is disabled
- cover the new diagnostics path with deterministic CLI tests

## Out of scope

- changing picker candidate semantics
- changing `fzf` terminal wiring
- adding SQLite lock timeouts or retries
- changing scan, stats, or doctor behavior
- README or broader documentation updates

## Relevant skills

- `SKILLS/go-cli.md`
- `SKILLS/testing.md`

## Acceptance criteria

- `histkit pick --debug` emits stage-level diagnostics before `fzf` launches
- standard `histkit pick` remains silent on stderr during successful runs
- CLI tests cover the debug flag and still pass

## Current repo state

Branch `045-pick-debug-diagnostics` contains the completed `pick` diagnostics slice and is ready for staging, publish, and merge workflow.

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

- Debug output is written to stderr only when explicitly enabled and must stay off the normal `pick` path.
- The diagnostics identify the blocking stage but do not yet change SQLite lock behavior or provide automatic recovery.

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

- intent: add opt-in `histkit pick` diagnostics to identify the pre-`fzf` blocking stage
- scope: `internal/cli/{root,pick}.go`, `internal/cli/pick_test.go`, `SESSION.md`, and `SESSIONS/045-pick-debug-diagnostics.md`
- constraints: preserve non-destructive behavior, keep normal `pick` stderr silent, avoid picker rewrites, avoid touching shell-history mutation paths
- files read:
  - `AGENTS.md`: required workflow, session-record rules, and publish expectations
  - `SESSION.md`: prior session state and structured working-state format
  - `ROADMAP.md`: roadmap boundaries and slice naming
  - `SKILLS/go-cli.md`: CLI implementation constraints
  - `SKILLS/testing.md`: test expectations
  - `internal/cli/root.go`: command dispatch path that needed stderr plumbing for `pick`
  - `internal/cli/pick.go`: current `pick` flag parsing and runtime sequence
  - `internal/cli/pick_test.go`: existing `pick` behavior and help coverage
  - `internal/picker/fzf.go`: confirmed `fzf` launch happens after candidate loading
  - `internal/picker/candidates.go`: confirmed candidate loading is the pre-`fzf` boundary
  - `internal/index/schema.go`: reviewed SQLite open and schema init path
  - `internal/index/picker.go`: reviewed recent-history query used by `pick`
  - `internal/snippets/store.go`: reviewed snippet store reads in the pre-`fzf` path
  - `SESSIONS/000-template.md`: session note template
  - `/home/opsman/.codex/plugins/cache/openai-curated/github/9d07fd08/skills/yeet/SKILL.md`: publish workflow requirements
- files changed:
  - `internal/cli/root.go`: passed stderr into `pick` command execution
  - `internal/cli/pick.go`: added `--debug` flag and stage/timing diagnostics for the pre-`fzf` path plus selection outcome logging
  - `internal/cli/pick_test.go`: covered debug flag help text and stderr diagnostics while preserving normal behavior assertions
  - `SESSION.md`: recorded the active and completed session state for slice 045
  - `SESSIONS/045-pick-debug-diagnostics.md`: recorded the completed session
- commands run:
  - `sed -n '1,260p' SESSION.md`: read prior session context
  - `sed -n '1,220p' ROADMAP.md`: confirmed roadmap boundaries and slice naming
  - `sed -n '1,220p' SKILLS/go-cli.md`: loaded CLI implementation constraints
  - `sed -n '1,220p' SKILLS/testing.md`: loaded test expectations
  - `git status --short --branch`: confirmed clean starting state on `main`
  - `git checkout -b 045-pick-debug-diagnostics`: created the session branch
  - `sed -n '1,260p' /home/opsman/.codex/plugins/cache/openai-curated/github/9d07fd08/skills/yeet/SKILL.md`: loaded publish workflow instructions
  - `sed -n '1,260p' internal/cli/pick.go`: read `pick` command implementation
  - `sed -n '1,260p' internal/cli/pick_test.go`: read `pick` tests
  - `sed -n '1,220p' SESSIONS/000-template.md`: loaded the session note template
  - `sed -n '1,260p' internal/cli/root.go`: read command dispatch for stderr plumbing
  - `sed -n '1,220p' internal/picker/fzf.go`: verified `fzf` launch position in the runtime path
  - `sed -n '1,260p' internal/picker/candidates.go`: verified candidate loading semantics before `fzf`
  - `sed -n '1,220p' internal/index/schema.go`: reviewed SQLite open and schema init behavior
  - `sed -n '1,220p' internal/index/picker.go`: reviewed recent-history query path
  - `sed -n '1,260p' internal/snippets/store.go`: reviewed snippet store reads
  - `gofmt -w internal/cli/root.go internal/cli/pick.go internal/cli/pick_test.go`: formatted the touched Go files
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`: verified CLI tests including the new debug path
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`: verified the full repository test suite
  - `git diff -- internal/cli/root.go internal/cli/pick.go internal/cli/pick_test.go`: reviewed the final implementation diff
  - `gh --version`: verified GitHub CLI availability for publish flow
  - `gh auth status`: verified authenticated GitHub CLI session
  - `git remote get-url origin`: confirmed repository remote
  - `gh repo view --json nameWithOwner,defaultBranchRef`: confirmed target repo and default branch
- tests:
  - added:
    - `TestExecutePickDebugWritesStageDiagnosticsToStderr`
  - changed:
    - `TestExecutePickHelp`
  - run:
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
  - skipped: none
  - failing: none
- decisions:
  - make diagnostics opt-in behind `histkit pick --debug` so normal interactive use stays unchanged
  - write diagnostics to stderr so stdout remains reserved for the selected command
  - instrument each pre-`fzf` stage with start and completion timing to identify the exact blocking boundary quickly
- assumptions:
  - `NON-BLOCKING`: a command-specific `--debug` flag is a safe CLI addition because it is opt-in, preserves normal behavior, and has low reversal cost if later replaced by a shared debug facility
- unresolved questions:
  - none currently recorded
- next step: use `histkit pick --debug` against the user environment and decide from the emitted stage boundary whether the next slice should target SQLite lock handling, path/config diagnostics, or `fzf` terminal wiring

## End-of-session notes

Summary:

- Added `histkit pick --debug` so the command now emits step-by-step progress diagnostics to stderr around home detection, config/path resolution, SQLite open/init, candidate loading, and `fzf` launch.
- Kept the default `pick` path unchanged: stdout still carries only the selected command, and successful non-debug runs remain silent on stderr.
- Added CLI coverage for the new debug diagnostics path.

Tests run:

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

Known failures:

- No repository test failures.

Next recommended session:

- `046-pick-stall-root-cause-follow-up`
