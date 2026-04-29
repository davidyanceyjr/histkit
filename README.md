# histkit

`histkit` is a Linux-native CLI for shell history hygiene, reusable command snippets, and fast fuzzy command recall.

It reads shell history from supported shells, builds a searchable local index, detects commands that may contain secrets or other sensitive material, and helps you review or apply cleanup rules safely. It also maintains a separate snippet library for reusable command templates and exposes both indexed history and snippets through an `fzf`-based picker.

## What histkit does

`histkit` is built around three separate data domains:

1. Raw shell history
2. Sanitized and indexed history
3. Reusable snippets

These domains stay separate by design. Snippets are not written into your real shell history by default.

Core capabilities:

- Read and normalize shell history from supported shells
- Detect, classify, redact, delete, or quarantine risky history entries
- Build a searchable local SQLite index
- Manage reusable command snippets outside shell history
- Provide interactive fuzzy recall through `fzf`
- Support backups, restore, and audit-friendly cleanup workflows

## Safe workflow

For production use, the intended workflow is:

1. Run `histkit scan`
2. Review findings and environment health with `histkit doctor` or structured output
3. Generate a cleanup preview with `histkit clean --dry-run`
4. Review quarantine candidates, backups, and audit records
5. Apply changes with `histkit clean --apply`
6. Restore from backup with `histkit restore` if needed

This keeps indexing, review, cleanup planning, and destructive apply steps clearly separated.

## Feature model

`histkit` separates two related but distinct concerns:

- Security sanitization: secrets, tokens, credentials, private key material, and other sensitive values
- History hygiene: duplicates, trivial commands, large paste blobs, stale entries, and noisy recall data

Rule matches should carry both an action and a confidence level so operators can review sensitive high-confidence findings differently from lower-risk cleanup suggestions.

## Command overview

```text
histkit scan
histkit clean [--apply] [--dry-run]
histkit pick
histkit doctor
histkit stats
histkit restore [backup-id]
histkit quarantine list
histkit quarantine show <entry-id>
histkit quarantine restore <entry-id>
histkit snippets list
histkit snippets add [options]
histkit snippets remove <snippet-id>
histkit audit list
```

## Commands

### `histkit scan`

Parse supported shell history files, normalize entries, evaluate rules, and update the local index.

`scan` is an ingest, indexing, and reporting command. It does not rewrite shell history.

Typical uses:

```sh
histkit scan
histkit scan --shell bash
histkit scan --config ~/.config/histkit/config.toml
```

### `histkit clean`

Generate or apply a cleanup plan for shell history.

`clean` has two operational modes:

- Planning mode: `histkit clean --dry-run`
- Apply mode: `histkit clean --apply`

By default, `clean` is non-destructive unless `--apply` is used. Matching entries may be kept, redacted, deleted, or quarantined depending on configuration and rules. Production apply mode assumes backups, restore support, and audit logging are enabled.

Typical uses:

```sh
histkit clean --dry-run
histkit clean --apply
histkit clean --apply --shell zsh
```

### `histkit pick`

Open an interactive fuzzy picker using `fzf` over indexed history and configured snippets.

The selected command is written to standard output. Shell integration wrappers can capture that output and place it into the interactive shell buffer. Snippet placeholders are returned unchanged unless a shell-side expansion workflow is configured.

```sh
histkit pick
```

### `histkit doctor`

Inspect the local environment and report configuration or runtime problems.

Checks include:

- Readable configuration file
- Writable state directory
- Shell history file detection
- Index availability
- `fzf` presence
- Backup path availability
- `systemd --user` unit visibility, when relevant

```sh
histkit doctor
```

```sh
histkit doctor --json
```

### `histkit stats`

Print local history and snippet statistics.

Typical output areas include:

- Indexed history entries
- Per-shell counts
- Deduplicated entries
- Rule matches
- Quarantined entries
- Snippet count
- Most frequent commands

```sh
histkit stats
```

```sh
histkit stats --json
```

### `histkit restore`

Restore history state from a backup or recorded snapshot.

If no backup identifier is given, `restore` lists available backups. Restore operations validate backup integrity before replacement and should refuse unsafe overwrite unless explicitly forced by implementation policy.

```sh
histkit restore
histkit restore 20260418T021533Z
```

### `histkit quarantine list`

List entries that were quarantined by rule actions instead of being removed immediately.

```sh
histkit quarantine list
```

### `histkit quarantine show`

Show a quarantined entry with its source, matched rule, and recovery metadata.

```sh
histkit quarantine show q_01HZX2T9X9Z3M
```

### `histkit quarantine restore`

Restore a quarantined entry to the appropriate history target or export it for manual recovery.

```sh
histkit quarantine restore q_01HZX2T9X9Z3M
```

### `histkit snippets list`

List available snippets from configured snippet stores.

```sh
histkit snippets list
```

### `histkit snippets add`

Add a snippet to the snippet store.

```sh
histkit snippets add --title "Delete .pyc files" \
  --command "find {{path}} -type f -name '*.pyc' -delete" \
  --tag find --tag python
```

### `histkit snippets remove`

Remove a snippet by identifier.

```sh
histkit snippets remove find-delete-pyc
```

Production snippet workflows may also include:

- `histkit snippets show <snippet-id>`
- `histkit snippets edit <snippet-id>`
- `histkit snippets validate`
- `histkit snippets import`

### `histkit audit list`

List cleanup runs, audit records, and apply history for review and compliance.

```sh
histkit audit list
```

## Global options

### `--config <path>`

Use an alternate configuration file instead of the default.

### `--shell <name>`

Restrict an operation to a specific shell source.

Typical supported values:

- `bash`
- `zsh`

### `--state-dir <path>`

Override the default local state directory.

### `--verbose`

Enable more detailed logging.

### `--quiet`

Reduce non-error output.

### `--json`

Emit machine-readable output for commands that support structured output, such as `scan`, `doctor`, `stats`, and audit-related reporting.

### `--no-color`

Disable ANSI color output.

### `--version`

Print version information and exit.

### `--help`

Print command or subcommand help and exit.

## Cleanup model

`histkit` supports four primary actions for matched history entries:

- `keep`: leave the command unchanged
- `redact`: preserve the command shape but remove or mask sensitive values
- `delete`: remove the command from the cleanup result
- `quarantine`: move the command into a recoverable quarantine record for later review

Rules may be based on exact matching, keyword groups, regular expressions, structured detectors, or heuristic checks.

For many security-sensitive detections, `redact` or `quarantine` is the preferred default over immediate deletion.

Typical matched material includes:

- Pasted private key blocks
- Inline passwords
- Bearer tokens
- Cloud credentials
- URL-embedded credentials
- Large accidental paste blobs

Broad rules against commands such as `ssh`, `sudo`, `openssl`, or `kubectl` should be avoided because they generate false positives and reduce trust.

For each matched command, production-ready output should make clear:

- Matched rule
- Confidence or severity
- Selected action
- Original value
- Transformed value when redacted
- Reason

## Snippets

Snippets are reusable command templates maintained outside your real shell history.

A snippet may contain:

- A stable identifier
- A title
- A command template
- Descriptive text
- Tags
- Placeholder values
- Shell compatibility metadata
- An optional safety rating

Example snippet:

```toml
[[snippets]]
id = "find-delete-pyc"
title = "Delete Python cache files"
command = "find {{path}} -type f -name '*.pyc' -delete"
description = "Delete .pyc files under a target path"
tags = ["find", "python", "cleanup"]
shells = ["bash", "zsh"]
safety = "medium"
```

Builtin snippets may be shipped with the tool and imported into the user snippet store as needed.

## Shell integration

`histkit` is designed to work with lightweight shell wrappers for interactive use.

A wrapper function can:

1. Invoke `histkit pick`
2. Capture the selected command
3. Place the result into the current shell editing buffer

The standalone binary writes the selected command to standard output and does not require direct manipulation of shell line-editor internals.

## Configuration

Default config path:

```text
~/.config/histkit/config.toml
```

Typical data layout:

```text
~/.config/histkit/config.toml
~/.local/share/histkit/history.db
~/.local/share/histkit/snippets.toml
~/.local/share/histkit/quarantine.log
~/.local/share/histkit/backups/
~/.cache/histkit/
```

Example configuration:

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

## Files

```text
~/.config/histkit/config.toml
~/.local/share/histkit/history.db
~/.local/share/histkit/snippets.toml
~/.local/share/histkit/quarantine.log
~/.local/share/histkit/backups/
~/.local/share/histkit/audit.log
~/.config/systemd/user/histkit-clean.service
~/.config/systemd/user/histkit-clean.timer
```

## Exit status

- `0`: successful completion
- `1`: general runtime error
- `2`: configuration error
- `3`: history source parsing error
- `4`: cleanup or rewrite operation failed
- `5`: interactive picker dependency or execution failure
- `6`: restore or backup operation failed

## Examples

### Scan history and update the index

```sh
histkit scan
```

### Show a non-destructive cleanup preview

```sh
histkit clean --dry-run
```

### List audit history

```sh
histkit audit list
```

### Apply configured cleanup rules

```sh
histkit clean --apply
```

### Pick a command through `fzf`

```sh
histkit pick
```

### Show environment diagnostics

```sh
histkit doctor
```

### Print local usage statistics

```sh
histkit stats
```

### Restore from a specific backup

```sh
histkit restore 20260418T021533Z
```

## Security notes

`histkit` processes shell history, which may contain sensitive material.

Assume shell history can include:

- Credentials
- Hostnames
- Internal paths
- Tokens
- Command arguments containing private data

Recommended operating practice:

- Review cleanup results before applying them
- Keep backups enabled
- Prefer quarantine over permanent deletion when testing rules
- Use dry-run mode before enabling automation
- Restrict file permissions on configuration, state, and backup directories

Because active shells may rewrite history files on exit, in-place cleanup should be approached carefully. Derived indexes, reviewable cleanup plans, backups, and audit logs are safer than blind direct rewrites.

## systemd integration

A typical user-level timer should default to periodic scans or report generation. Automated cleanup should remain conservative and should not imply unattended destructive apply mode by default.

Example timer:

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

Matching one-shot service for scheduled scanning:

```ini
[Unit]
Description=Scan and index shell history

[Service]
Type=oneshot
ExecStart=%h/.local/bin/histkit scan --config %h/.config/histkit/config.toml
```

Install and enable as the target user, not as root, unless your deployment specifically requires otherwise.

## Dependencies

Practical interactive use requires:

- A supported shell history source
- Write access to the configured state directory
- `fzf` for the interactive picker
- SQLite support in the build or runtime configuration
- `systemd --user` for timer-based automation, if used

## Failure modes

Likely failure modes include:

- Concurrent shell sessions rewriting history after cleanup
- Shell-specific history format edge cases
- False positives from overbroad sanitization rules
- Invalid snippet placeholders
- External picker failures when `fzf` is missing or misconfigured
- Incomplete audit or backup retention policies in multi-shell environments

## See also

`bash(1)`, `zsh(1)`, `fzf(1)`, `systemd(1)`, `systemd.timer(5)`, `systemd.service(5)`

## License

Distributed under the terms of the project license.
