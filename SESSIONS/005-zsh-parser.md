# Session 005: Zsh Parser

## Objective

Implement the initial Zsh history parser for extended history lines.

## Completed

- Added `internal/history/zsh.go` with Zsh extended-history parsing from an `io.Reader`.
- Parsed reliable timestamps from Zsh metadata.
- Preserved command text after the first metadata separator semicolon.
- Left `ExitCode` unset because Zsh extended history stores duration, not exit status.
- Reported malformed metadata prefixes as parse warnings.
- Added fixture-driven Zsh parser tests and a Zsh history fixture.

## Files changed

- internal/history/zsh.go
- internal/history/zsh_test.go
- testdata/history/zsh/extended.hist
- SESSION.md
- SESSIONS/005-zsh-parser.md

## Tests added

- TestParseZshFixture
- TestParseZshMalformedPrefixWarning
- TestParseZshEmptyInput
- TestParseZshRequiresSourceFile
- TestParseZshRequiresReader

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Decisions

- Parse the first Zsh metadata field as a reliable timestamp.
- Validate the duration field for structure but do not map it onto `ExitCode`.
- Preserve everything after the first `;` as the command text, including additional semicolons.

## Known issues

- Multiline Zsh history remains deferred.
- Source detection and parser selection remain deferred.

## Next recommended session

`006-source-detection`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
