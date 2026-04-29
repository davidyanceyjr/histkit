# HUMAN_GATES.md

## Purpose

This document defines how the implementation agent handles unknowns.

The core rule:

```text
Every open question must be answered or documented if unanswered.
```

No unresolved question may silently disappear.

## Classification

Every open question must be classified as either:

- `BLOCKING`
- `NON-BLOCKING`

## Blocking questions

A question is `BLOCKING` when the current work cannot safely continue without an answer.

Use this classification when the answer may affect:

- data loss
- destructive behavior
- security posture
- user-visible defaults
- public CLI behavior
- storage schema compatibility
- audit behavior
- backup behavior
- restore behavior
- external integration behavior
- irreversible implementation choices

### Required behavior

When a blocking question appears, the agent must:

1. Stop the affected work.
2. Record the question in `SESSION.md`.
3. Label it `BLOCKING`.
4. Explain the decision needed.
5. Explain why guessing is unsafe.
6. Ask for a direct answer, source document, or permission to defer.
7. Leave the blocked implementation incomplete or safely stubbed.

## Non-blocking questions

A question is `NON-BLOCKING` only when the current session can continue safely.

The temporary assumption must be:

- safe
- reversible
- cheap to change
- not part of a destructive path
- not part of a public contract unless explicitly marked provisional

### Required behavior

When a non-blocking question appears, the agent must:

1. Record the question in `SESSION.md`.
2. Label it `NON-BLOCKING`.
3. Record the temporary assumption.
4. Explain why continuing is safe.
5. State the reversal cost.
6. Continue only inside that safe assumption.

## Required open-question format

```markdown
## Open questions

### BLOCKING

#### Q001: <question>

- Area:
- Decision needed:
- Why blocking:
- Why guessing is unsafe:
- Needed from human:
- Source requested:
- Status: blocked

### NON-BLOCKING

#### Q002: <question>

- Area:
- Temporary assumption:
- Why non-blocking:
- Why assumption is safe:
- Reversal cost:
- Status: assumed-non-blocking
```

## Required answer format

```markdown
## Answer log

### Q001: <question>

- Status: answered
- Answer:
- Source:
- Date answered:
- Decision file updated: yes/no
- Risk file updated: yes/no
```

## Allowed statuses

Open or historical questions must use one of these statuses:

- `blocked`
- `answered`
- `deferred`
- `assumed-non-blocking`
- `superseded`

## Movement rules

A question may be removed from `SESSION.md` only when:

1. it has been answered and logged in the session answer log, or
2. it has been moved to `docs/OPEN_QUESTIONS.md`, or
3. it has been superseded and the superseding question is documented.

## Durable unanswered questions

If a question remains unanswered at the end of a session, it must be copied to `docs/OPEN_QUESTIONS.md`.

Blocking questions must remain visible until resolved or explicitly deferred by a human.

## Durable answered questions

If an answer changes project policy, architecture, CLI behavior, schema, safety posture, or user-facing semantics, the answer must also be recorded in `DECISIONS.md`.

If the answer reduces or introduces risk, update `RISKS.md`.

## Examples

### Blocking example

```markdown
#### Q001: Should `clean --apply` be allowed to delete entries directly?

- Area: cleanup safety
- Decision needed: whether delete is an allowed apply action
- Why blocking: affects irreversible user data mutation
- Why guessing is unsafe: wrong behavior may destroy user history
- Needed from human: direct answer or cleanup policy document
- Source requested: cleanup safety policy
- Status: blocked
```

### Non-blocking example

```markdown
#### Q002: Should generated config examples include comments?

- Area: configuration UX
- Temporary assumption: yes, include short comments
- Why non-blocking: does not affect runtime behavior or data safety
- Why assumption is safe: comments can be changed without migration
- Reversal cost: low
- Status: assumed-non-blocking
```
