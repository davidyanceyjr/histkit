# ROADMAP.md

## Milestone 1: Parser and index foundation

Goal: `histkit scan` and `histkit stats`.

### Slices

1. `001-bootstrap-cli`
2. `002-config-and-paths`
3. `003-history-model`
4. `004-bash-parser`
5. `005-zsh-parser`
6. `006-source-detection`
7. `007-sqlite-schema`
8. `008-index-writer`
9. `009-scan-pipeline`
10. `010-stats-command`
11. `011-doctor-command-v1`

### Exit criteria

- Bash and Zsh fixture histories parse.
- Normalized entries are stored in SQLite.
- `stats` shows counts by shell/source.
- No history files are modified.

---

## Milestone 2: Fuzzy recall foundation

Goal: `histkit pick`.

### Slices

1. `012-snippet-model`
2. `013-snippet-store`
3. `014-builtin-snippets`
4. `015-pick-candidates`
5. `016-fzf-picker`
6. `017-shell-wrappers`

### Exit criteria

- History and snippets appear together in the picker.
- Snippets remain separate from real history.
- The selected command emits cleanly to stdout.

---

## Milestone 3: Sanitizer engine

Goal: `histkit clean --dry-run`.

### Slices

1. `018-rule-model`
2. `019-rule-matching`
3. `020-redaction-transforms`
4. `021-secret-rules`
5. `022-trivial-command-rules`
6. `023-dry-run-preview`
7. `024-quarantine-records`

### Exit criteria

- Matched entries show rule, reason, and action.
- Dry-run performs no mutation.
- Quarantine records are recoverable.

---

## Milestone 4: Safe apply and restore

Goal: controlled mutation with recovery.

### Slices

1. `025-backup-model`
2. `026-backup-creation`
3. `027-atomic-rewrite`
4. `028-audit-log`
5. `029-clean-apply`
6. `030-restore-command`
7. `031-failure-recovery-tests`

### Exit criteria

- Every apply has a backup and audit record.
- Restore works by backup ID.
- Unsafe overwrites are refused.

---

## Milestone 5: Automation and polish

Goal: optional `systemd --user` automation.

### Slices

1. `032-systemd-user-service`
2. `033-systemd-user-timer`
3. `034-doctor-systemd-checks`
4. `035-shell-wrapper-polish`
5. `036-readme-usage-flow`
6. `037-release-readiness-pass`

### Exit criteria

- Timer can run scan.
- `doctor` catches missing dependencies.
- Documentation explains the safe workflow.

---

## Hard boundaries

Before Milestone 1 exits:

- no sanitizer
- no fzf
- no `clean --apply`

Before Milestone 2 exits:

- no history mutation
- no backup system except harmless metadata scaffolding

Before Milestone 3 exits:

- no destructive cleanup

Before Milestone 4 exits:

- no automation that applies cleanup

Before Milestone 5 exits:

- no default timer that mutates user files
