# histkit

`histkit` is a Linux-native CLI for shell history hygiene, reusable command snippets, and fast fuzzy command recall.

It reads shell history from supported shells, builds a searchable local index, applies built-in cleanup rules conservatively, and exposes indexed history plus snippets through an `fzf`-based picker.

## What histkit does

`histkit` is built around three separate data domains:

1. Raw shell history
2. Sanitized and indexed history
3. Reusable snippets

These domains stay separate by design. Snippets are not written into your real shell history by default.

Core capabilities:

- Read and normalize shell history from supported shells
- Build and query a local SQLite history index
- Preview built-in cleanup actions before rewriting anything
- Apply cleanup changes with per-run backups and audit logging
- Restore from recorded backups
- Provide interactive fuzzy recall through `fzf`
- Load reusable snippets from a separate snippet store

## Safe workflow

For production use, the intended workflow is:

1. Run `histkit doctor` to check config, state paths, history detection, `fzf`, and optional `systemd --user` units.
2. Run `histkit scan` to parse history sources and refresh the local index.
3. Optionally inspect the indexed result with `histkit stats` or `histkit pick`.
4. Generate a cleanup preview with `histkit clean --dry-run`.
5. Apply changes only after review with `histkit clean --apply`.
6. If needed, list available backups with `histkit restore`, then restore a specific backup ID.

This keeps environment checks, indexing, preview, destructive apply, and recovery clearly separated.

## Feature model

`histkit` separates two related but distinct concerns:

- Security sanitization: secrets, tokens, credentials, private key material, and other sensitive values
- History hygiene: duplicates, trivial commands, large paste blobs, stale entries, and noisy recall data

Rule matches should carry both an action and a confidence level so operators can review sensitive high-confidence findings differently from lower-risk cleanup suggestions.

## Current command surface

```text
histkit scan
histkit clean [--apply] [--dry-run]
histkit pick
histkit doctor
histkit stats
histkit restore [backup-id]
```

## Commands

### `histkit scan`

Parse supported shell history files and update the local index.

`scan` is an ingest, indexing, and reporting command. It does not rewrite shell history.

Typical uses:

```sh
histkit scan
histkit scan --shell bash
histkit scan --config ~/.config/histkit/config.toml
```

### `histkit clean`

Preview or apply built-in cleanup actions for shell history.

`clean` has two operational modes:

- Planning mode: `histkit clean` or `histkit clean --dry-run`
- Apply mode: `histkit clean --apply`

`--dry-run` renders the planned actions without changing files. `--apply` rewrites the detected history source, creates a backup under the histkit state directory, and appends an audit record. Apply mode requires `backup_history = true` in config.

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
- `systemd --user` unit visibility when histkit automation is installed

```sh
histkit doctor
```

### `histkit stats`

Print local index statistics.

Typical output areas include:

- Indexed history entries
- Per-shell counts
- Per-source counts

```sh
histkit stats
```

### `histkit restore`

Restore history state from a recorded backup.

If no backup identifier is given, `restore` lists available backups. When a backup ID is provided, histkit restores that snapshot and appends a restore record to the audit log. As with the other commands, flags must appear before positional arguments.

```sh
histkit restore
histkit restore --config ~/.config/histkit/config.toml
histkit restore b_20260501T130200Z_001
histkit restore --config ~/.config/histkit/config.toml b_20260501T130200Z_001
```

## Available flags

- `--config <path>`: supported by `scan`, `clean`, `pick`, `doctor`, `stats`, and `restore`
- `--shell <name>`: supported by `scan` and `clean`; current values are `bash` and `zsh`
- `--apply`: supported by `clean`
- `--dry-run`: supported by `clean`
- `--help` or `-h`: supported by every command

## Cleanup model

`histkit` supports four primary actions for matched history entries:

- `keep`: leave the command unchanged
- `redact`: preserve the command shape but remove or mask sensitive values
- `delete`: remove the command from the cleanup result
- `quarantine`: rewrite the command to a `[QUARANTINED]` placeholder during apply

Rules may be based on exact matching, keyword groups, regular expressions, structured detectors, or heuristic checks.

For many security-sensitive detections, `redact` or `quarantine` is preferred over immediate deletion.

Typical matched material includes:

- Pasted private key blocks
- Inline passwords
- Bearer tokens
- Cloud credentials
- URL-embedded credentials
- Large accidental paste blobs

Broad rules against commands such as `ssh`, `sudo`, `openssl`, or `kubectl` should be avoided because they generate false positives and reduce trust.

Dry-run output shows the matched rule, confidence, selected action, and a preview of the rewritten command where applicable.

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

The current CLI reads snippets from the configured snippet file and builtin snippet set for `histkit pick`. Snippet management commands are not part of the current command surface.

## Shell integration

`histkit` is designed to work with lightweight shell wrappers for interactive use.

A wrapper function can:

1. Invoke `histkit pick`
2. Capture the selected command
3. Place the result into the current shell editing buffer

The standalone binary writes the selected command to standard output and does not require direct manipulation of shell line-editor internals.

Example Bash setup:

```sh
source /path/to/histkit/contrib/histkit.bash
histkit_bind_bash_pick
```

Pass an alternate readline key sequence to use a different binding:

```sh
histkit_bind_bash_pick '\C-x\C-r'
```

Example Zsh setup:

```sh
source /path/to/histkit/contrib/histkit.zsh
histkit_bind_zsh_pick
```

Pass an alternate ZLE key sequence to use a different binding:

```sh
histkit_bind_zsh_pick '^X^R'
```

Both wrappers default to `Ctrl-R`, invoke `histkit pick`, capture the selected command, and replace the current shell editing buffer with that command. Run `histkit scan` first so the picker has indexed history to display.

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
~/.local/share/histkit/backups/
~/.local/share/histkit/audit.log
~/.cache/histkit/
```

Example configuration:

```toml
[general]
default_shell = "bash"
backup_history = true
dry_run = true
preview_diff = true

[snippets]
enabled = true
builtin = true
user_file = "~/.local/share/histkit/snippets.toml"
```

## Files

```text
~/.config/histkit/config.toml
~/.local/share/histkit/history.db
~/.local/share/histkit/snippets.toml
~/.local/share/histkit/backups/
~/.local/share/histkit/audit.log
~/.config/systemd/user/histkit-scan.service
~/.config/systemd/user/histkit-scan.timer
```

## Examples

### Check the environment before first use

```sh
histkit doctor
```

### Scan history and update the index

```sh
histkit scan
```

### Review indexed history counts

```sh
histkit stats
```

### Show a non-destructive cleanup preview

```sh
histkit clean --dry-run
```

### Apply configured cleanup rules

```sh
histkit clean --apply
```

### Pick a command through `fzf`

```sh
histkit pick
```

### List available backups

```sh
histkit restore
```

### Restore from a specific backup

```sh
histkit restore b_20260501T130200Z_001
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
- Use dry-run mode before enabling automation
- Restrict file permissions on configuration, state, and backup directories

Because active shells may rewrite history files on exit, in-place cleanup should be approached carefully. Derived indexes, reviewable cleanup plans, backups, and audit logs are safer than blind direct rewrites.

## systemd integration

A typical user-level timer should default to periodic scans or report generation. Automated cleanup should remain conservative and should not imply unattended destructive apply mode by default.

The repository ships matching `systemd --user` timer and service templates for scheduled scanning:

```ini
[Unit]
Description=Run histkit scan periodically

[Timer]
OnBootSec=5m
OnUnitActiveSec=12h
Persistent=true

[Install]
WantedBy=timers.target
```

```ini
[Unit]
Description=Scan and index shell history with histkit

[Service]
Type=oneshot
ExecStart=%h/.local/bin/histkit scan --config %h/.config/histkit/config.toml
```

Install `contrib/histkit-scan.timer` and `contrib/histkit-scan.service` as the target user, not as root, unless your deployment specifically requires otherwise. The shipped timer only runs `histkit scan`; it does not schedule `clean --apply`.

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
- External picker failures when `fzf` is missing or misconfigured
- Incomplete audit or backup retention policies in multi-shell environments

## See also

`bash(1)`, `zsh(1)`, `fzf(1)`, `systemd(1)`, `systemd.timer(5)`, `systemd.service(5)`

## License

Distributed under the terms of the project license.
