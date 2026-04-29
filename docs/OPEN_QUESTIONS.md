# OPEN_QUESTIONS.md

## Purpose

This file stores open questions that remain unresolved after a session ends.

Every open question must be answered or documented if unanswered.

## BLOCKING

No blocking questions currently recorded.

## NON-BLOCKING

#### Q001: Should initial source detection support only canonical `~/.bash_history` and `~/.zsh_history` paths?

- Area: source detection defaults
- Temporary assumption: yes, detect only canonical Bash and Zsh history paths in slice `006`
- Why non-blocking: configurable/custom history path support can be added later without changing parser correctness
- Why assumption is safe: it narrows behavior conservatively and does not mutate data
- Reversal cost: low
- Status: assumed-non-blocking

## Answered / historical

Resolved questions may be moved here from session files when they remain useful for project history.

No historical questions currently recorded.
