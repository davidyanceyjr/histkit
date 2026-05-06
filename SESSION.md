# SESSION.md

## Current session

ID: `047-pick-fzf-tty-wiring`

Status: completed

## Objective

Make `histkit pick` mirror `fzf` interactive output to the controlling terminal while preserving stdout capture for the selected command.

## Scope

Implement:

- update `internal/picker/fzf.go` so `fzf` stderr is mirrored to the controlling TTY when available
- preserve the existing picker contract: selected command returned to stdout, no shell-history mutation, no snippet expansion
- extend picker tests to cover TTY mirroring and fallback behavior without requiring a real interactive terminal
- carry the uncommitted `046` planning artifacts forward in this slice

## Out of scope

- changing picker candidate semantics
- changing shell wrapper behavior in `contrib/`
- changing SQLite behavior, retries, or lock handling
- broad CLI debug infrastructure beyond the existing `pick --debug`
- README or broader documentation updates

## Relevant skills

- `SKILLS/go-cli.md`
- `SKILLS/testing.md`

## Acceptance criteria

- interactive `fzf` output is no longer confined to an internal stderr buffer when a controlling TTY is available
- `histkit pick` still returns the selected command cleanly on stdout
- picker tests cover success, abort, missing `fzf`, tty mirroring, and no-tty fallback behavior

## Current repo state

Branch `047-pick-fzf-tty-wiring` contains the completed `fzf` TTY-wiring slice and the carried-forward `046` planning artifacts. An unrelated untracked file `1` exists in the worktree and was left untouched.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target
- Default automation runs `scan`, not destructive apply
- Wrapper logic stays outside the Go binary under `contrib/`
- `pick --debug` remains the existing diagnostic path for identifying pre- and post-`fzf` boundaries

## Risks to watch

- The current fix mirrors `fzf` stderr to the controlling TTY but still feeds candidates through stdin; if a specific environment needs stronger TTY ownership than stderr mirroring, a follow-up slice may still be required.
- Tests validate the launch seam with fake scripts and fake TTY files, not a full interactive terminal session.

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

No non-blocking questions currently recorded.

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

- Question: Should the implementation prefer inheriting parent stdio directly or opening `/dev/tty` explicitly for `fzf`?
  - Answer: Open `/dev/tty` explicitly and mirror `fzf` stderr to it while keeping stdout buffered for the selected command.
  - Source: `047-pick-fzf-tty-wiring` implementation in `internal/picker/fzf.go`

## Working state

- intent: make `histkit pick` show `fzf` interactive output through the controlling terminal without breaking selected-command stdout capture
- scope: `internal/picker/fzf.go`, `internal/picker/fzf_test.go`, `SESSION.md`, `SESSIONS/046-pick-fzf-tty-wiring-plan.md`, and `SESSIONS/047-pick-fzf-tty-wiring.md`
- constraints: preserve non-destructive behavior, keep selected-command stdout contract intact, avoid shell-history mutation, avoid broad picker rewrites, leave unrelated untracked file `1` untouched
- files read:
  - `AGENTS.md`: required workflow, session-record rules, and open-question protocol
  - `SESSION.md`: carried-forward planning context and implementation target
  - `ROADMAP.md`: roadmap boundaries for `pick`
  - `SKILLS/go-cli.md`: CLI constraints relevant to keeping picker behavior narrow
  - `SKILLS/testing.md`: test expectations and deterministic-fixture requirements
  - `internal/picker/fzf.go`: current buffered `fzf` launch implementation
  - `internal/picker/fzf_test.go`: existing picker tests and fake `fzf` seam
- files changed:
  - `internal/picker/fzf.go`: mirrored `fzf` stderr to `/dev/tty` when available while preserving buffered stderr capture for error reporting
  - `internal/picker/fzf_test.go`: added fake-TTY coverage and no-TTY fallback assertions; updated existing selection/abort tests to use the new seam
  - `SESSION.md`: recorded the completed session state for slice 047
  - `SESSIONS/046-pick-fzf-tty-wiring-plan.md`: carried forward the planning note created in the prior uncommitted slice
  - `SESSIONS/047-pick-fzf-tty-wiring.md`: recorded the completed implementation session
- commands run:
  - `git checkout -b 047-pick-fzf-tty-wiring`: created the implementation branch from the planning branch
  - `sed -n '1,260p' SESSION.md`: refreshed planning context and constraints
  - `sed -n '1,220p' internal/picker/fzf.go`: reviewed the current `fzf` launch path
  - `sed -n '1,260p' internal/picker/fzf_test.go`: reviewed picker tests
  - `gofmt -w internal/picker/fzf.go internal/picker/fzf_test.go`: formatted the touched Go files
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/picker`: verified picker tests including the new TTY seam coverage
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`: verified `pick` CLI behavior still passes
  - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`: verified the full repository test suite
  - `git status --short --branch`: reviewed the worktree and preserved unrelated untracked file `1`
  - `git diff -- internal/picker/fzf.go internal/picker/fzf_test.go SESSION.md SESSIONS/046-pick-fzf-tty-wiring-plan.md`: reviewed the final implementation diff
- tests:
  - added:
    - `TestSelectMirrorsFZFStderrToTTYWhenAvailable`
    - `TestSelectReturnsCapturedErrorWhenTTYUnavailable`
  - changed:
    - `TestSelectReturnsChosenCandidate`
    - `TestSelectReturnsNoSelectionForAbort`
  - run:
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/picker`
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
    - `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
  - skipped: none
  - failing: none
- decisions:
  - open `/dev/tty` explicitly for the interactive `fzf` output channel instead of inheriting the process stderr blindly
  - mirror `fzf` stderr to both the TTY and an internal buffer so user-visible UI and returned error text can coexist
  - keep the fix inside the picker layer and avoid any shell-wrapper contract changes
- assumptions:
  - `NON-BLOCKING`: mirroring `fzf` stderr to the TTY is sufficient to fix the observed hidden-UI behavior while keeping stdout capture intact
- unresolved questions:
  - none currently recorded
- next step: push the slice branch, open a PR, and wait for human approval before merge and cleanup; if user testing still shows hidden `fzf`, investigate whether stdin or broader stdio inheritance also needs TTY ownership

## End-of-session notes

Summary:

- Updated the picker so `fzf` interactive stderr is mirrored to `/dev/tty` when available, while stdout remains buffered for the selected command.
- Added picker tests for TTY mirroring and no-TTY fallback, and kept the existing selection and abort behavior intact.
- Carried the prior `046` planning note forward with this implementation slice.

Tests run:

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/picker`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

Known failures:

- No repository test failures.

Next recommended session:

- `048-pick-tty-follow-up` if user testing shows that stderr mirroring alone is insufficient in some terminals
