# Session 006: Source Detection

## Objective

Implement initial shell history source detection for Bash and Zsh.

## Completed

- Added `internal/history/detect.go` with canonical Bash and Zsh source candidates.
- Added detection that returns only existing history files.
- Added shell filtering for supported source types.
- Added parser resolution for supported shells.
- Added deterministic tests for detection behavior.
- Recorded the canonical-path assumption as a non-blocking open question.

## Files changed

- internal/history/detect.go
- internal/history/detect_test.go
- docs/OPEN_QUESTIONS.md
- SESSION.md
- SESSIONS/006-source-detection.md

## Tests added

- TestCandidateSources
- TestCandidateSourcesRequiresHome
- TestDetectSourcesFindsExistingFiles
- TestDetectSourcesFiltersByShell
- TestDetectSourcesIgnoresMissingFiles
- TestDetectSourcesRejectsUnsupportedShell
- TestParserForShell
- TestParserForShellRejectsUnsupportedShell

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Decisions

- Detect only canonical `~/.bash_history` and `~/.zsh_history` files in the initial source-detection slice.
- Keep parser lookup alongside source detection so later scan wiring can resolve supported parsers cleanly.

## Known issues

- Custom shell history path discovery remains deferred.
- CLI wiring for shell filters remains deferred.

## Next recommended session

`007-sqlite-schema`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

- Q001 copied to `docs/OPEN_QUESTIONS.md`.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

- Q001: Should initial source detection support only canonical `~/.bash_history` and `~/.zsh_history` paths?
