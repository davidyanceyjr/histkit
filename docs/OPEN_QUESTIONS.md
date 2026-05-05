# OPEN_QUESTIONS.md

## Purpose

This file stores open questions that remain unresolved after a session ends.

Every open question must be answered or documented if unanswered.

## BLOCKING

No blocking questions currently recorded.

## NON-BLOCKING

## Answered / historical

Resolved questions may be moved here from session files when they remain useful for project history.

#### Q001: Should initial source detection support only canonical `~/.bash_history` and `~/.zsh_history` paths?

- Area: source detection defaults
- Initial answer in slice `006`: yes, canonical paths only
- Updated answer: no. Detection now also honors `HISTFILE` for the active `bash` or `zsh` shell while keeping canonical paths as fallback candidates for the other supported shell.
- Why the answer changed: real-system testing showed canonical-only detection misses valid shell history locations when users customize `HISTFILE`.
