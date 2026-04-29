# histkit Implementation Plan

## Overview

This document describes a practical implementation plan for **histkit**, a Linux-native CLI application for shell history hygiene, reusable command snippets, and fuzzy command recall.

The implementation assumes the design constraints already established:

- real shell history remains distinct from snippets
- destructive history rewriting is not the starting point
- snippet templates are stored separately from actual user history
- the application is written in Go
- the tool operates cleanly as a user-level Linux utility
- automation is done with `systemd --user`

---

## Product Goals

histkit should provide four core capabilities:

1. read and normalize shell history from supported shells
2. detect, classify, and optionally sanitize risky or noisy history entries
3. store reusable command snippets separately from real history
4. provide fast fuzzy recall through `fzf`

The first release should prioritize safety, predictability, and auditability over aggressive automation.

---

## Safe Workflow

Production use should follow a conservative sequence:

1. run `histkit scan`
2. review findings and environment health
3. generate a cleanup preview with `histkit clean --dry-run`
4. review quarantine candidates, backups, and audit records
5. apply changes with `histkit clean --apply`
6. restore from backup if needed

This keeps indexing, review, cleanup planning, and destructive apply steps clearly separated.

---

## Core Design Rules

### 1. Keep data domains separate

The application should maintain three distinct domains:

- raw history
- sanitized/indexed history
- snippets/templates

These should only be merged at presentation time, such as inside an `fzf` picker.

### 2. Default to non-destructive behavior

The first working version should not immediately rewrite shell history files in place. It should:

- parse history
- index it
- evaluate cleanup rules
- preview proposed changes
- allow explicit apply later

### 3. Every destructive action must be recoverable

If `clean --apply` exists, then backups, audit logs, and restore support must exist before it is treated as production ready.

### 4. Treat shell history files as unstable inputs

Interactive shells may append or overwrite history after the tool runs. All logic should account for:

- concurrent terminal sessions
- shell-specific history formats
- delayed history flushes on shell exit
- partial writes or race conditions

---

## Recommended Tech Stack

### Language

**Go**

### Supporting components

- CLI framework: `cobra` or stdlib `flag`
- config parsing: `BurntSushi/toml`
- local database: SQLite
- SQLite driver: `modernc.org/sqlite` preferred for simpler builds
- logging: `log/slog`
- fuzzy picker: external `fzf`
- automation: `systemd --user`

---

## Filesystem Layout

```text
~/.config/histkit/config.toml
~/.local/share/histkit/history.db
~/.local/share/histkit/snippets.toml
~/.local/share/histkit/quarantine.log
~/.local/share/histkit/backups/
~/.local/share/histkit/audit.log
~/.cache/histkit/
```

---

## Repository Layout

```text
histkit/
├── cmd/
│   └── histkit/
│       └── main.go
├── internal/
│   ├── app/
│   ├── cli/
│   ├── config/
│   ├── history/
│   │   ├── bash.go
│   │   ├── zsh.go
│   │   ├── csh.go
│   │   ├── detect.go
│   │   └── model.go
│   ├── sanitize/
│   │   ├── rules.go
│   │   ├── redact.go
│   │   ├── delete.go
│   │   ├── trivial.go
│   │   └── entropy.go
│   ├── snippets/
│   │   ├── store.go
│   │   ├── builtin.go
│   │   └── model.go
│   ├── index/
│   │   ├── sqlite.go
│   │   ├── search.go
│   │   └── schema.go
│   ├── picker/
│   │   ├── fzf.go
│   │   └── preview.go
│   ├── service/
│   │   └── cleanup.go
│   ├── audit/
│   │   ├── log.go
│   │   └── restore.go
│   └── util/
│       ├── atomic.go
│       ├── fs.go
│       └── lock.go
├── configs/
│   └── config.example.toml
├── contrib/
│   ├── histkit-clean.service
│   ├── histkit-clean.timer
│   ├── histkit.bash
│   └── histkit.zsh
├── go.mod
└── README.md
```

---

## Command Surface

The production CLI should expose a coherent, reviewable command surface without collapsing too much behavior into one command. Implementation can still land in phases.

```text
histkit scan
histkit clean
histkit pick
histkit doctor
histkit stats
histkit restore
histkit quarantine list
histkit quarantine show
histkit quarantine restore
histkit snippets list
histkit snippets add
histkit snippets remove
histkit audit list
```

### Command responsibilities

#### `histkit scan`

- detect shell history sources
- parse entries
- normalize and hash commands
- evaluate rules
- update SQLite index
- report suspicious entries
- never rewrite history files

#### `histkit clean`

- generate cleanup plan
- show a diff or structured preview
- optionally apply actions
- support `--dry-run`
- support backups and audit entries
- behave as a two-phase workflow: plan first, apply second

#### `histkit pick`

- retrieve candidates from index and snippet store
- format them for `fzf`
- emit selected command to stdout

#### `histkit doctor`

- validate config
- verify writable data directories
- verify history source readability
- verify database access
- verify `fzf`
- verify systemd user unit presence when configured

#### `histkit stats`

- show indexed entry counts
- show counts by shell
- show duplicate reduction
- show quarantine totals
- show snippet totals
- show common commands

#### `histkit restore`

- list available backups if no identifier is given
- restore a selected snapshot
- verify integrity before replacement
- refuse unsafe overwrite unless forced

#### `histkit quarantine list`

- list recoverable quarantined entries

#### `histkit quarantine show`

- show a quarantined entry with source and rule metadata

#### `histkit quarantine restore`

- restore or export a quarantined entry for recovery

#### `histkit audit list`

- list cleanup runs and apply history
- expose audit-friendly metadata for review and compliance

---

## Internal Data Model

### HistoryEntry

```go
type HistoryEntry struct {
	ID         string
	Shell      string
	SourceFile string
	RawLine    string
	Command    string
	Timestamp  *time.Time
	ExitCode   *int
	SessionID  string
	Hash       string
}
```

### RuleMatch

```go
type ActionType string

const (
	ActionKeep       ActionType = "keep"
	ActionDelete     ActionType = "delete"
	ActionRedact     ActionType = "redact"
	ActionQuarantine ActionType = "quarantine"
)

type RuleMatch struct {
	RuleName string
	Reason   string
	Confidence string
	Action   ActionType
	Before   string
	After    string
}
```

### Snippet

```go
type Snippet struct {
	ID           string
	Title        string
	Command      string
	Description  string
	Tags         []string
	Placeholders map[string]string
	Shells       []string
	Safety       string
}
```

---

## Configuration Plan

Use TOML for user-facing configuration.

### Required configuration sections

- `[general]`
- `[cleanup]`
- `[security]`
- `[fzf]`
- `[snippets]`
- `[[rules]]`

### Example configuration

```toml
[general]
default_shell = "bash"
backup_history = true
dry_run = true
preview_diff = true

[cleanup]
dedupe = true
drop_trivial = true
drop_multiline_paste = true
max_history_age_days = 365

[cleanup.trivial]
commands = ["clear", "pwd", "ls", "ll"]

[security]
enabled = true
action = "quarantine"
remove_inline_passwords = true
remove_private_keys = true
remove_tokens = true
default_confidence_threshold = "medium"

[fzf]
preview = true
layout = "reverse"
show_source_labels = true

[snippets]
enabled = true
builtin = true
user_file = "~/.local/share/histkit/snippets.toml"

[[rules]]
name = "remove-private-key-block"
type = "contains"
pattern = "BEGIN OPENSSH PRIVATE KEY"
action = "delete"

[[rules]]
name = "redact-kubectl-token"
type = "regex"
pattern = '''kubectl.*--token[ =][^ ]+'''
action = "redact"
```

---

## Storage Plan

Use SQLite as the canonical local index and metadata store.

### Why SQLite

- fast local indexed queries
- no daemon required
- stable for single-user CLI applications
- easy to back up
- supports auditing and statistics
- better than reparsing large history files on every interaction

### Suggested tables

#### `history_entries`

Stores normalized history entries.

Suggested columns:

- `id`
- `shell`
- `source_file`
- `raw_line`
- `command`
- `timestamp`
- `exit_code`
- `session_id`
- `hash`
- `ingested_at`

#### `history_actions`

Stores decisions produced by cleanup rules.

Suggested columns:

- `id`
- `history_entry_id`
- `rule_name`
- `action`
- `before_value`
- `after_value`
- `reason`
- `created_at`

#### `snippets`

Stores reusable snippet records.

Suggested columns:

- `id`
- `title`
- `command`
- `description`
- `tags`
- `shells`
- `safety`
- `created_at`
- `updated_at`

#### `runs`

Stores scan or cleanup execution metadata.

Suggested columns:

- `id`
- `command`
- `started_at`
- `finished_at`
- `status`
- `notes`

#### `backups`

Stores metadata about file backups and restore points.

Suggested columns:

- `id`
- `source_file`
- `backup_path`
- `created_at`
- `checksum`

---

## History Parsing Plan

### Supported shells for initial release

- Bash
- Zsh

### Deferred shell support

- Csh/Tcsh

The parser layer should be abstracted behind shell-specific implementations.

### Parser interface

```text
HistorySource
  - detect()
  - parse()
  - serialize()
  - lock_strategy()
```

### Notes by shell

#### Bash

Usually line-oriented plain history.

Issues:

- timestamps may be absent or configured differently
- multiline commands may be awkwardly represented
- command provenance is limited

#### Zsh

Often uses extended history format.

Example:

```text
: 1712959000:0;command
```

Issues:

- timestamps and duration metadata may exist
- parser must correctly split metadata from command text
- escaped or multiline content must be handled carefully

---

## Sanitization Engine Plan

The sanitizer should operate as a rule engine over normalized entries.

It should treat two related but distinct concerns separately:

- security sanitization: credentials, tokens, private key material, and other sensitive values
- history hygiene: duplicates, trivial commands, large paste blobs, stale entries, and noisy recall data

### Rule classes

- exact match
- contains
- regex
- keyword group
- heuristic detector

### Supported actions

- keep
- redact
- delete
- quarantine

For many security-sensitive detections, `redact` or `quarantine` should be preferred over immediate deletion.

### Good initial detections

- `BEGIN OPENSSH PRIVATE KEY`
- bearer tokens
- inline passwords
- URLs with embedded credentials
- cloud access key patterns
- suspicious high-entropy tokens
- accidental large paste blobs

### Bad initial detections

Do not start with broad rules like:

- every command containing `ssh`
- every command containing `sudo`
- every command containing `openssl`
- every command containing `kubectl`

Those will create false positives and destroy trust.

### Sanitization output behavior

For each matched command, provide:

- matched rule
- confidence or severity
- action
- original value
- transformed value if redacted
- reason

This should be visible during dry-run previews and logged for audit purposes.

---

## Snippet System Plan

Snippets must not be stored in real shell history.

### Snippet requirements

Each snippet should support:

- unique ID
- title
- command template
- description
- tags
- shell compatibility
- optional placeholders
- safety classification

### Example snippet record

```toml
[[snippets]]
id = "find-delete-pyc"
title = "Delete Python cache files"
command = "find {{path}} -type f -name '*.pyc' -delete"
description = "Delete .pyc files under a path"
tags = ["find", "python", "cleanup"]
shells = ["bash", "zsh"]
safety = "medium"
```

### Snippet operations

- list
- add
- remove
- show
- edit
- validate
- import builtins
- search through `pick`

Placeholder expansion can be deferred if needed. The first version may simply return the template as-is.

---

## Picker Plan

Do not build an internal fuzzy finder. Use `fzf`.

### Picker flow

1. query recent and frequent history from SQLite
2. load snippets
3. merge items into one candidate stream
4. prefix each entry with source labels such as `[history]` or `[snippet]`
5. invoke `fzf`
6. parse the selected line
7. emit the selected command to stdout

### Candidate format

```text
[history]  find . -type f -name '*.tmp' -delete
[snippet]  find {{path}} -type f -name '{{pattern}}' -exec {{cmd}} {} \;
```

### Preview pane

The preview text should eventually show:

- full command
- source type
- tags
- last used time
- safety level
- snippet description

Snippet placeholders may be returned unchanged from `pick` unless a shell-side expansion workflow is configured.

---

## Shell Integration Plan

The Go binary should remain shell-agnostic.

Shell-specific integration belongs in lightweight wrapper scripts.

### Bash/Zsh wrapper behavior

A shell function can:

1. call `histkit pick`
2. capture stdout
3. place the resulting command into the shell buffer

This keeps the binary portable and avoids embedding shell-line-editor coupling into the main executable.

---

## Backup and Restore Plan

Before any in-place history mutation is allowed, backups must exist.

### Backup requirements

- per-run backup directory
- timestamped identifiers
- checksum for integrity
- metadata stored in SQLite and/or audit log
- explicit restore command

### Restore requirements

- enumerate available backups
- validate backup readability
- replace target file atomically
- log restore event
- refuse unsafe overwrite unless forced

---

## Audit Plan

Every cleanup run should be traceable.

### Audit records should include

- run identifier
- timestamps
- shell source
- matched rules
- confidence or severity summaries
- counts by action
- backup identifier
- whether apply mode was used

A human-readable log file plus structured DB records is the right compromise.

---

## systemd User-Service Plan

Automation should be done with `systemd --user`.

### Timer unit

```ini
[Unit]
Description=Run histkit cleanup periodically

[Timer]
OnBootSec=5m
OnUnitActiveSec=12h
Persistent=true

[Install]
WantedBy=timers.target
```

### Service unit

```ini
[Unit]
Description=Scan and index shell history

[Service]
Type=oneshot
ExecStart=%h/.local/bin/histkit scan --config %h/.config/histkit/config.toml
```

### Automation guidance

Default automation should remain conservative. Prefer scheduled scans and report generation before enabling automatic apply mode.

---

## Security Model

histkit handles highly sensitive user data. The implementation should assume shell history may contain:

- passwords
- API keys
- SSH private material
- internal hostnames
- sensitive filesystem paths
- production commands

### Security requirements

- file permissions should be restrictive
- backups should be local and private
- dry-run should be the default
- quarantine should be preferred over deletion during initial adoption
- logs must avoid re-exposing secrets when possible
- redaction should be applied before storing visible previews if needed

---

## Milestone Plan

## Milestone 1: parser and index foundation

### Scope

- config loading
- bash parser
- zsh parser
- SQLite schema
- `scan`
- `stats`

### Success criteria

- parses bash and zsh history
- stores normalized entries in SQLite
- prints basic counts and shell breakdown

---

## Milestone 2: fuzzy recall foundation

### Scope

- snippet store
- builtin snippets
- `pick`
- `fzf` integration
- shell labels in picker output

### Success criteria

- `histkit pick` can return indexed history entries
- snippets appear alongside history
- selected command is emitted to stdout cleanly

---

## Milestone 3: sanitizer engine

### Scope

- rule engine
- dry-run cleanup preview
- trivial command filtering
- token/private-key detection
- quarantine storage

### Success criteria

- matched entries are classified consistently
- dry-run preview shows what would happen and why
- quarantine records are recoverable

---

## Milestone 4: safe apply and restore

### Scope

- backup creation
- atomic rewrite support
- audit logs
- restore command
- apply mode safeguards

### Success criteria

- `clean --apply` is backed by recoverable snapshots
- restore works from backup IDs
- all destructive changes are logged

---

## Milestone 5: automation and polish

### Scope

- `systemd --user` units
- doctor command improvements
- shell wrapper scripts
- ranking improvements
- preview pane polish

### Success criteria

- periodic scans or cleanups work as a user service
- shell wrappers are usable in Bash and Zsh
- environment diagnostics are clear and actionable

---

## Initial Coding Order

Do not start by trying to implement everything.

### First coding session target

This is the recommended foundation for implementation order, not the final user-facing product boundary.

Build:

1. `histkit scan`
2. bash parser
3. zsh parser
4. SQLite storage
5. `histkit stats`

### Second coding session target

Build:

1. snippet store
2. builtin snippets
3. `histkit pick`
4. `fzf` integration

### Third coding session target

Build:

1. rule engine
2. dry-run preview
3. quarantine support

Only after those stages should in-place history mutation and broader recovery workflows enter the picture.

---

## Risks and Failure Modes

### 1. History rewrite races

Shells may overwrite cleaned history after the tool runs.

Mitigation:

- delay in-place rewrite support
- prefer derived index first
- warn about active shell sessions
- keep backups

### 2. False positives in secret detection

Overbroad rules will remove legitimate commands.

Mitigation:

- start with narrow high-confidence rules
- prefer quarantine over delete
- require review during dry-run
- keep audit trail

### 3. Configuration sprawl

Too many knobs too early will make the app fragile.

Mitigation:

- keep config minimal in v1
- add advanced rules later
- provide a clean example config

### 4. Shell-specific edge cases

History formats vary and can be ugly.

Mitigation:

- support Bash and Zsh first
- build parser tests from real-world fixtures
- defer Csh/Tcsh until the main pipeline is stable

### 5. Loss of user trust

One bad cleanup can kill adoption.

Mitigation:

- default to conservative mode
- make dry-run explicit and useful
- never silently inject snippets into real history
- make restore simple

---

## Testing Plan

### Unit tests

- bash history parsing
- zsh extended history parsing
- rule matching
- redaction transforms
- config loading
- snippet serialization

### Integration tests

- full scan against fixture history files
- database population
- dry-run cleanup generation
- backup creation and restore
- picker candidate generation

### Manual validation

- multiple interactive shell sessions open
- shell exit after scan
- cleanup preview against realistic sensitive history
- `fzf` wrapper behavior in Bash and Zsh
- systemd user timer execution

---

## Release Readiness Checklist

Before calling the first release usable, verify:

- Bash and Zsh history parsing are stable
- dry-run cleanup is trustworthy
- snippets are clearly separate from history
- backups and restore are functional
- audit records exist for destructive actions
- `doctor` catches missing `fzf`, bad config, and unwritable paths
- shell wrappers behave sensibly
- systemd user automation is optional and documented

---

## Final Recommendation

The implementation should be conservative, boring, and auditable.

That is the right tradeoff for a tool that touches shell history and potentially handles secrets. The early goal is not maximum cleverness. The early goal is to become trusted.

The strongest implementation sequence is:

1. parse and index
2. search and recall
3. classify and preview
4. back up and apply
5. automate carefully

Anything more aggressive than that at the start is asking for trouble.
