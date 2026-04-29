# History Parsing Skill

## Goal

Parse Bash and Zsh history into normalized `HistoryEntry` values.

## Constraints

- Never mutate history files in parser code.
- Preserve the raw line.
- Extract command text.
- Extract timestamp only when reliable.
- Do not discard malformed lines silently; report parse warnings.
- Do not assume shell history is stable while reading.

## Parser interface expectation

The parser layer should support shell-specific implementations behind a common interface.

Conceptually:

```text
HistorySource
  - detect()
  - parse()
  - serialize()
  - lock_strategy()
```

Serialization can be deferred until safe-apply work.

## Bash

Support plain line-oriented history first.

Deferred:

- `HISTTIMEFORMAT` ambiguity
- multiline reconstruction
- shell-session provenance

## Zsh

Support extended history format:

```text
: 1712959000:0;command
```

Parser must split metadata from command correctly.

## Required tests

- plain Bash lines
- empty lines
- Zsh extended history
- malformed Zsh prefix
- commands containing semicolons
- commands with leading/trailing whitespace
