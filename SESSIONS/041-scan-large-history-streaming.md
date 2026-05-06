# 041-scan-large-history-streaming

Status: completed

## Summary

Reworked `histkit scan` to stream parsed history entries and write them in bounded SQLite batches so large history files no longer require full-file in-memory accumulation before indexing.

## Objective completed or not completed

Completed.

## Files read

- `AGENTS.md` - confirmed required session workflow and closeout requirements.
- `SESSION.md` - established the active working-state format that had to be preserved.
- `ROADMAP.md` - confirmed this was a post-roadmap bug-fix slice rather than new milestone feature work.
- `SKILLS/testing.md` - required deterministic tests and repository verification.
- `SKILLS/history-parsing.md` - reinforced parser safety and warning handling constraints.
- `SKILLS/go-cli.md` - reinforced keeping command logic thin and pushing behavior under `internal/`.
- `internal/cli/scan.go` - identified full-file parse accumulation in the scan hot path.
- `internal/history/bash.go` - identified whole-slice Bash parsing and default scanner limits.
- `internal/history/zsh.go` - identified whole-slice Zsh parsing and default scanner limits.
- `internal/history/detect.go` - extended parser selection to support streaming parsers.
- `internal/history/model.go` - confirmed entry and warning validation expectations.
- `internal/index/writer.go` - confirmed batched writes could reuse the existing writer safely.
- `internal/index/writer_test.go` - verified duplicate handling and rollback behavior expectations.
- `internal/cli/scan_test.go` - extended command-level regression coverage.
- `internal/history/detect_test.go` - extended parser selector coverage.
- `internal/history/bash_test.go` - validated expected Bash parser behavior to preserve.
- `internal/history/zsh_test.go` - validated expected Zsh parser behavior to preserve.

## Files changed

- `internal/cli/scan.go` - replaced whole-file parse accumulation with streaming parse callbacks plus 1000-entry write batching.
- `internal/cli/scan_test.go` - added regression tests for multi-batch scans and long history lines.
- `internal/history/bash.go` - added a streaming Bash parser and raised scanner buffer capacity.
- `internal/history/detect.go` - added stream-parser selection alongside the existing parser API.
- `internal/history/detect_test.go` - added stream-parser selector tests.
- `internal/history/zsh.go` - added a streaming Zsh parser and replaced fixed scanner limits with buffered line reads.
- `internal/history/lines.go` - added shared line-reading helpers for unbounded history-line handling.
- `SESSION.md` - updated the active working state for this session.
- `SESSIONS/041-scan-large-history-streaming.md` - recorded this session artifact.

## Tests added

- `TestExecuteScanStreamsLargeHistoryInBatches`
- `TestExecuteScanAcceptsLongHistoryLine`
- `TestStreamParserForShell`
- `TestStreamParserForShellRejectsUnsupportedShell`

## Tests run

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/history ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`

## Known failures

- No repository test failures.
- No known parser size ceiling remains from fixed scanner token limits.

## Commands run

- `git status --short --branch`
- `sed -n '1,260p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '1,220p' SKILLS/history-parsing.md`
- `sed -n '1,220p' SKILLS/go-cli.md`
- `git checkout -b 041-scan-large-history-streaming`
- `sed -n '1,240p' internal/cli/scan.go`
- `sed -n '1,260p' internal/history/bash.go`
- `sed -n '1,260p' internal/history/zsh.go`
- `sed -n '1,260p' internal/index/writer.go`
- `sed -n '1,260p' internal/cli/scan_test.go`
- `rg -n "type .*Parser|ParserForShell|ParseBash|ParseZsh|WriteHistoryEntries\\(" internal -S`
- `sed -n '1,260p' internal/history/detect.go`
- `sed -n '1,260p' internal/history/model.go`
- `sed -n '1,260p' internal/index/writer_test.go`
- `sed -n '1,260p' internal/history/bash_test.go`
- `sed -n '1,260p' internal/history/zsh_test.go`
- `sed -n '140,240p' internal/history/detect_test.go`
- `gofmt -w internal/history/detect.go internal/history/bash.go internal/history/zsh.go internal/cli/scan.go internal/history/detect_test.go internal/cli/scan_test.go`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./internal/history ./internal/cli`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
- `git diff --stat`

## Decisions made

- Preserve the existing slice-returning parsers for compatibility and add streaming variants for scan-specific performance work.
- Keep batching inside `scan` so the index writer contract stays stable for the rest of the codebase.
- Replace `bufio.Scanner` with a buffered line reader so very large single-line history entries are accepted without a fixed token ceiling.

## Assumptions made

- None.

## Unresolved questions

- No active blocking questions.
- No active non-blocking questions.

## Risks introduced or reduced

- Reduced: large history files no longer require full-file entry accumulation before index writes begin.
- Reduced: very large single-line history entries no longer fail at a fixed scanner token limit.
- Remaining: manual testing should still confirm acceptable latency on the original large history file.

## Next slice recommendation

- Manually test `histkit scan` against the real large history file that exposed the issue, then decide whether progress reporting or a larger single-line ceiling is justified.
