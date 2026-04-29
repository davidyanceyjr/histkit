# systemd User Skill

## Goal

Support optional user-level automation.

## Constraints

- Use `systemd --user`.
- Default automation should run `scan`, not destructive `clean`.
- Units belong under `contrib/` first.
- `doctor` should verify unit availability when configured.
- Do not enable automatic `clean --apply` by default.

## Initial service

```ini
[Unit]
Description=Scan and index shell history with histkit

[Service]
Type=oneshot
ExecStart=%h/.local/bin/histkit scan --config %h/.config/histkit/config.toml
```

## Initial timer

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

## Required tests or checks

- unit files are syntactically plausible
- `doctor` can detect missing unit files
- documentation makes automation optional
