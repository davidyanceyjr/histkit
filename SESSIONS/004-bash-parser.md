# Session 004: Bash Parser

## Objective

Implement the initial Bash history parser for plain line-oriented history files.

## Completed

- Added `internal/history/bash.go` with plain Bash history parsing from an `io.Reader`.
- Preserved raw lines and command text for valid Bash history entries.
- Ignored empty lines deterministically.
- Reported whitespace-only lines as parse warnings instead of producing invalid entries.
- Added a fixture-driven Bash parser test suite.

## Files changed

- internal/history/bash.go
- internal/history/bash_test.go
- testdata/history/bash/plain.hist
- SESSION.md
- SESSIONS/004-bash-parser.md

## Tests added

- TestParseBashFixture
- TestParseBashEmptyInput
- TestParseBashRequiresSourceFile
- TestParseBashRequiresReader

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Decisions

- Treat empty Bash history lines as ignorable input.
- Treat whitespace-only Bash history lines as parse warnings so they are not silently converted into invalid commands.

## Known issues

- `HISTTIMEFORMAT` and multiline Bash history remain deferred.
- Source detection and parser selection remain deferred.

## Next recommended session

`005-zsh-parser`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
