# Testing Skill

## Goal

Keep every implementation slice verifiable.

## General rules

- Add tests with the feature, not later.
- Prefer deterministic fixtures.
- Do not test against the user's real history files.
- Do not require `fzf`, systemd, or a real shell session in unit tests.
- Use integration tests only when they can run safely in temporary directories.

## Go test command

```bash
go test ./...
```

## Fixture guidance

Use project-local fixtures such as:

```text
testdata/history/bash/plain.hist
testdata/history/zsh/extended.hist
```

## Required session ending

Every session note should include:

- tests added
- tests run
- pass/fail result
- known untested behavior
