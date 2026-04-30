# SESSION.md

## Current session

ID: `017-shell-wrappers`

Status: completed

## Objective

Add lightweight Bash and Zsh wrapper scripts for interactive `histkit pick` usage.

## Scope

Implement:

- Bash wrapper script for loading picker output into the readline buffer
- Zsh wrapper script for loading picker output into the shell editing buffer
- README usage examples for sourcing and binding the wrapper scripts
- wrapper-focused tests

## Out of scope

- changes to `histkit pick` core command behavior
- preview panes or placeholder expansion
- automatic shell-profile installation
- sanitizer or cleanup features

## Relevant skills

- `SKILLS/fzf-picker.md`
- `SKILLS/testing.md`

## Acceptance criteria

- repository contains sourceable Bash and Zsh wrapper scripts
- wrappers call `histkit pick`, capture stdout, and place the result into the active shell buffer
- wrappers are documented in `README.md`
- automated validation covers wrapper behavior where feasible
- `go test ./...` passes

## Current repo state

The repository now includes Bash and Zsh wrapper scripts under `contrib/` so interactive shells can bind `Ctrl-R` to `histkit pick` and replace the current editing buffer with the selected command.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Zsh is not installed in this environment, so runtime validation here covers Bash directly and verifies the Zsh wrapper structurally.
- Wrapper installation remains manual; users must source the script and bind the helper in their shell config.

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

No non-blocking questions currently recorded.

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

No questions answered yet.

## End-of-session notes

Summary:

- Added `contrib/histkit.bash` and `contrib/histkit.zsh` as lightweight sourceable wrappers around `histkit pick`.
- Added `contrib/wrappers_test.go` to exercise the Bash wrapper directly and verify key Zsh wrapper behavior from file contents.
- Updated `README.md` with shell-wrapper setup examples for Bash and Zsh.

Files changed:

- README.md
- contrib/histkit.bash
- contrib/histkit.zsh
- contrib/wrappers_test.go
- SESSION.md
- SESSIONS/017-shell-wrappers.md

Files read:

- SESSION.md
- ROADMAP.md
- SKILLS/fzf-picker.md
- docs/histkit-implementation-plan.md
- README.md
- configs/config.example.toml

Tests added:

- `TestBashWrapperLoadsSelectionIntoReadlineBuffer`
- `TestBashWrapperLeavesBufferUntouchedOnEmptySelection`
- `TestZshWrapperScriptContainsBindingHelper`

Tests run:

- `bash -n contrib/histkit.bash`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./contrib`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Keep shell integration outside the Go binary in sourceable wrapper scripts under `contrib/`.
- Provide explicit binding helper functions rather than auto-modifying shell key bindings when sourced.
- Replace the current editing buffer with the selected command rather than attempting partial insertion.

Commands run:

- `git checkout -b 017-shell-wrappers`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/fzf-picker.md`
- `sed -n '665,710p' docs/histkit-implementation-plan.md`
- `sed -n '330,370p' README.md`
- `sed -n '460,478p' README.md`
- `rg --files | rg '^contrib/|shell|bash|zsh'`
- `sed -n '1,220p' configs/config.example.toml`
- `bash --version | head -n 1`
- `zsh --version`
- `rg -n "contrib/|histkit.bash|histkit.zsh|Shell integration|wrapper" README.md docs -S`
- `sed -n '850,885p' docs/histkit-implementation-plan.md`
- `sed -n '990,1018p' docs/histkit-implementation-plan.md`
- `gofmt -w contrib/wrappers_test.go`
- `bash -n contrib/histkit.bash`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./contrib`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- A manual source-and-bind workflow is sufficient for the first shell-wrapper slice.
- Structural verification of the Zsh wrapper is acceptable in this environment because `zsh` is not installed locally.

Risks introduced or reduced:

- Reduced: interactive shell integration now exists without coupling the Go binary to shell editor internals.
- Ongoing: wrapper ergonomics may need polish later, especially once preview panes or richer shell integration exist.

Next recommended session:

- `018-rule-model`
