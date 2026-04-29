# Go CLI Skill

## Goal

Implement histkit CLI commands in small, testable units.

## Constraints

- Keep command handlers thin.
- Put real logic under `internal/` packages.
- Commands should return errors, not call `os.Exit` deep inside logic.
- User-facing errors should be clear.
- Machine-readable output can be added later; plain text is fine initially.
- Avoid destructive behavior in command placeholders.

## Preferred layout

```text
cmd/histkit/main.go
internal/cli/root.go
internal/cli/scan.go
internal/cli/stats.go
internal/cli/doctor.go
```

## Command expectations

### `histkit scan`

Eventually:

- detect shell history sources
- parse entries
- normalize entries
- update SQLite index
- report findings

Initial placeholder:

- explain that scan is not implemented yet
- perform no mutation

### `histkit stats`

Eventually:

- report indexed entry counts
- report counts by shell/source
- report duplicate reduction

Initial placeholder:

- explain that stats are not available until indexing exists

### `histkit doctor`

Eventually:

- validate config
- verify data directories
- verify history source readability
- verify SQLite access
- verify `fzf`
- verify systemd units when configured

Initial placeholder:

- return basic environment message only

## Testing expectations

- Unit test command routing where practical.
- Keep command behavior deterministic.
- Avoid tests that depend on the user's real shell history.
