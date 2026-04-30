# Session 017: Shell Wrappers

## Objective

Add lightweight Bash and Zsh wrapper scripts for interactive `histkit pick` usage.

## Completed

- Added `contrib/histkit.bash` with a helper that runs `histkit pick`, captures stdout, and replaces the Bash readline buffer.
- Added `contrib/histkit.zsh` with a helper that runs `histkit pick`, captures stdout, and replaces the Zsh editing buffer.
- Added explicit binding helper functions for both shells so users can opt into the `Ctrl-R` binding from their shell config.
- Updated `README.md` with shell-wrapper sourcing and binding examples.
- Added wrapper-focused tests, including direct Bash wrapper execution.

## Files changed

- SESSION.md
- SESSIONS/017-shell-wrappers.md
- README.md
- contrib/histkit.bash
- contrib/histkit.zsh
- contrib/wrappers_test.go

## Files read

- SESSION.md
- ROADMAP.md
- SKILLS/fzf-picker.md
- docs/histkit-implementation-plan.md
- README.md
- configs/config.example.toml

## Tests added

- `TestBashWrapperLoadsSelectionIntoReadlineBuffer`
- `TestBashWrapperLeavesBufferUntouchedOnEmptySelection`
- `TestZshWrapperScriptContainsBindingHelper`

## Tests run

```bash
bash -n contrib/histkit.bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./contrib
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 017-shell-wrappers
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SKILLS/fzf-picker.md
sed -n '665,710p' docs/histkit-implementation-plan.md
sed -n '330,370p' README.md
sed -n '460,478p' README.md
rg --files | rg '^contrib/|shell|bash|zsh'
sed -n '1,220p' configs/config.example.toml
bash --version | head -n 1
zsh --version
rg -n "contrib/|histkit.bash|histkit.zsh|Shell integration|wrapper" README.md docs -S
sed -n '850,885p' docs/histkit-implementation-plan.md
sed -n '990,1018p' docs/histkit-implementation-plan.md
gofmt -w contrib/wrappers_test.go
bash -n contrib/histkit.bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./contrib
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Keep wrapper logic outside the Go binary and store it under `contrib/`.
- Expose binding helper functions instead of auto-binding on source.
- Replace the current shell buffer with the selected command.

## Assumptions

- Manual source-and-bind setup is sufficient for the first wrapper slice.
- Zsh wrapper structure checks are acceptable when `zsh` is unavailable locally.

## Known issues

- Wrapper installation is not automated.
- Zsh runtime behavior was not executed in this environment because `zsh` is not installed.

## Risks reduced

- Interactive shell integration now exists without making the binary shell-aware.

## Next recommended session

`018-rule-model`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
