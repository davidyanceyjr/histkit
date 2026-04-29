# Implementation Operating Model

## Purpose

This project uses an agent/session/skills model to keep implementation work focused and context-efficient.

## Main idea

```text
ROADMAP.md decides sequence.
SESSION.md decides current work.
SKILLS/*.md decide implementation behavior.
SESSIONS/*.md preserve history.
DECISIONS.md prevents churn.
RISKS.md keeps danger visible.
```

## Recommended workflow

1. Pick the next roadmap slice.
2. Rewrite `SESSION.md` for that slice.
3. Load only the relevant skills.
4. Implement the slice.
5. Run tests.
6. Save a completed session note.
7. Update `SESSION.md`.

## Session sizing

Ideal session:

- one objective
- three to eight files changed
- five to twenty tests added
- no unrelated refactors
- no milestone boundary crossing

## Context discipline

The implementation agent should not load all project documents every time.

For most sessions, load only:

- `AGENT.md`
- `ROADMAP.md`
- `SESSION.md`
- relevant `SKILLS/*.md`
- source files directly involved in the slice

## When to update durable files

Update `DECISIONS.md` only when a durable architectural or product decision is made.

Update `RISKS.md` only when a new meaningful project risk is discovered.

Create a new file under `SESSIONS/` at the end of every implementation session.
