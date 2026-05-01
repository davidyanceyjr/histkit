# SESSION.md

## Current session

ID: `029-clean-apply`

Status: completed

## Objective

Implement the first safe `histkit clean --apply` slice with backup, atomic rewrite, and audit logging.

## Scope

Implement:

- initial `clean` CLI wiring for `--dry-run` and `--apply`
- apply-path plumbing that reads detected shell history, evaluates built-in cleanup rules, creates a backup, rewrites the history file atomically, and appends an audit record
- focused tests for safe apply behavior in temporary history files

## Out of scope

- restore command wiring
- audit-list CLI wiring
- policy expansion beyond the current built-in rules

## Relevant skills

- `SKILLS/backup-restore.md`
- `SKILLS/testing.md`

## Acceptance criteria

- `histkit clean --dry-run` renders the existing cleanup preview through the CLI
- `histkit clean --apply` removes matching history lines for `delete` rules while still requiring backup, atomic rewrite, and audit logging
- every successful apply creates a backup, rewrites history atomically, and appends an audit record
- `go test ./...` passes

## Current repo state

The repository now has a working `clean` command with dry-run preview plus the first apply path for built-in Bash and Zsh cleanup rules.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred unless explicitly invoked
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Shell history rewrite logic must preserve valid shell-specific line formats after redaction or quarantine transforms.
- Later restore and failure-recovery slices still need to complete end-to-end rollback behavior for apply failures beyond per-file atomic writes.

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

No non-blocking questions currently recorded.

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

### Q001: Should `clean --apply` be allowed to delete entries directly for rules whose action is `delete`?

- Status: answered
- Answer: yes, remove matching history lines
- Source: direct user answer in this session
- Date answered: 2026-05-01
- Decision file updated: yes
- Risk file updated: yes

## End-of-session notes

Summary:

- Added the first `clean` CLI command with default dry-run behavior plus explicit `--apply` mode and `--shell` or `--config` support.
- Added a shell-aware cleanup apply helper that preserves unmatched raw lines, removes `delete` matches, redacts secret-bearing commands, quarantines risky commands with a placeholder, and preserves Zsh extended-history metadata during rewrites.
- Wired apply mode to create per-source backups, rewrite history atomically, and append audit records for successful applies.
- Recorded the direct-deletion policy decision in `DECISIONS.md` and `RISKS.md`.

Files changed:

- DECISIONS.md
- RISKS.md
- SESSION.md
- SESSIONS/029-clean-apply.md
- internal/cli/clean.go
- internal/cli/clean_test.go
- internal/cli/root.go
- internal/cli/root_test.go
- internal/sanitize/apply.go
- internal/sanitize/apply_test.go

Files read:

- SESSION.md
- ROADMAP.md
- SKILLS/backup-restore.md
- SKILLS/testing.md
- README.md
- docs/HUMAN_GATES.md
- docs/histkit-implementation-plan.md
- docs/OPEN_QUESTIONS.md
- DECISIONS.md
- RISKS.md
- cmd/histkit/main.go
- internal/audit/model.go
- internal/backup/create.go
- internal/backup/model.go
- internal/cli/pick.go
- internal/cli/root.go
- internal/cli/root_test.go
- internal/cli/scan.go
- internal/cli/scan_test.go
- internal/cli/stats.go
- internal/config/config.go
- internal/history/bash.go
- internal/history/detect.go
- internal/history/model.go
- internal/history/zsh.go
- internal/history/zsh_test.go
- internal/sanitize/matcher.go
- internal/sanitize/model.go
- internal/sanitize/preview.go
- internal/sanitize/preview_test.go
- internal/sanitize/quarantine.go
- internal/sanitize/redact.go
- internal/sanitize/secrets.go
- internal/sanitize/secrets_test.go
- internal/sanitize/trivial.go
- internal/sanitize/trivial_test.go

Tests added:

- `TestApplyToSourceBashRewritesDeleteRedactAndQuarantine`
- `TestApplyToSourceZshExtendedHistoryPreservesMetadata`
- `TestApplyToSourcePreservesUnparsedZshLines`
- `TestApplyToSourceRejectsUnsupportedShell`
- `TestExecuteCleanDryRunOutputsPreviewWithoutMutatingHistory`
- `TestExecuteCleanApplyRewritesHistoryCreatesBackupAndAudit`
- `TestExecuteCleanApplyRequiresBackupHistory`

Tests run:

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/sanitize ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

Known failures:

- None currently recorded.

Decisions made:

- `delete` rules remove matching history lines during `clean --apply`.
- Apply mode preserves unmatched raw lines and malformed or unparsed Zsh lines rather than dropping them.
- Final per-entry action precedence is conservative: quarantine, then redact, then delete.

Commands run:

- `git status --short --branch`
- `sed -n '1,220p' SESSION.md`
- `sed -n '1,220p' ROADMAP.md`
- `rg -n "clean|apply|backup|audit|sanitize|rewrite" internal cmd -g '!**/*_test.go'`
- `sed -n '1,240p' SKILLS/backup-restore.md`
- `sed -n '1,220p' SKILLS/testing.md`
- `sed -n '240,340p' docs/histkit-implementation-plan.md`
- `sed -n '1,260p' internal/sanitize/preview.go`
- `sed -n '1,260p' internal/sanitize/model.go`
- `rg -n "clean --dry-run|clean|apply" cmd internal README.md docs -g '!**/*_test.go'`
- `rg --files cmd internal | rg 'clean|history|scan|config|audit|backup'`
- `sed -n '1,260p' cmd/histkit/main.go`
- `sed -n '1,260p' internal/history/model.go`
- `sed -n '1,260p' internal/config/config.go`
- `sed -n '700,760p' docs/histkit-implementation-plan.md`
- `sed -n '1,260p' internal/cli/scan.go`
- `sed -n '1,260p' internal/cli/scan_test.go`
- `rg -n "Execute\\(|cobra|flag|scan|clean" internal/cli -S`
- `ls -la internal/cli`
- `sed -n '1,220p' internal/cli/root.go`
- `sed -n '1,260p' internal/history/bash.go`
- `sed -n '1,260p' internal/history/zsh.go`
- `sed -n '1,260p' internal/history/detect.go`
- `sed -n '1,260p' internal/audit/model.go`
- `sed -n '1,240p' internal/backup/create.go`
- `sed -n '150,175p' docs/HUMAN_GATES.md`
- `sed -n '80,120p' README.md`
- `sed -n '464,490p' README.md`
- `sed -n '1,260p' internal/sanitize/quarantine.go`
- `sed -n '1,260p' internal/sanitize/matcher.go`
- `sed -n '1,260p' internal/backup/model.go`
- `sed -n '1,240p' internal/cli/root_test.go`
- `sed -n '1,260p' internal/sanitize/secrets.go`
- `sed -n '1,260p' internal/sanitize/trivial.go`
- `rg -n "ActionDelete|ActionQuarantine|ActionRedact|ActionKeep" internal/sanitize -S`
- `sed -n '1,260p' docs/HUMAN_GATES.md`
- `rg -n "Q001|HUMAN_GATES|delete entries directly|apply action" -S .`
- `git checkout -b 029-clean-apply`
- `ls -1`
- `rg --files . | rg 'DECISIONS|RISKS|OPEN_QUESTIONS|internal/cli|internal/history|internal/sanitize'`
- `sed -n '1,220p' internal/cli/stats.go`
- `sed -n '1,260p' internal/cli/pick.go`
- `sed -n '1,220p' docs/OPEN_QUESTIONS.md`
- `sed -n '1,260p' internal/sanitize/redact.go`
- `sed -n '1,260p' internal/sanitize/preview_test.go`
- `sed -n '1,260p' internal/history/zsh_test.go`
- `sed -n '1,260p' DECISIONS.md`
- `sed -n '1,260p' RISKS.md`
- `rg -n "precedence|multiple matches|quarantine before|redact before|delete" README.md docs internal/sanitize -S`
- `sed -n '1,260p' internal/sanitize/secrets_test.go`
- `sed -n '1,260p' internal/sanitize/trivial_test.go`
- `gofmt -w internal/sanitize/apply.go internal/sanitize/apply_test.go internal/cli/clean.go internal/cli/clean_test.go internal/cli/root.go internal/cli/root_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/sanitize ./internal/cli`
- `gofmt -w internal/sanitize/apply_test.go`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`
- `git status --short`
- `date -u +%Y-%m-%d`

Assumptions made:

- Preserve malformed or otherwise unparsed Zsh lines verbatim during apply instead of rejecting the whole rewrite.

Risks introduced or reduced:

- Reduced: `clean --apply` now has backup, atomic rewrite, and audit coverage for built-in cleanup behavior.
- Reduced: malformed or unmatched lines are preserved verbatim instead of being lost during rewrite.
- Ongoing: later restore and failure-recovery slices still need to complete end-to-end recovery behavior.

Next recommended session:

- `030-restore-command`
