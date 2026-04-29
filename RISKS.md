# RISKS.md

## History rewrite races

Shells may append or overwrite history after histkit runs.

### Mitigation

- Delay mutation.
- Warn about active sessions.
- Require backups before apply.
- Prefer derived index first.

## False positives

Secret detection can quarantine legitimate commands.

### Mitigation

- Start with narrow rules.
- Use dry-run preview.
- Prefer quarantine before delete.
- Add false-positive guard tests.

## Trust loss

One bad cleanup can kill the tool.

### Mitigation

- Dry-run default.
- Explicit apply.
- Audit logs.
- Restore support.
- Plain explanations for every destructive action.

## Config sprawl

Too much config too early will slow implementation.

### Mitigation

- Minimal v1 config.
- Advanced rules later.
- Keep defaults boring and safe.

## Shell edge cases

History formats vary.

### Mitigation

- Bash and Zsh first.
- Fixture-driven parser tests.
- Defer Csh/Tcsh.
- Preserve raw lines.

## Sensitive output leakage

Dry-run, audit, logs, and previews can accidentally re-expose secrets.

### Mitigation

- Redact before display where possible.
- Store raw sensitive values only when absolutely required.
- Keep permissions restrictive.
- Avoid logging full secret-bearing commands unless needed for recoverability.
