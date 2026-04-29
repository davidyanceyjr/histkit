# Config Skill

## Goal

Provide a minimal, predictable TOML configuration for histkit.

## Constraints

- Keep v1 config small.
- Defaults should be safe.
- Dry-run should be the default cleanup behavior.
- Do not require config for basic scan/stats commands if sane defaults exist.
- Path expansion must be explicit and tested.

## Required sections

Eventually:

- `[general]`
- `[cleanup]`
- `[security]`
- `[fzf]`
- `[snippets]`
- `[[rules]]`

## Initial config work

Start with:

```toml
[general]
default_shell = "bash"
backup_history = true
dry_run = true
preview_diff = true
```

## Required tests

- load default config
- load config from path
- reject invalid TOML
- expand user paths predictably
