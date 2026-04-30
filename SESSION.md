# SESSION.md

## Current session

ID: `016-session-git-workflow`

Status: completed

## Objective

Add an explicit git/GitHub workflow requirement to every session.

## Scope

Implement:

- workflow updates in `AGENT.md`
- session-start/session-end prompt updates in `SESSION_PROMPT.md`
- operating-model documentation updates
- session record and note updates

## Out of scope

- automation of GitHub publishing
- CI enforcement of session workflow
- changes to product runtime behavior

## Relevant skills

- `SKILLS/testing.md`

## Acceptance criteria

- repository workflow docs require creating the branch at session start
- repository workflow docs require staging, committing, pushing, and PR creation at session end
- repository workflow docs require human approval before merge and cleanup
- `go test ./...` still passes

## Current repo state

The repository now documents a full session-close workflow that includes branch creation, commit/push, PR creation, human approval, merge, and cleanup.

## Decisions already made

- Language: Go
- CLI should be conservative and auditable
- Destructive cleanup is deferred
- SQLite is the local index and metadata store
- `fzf` is the picker engine
- `systemd --user` is the automation target

## Risks to watch

- Keep workflow docs consistent across all session entry points.
- Do not imply automatic merge authority; preserve explicit human approval.

## Open questions

Every open question discovered during this session must be recorded here.

### BLOCKING

No blocking questions currently recorded.

### NON-BLOCKING

No non-blocking questions currently recorded.

## Answer log

Every answered question must be recorded here before it is removed from the active open-question list.

### Answered this session

No questions answered yet.

## End-of-session notes

Summary:

- Updated `AGENT.md` so sessions now require branch creation at start and add/commit/push/PR/human approval/merge/cleanup at end.
- Updated `SESSION_PROMPT.md` to include the same GitHub workflow expectations in the startup prompt.
- Updated `docs/IMPLEMENTATION_OPERATING_MODEL.md` so the recommended workflow now includes branch creation, PR creation, and human approval before merge.

Files changed:

- AGENT.md
- SESSION_PROMPT.md
- docs/IMPLEMENTATION_OPERATING_MODEL.md
- SESSION.md
- SESSIONS/016-session-git-workflow.md

Files read:

- AGENT.md
- SESSION_PROMPT.md
- docs/IMPLEMENTATION_OPERATING_MODEL.md
- SESSION.md
- ROADMAP.md

Tests added:

- None.

Tests run:

- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Known failures:

- None.

Decisions made:

- Every session must start by creating or switching to the session branch before implementation work begins.
- Every session must end by staging intended files, committing with a human-readable message, pushing, and opening a pull request.
- Merge and cleanup require explicit human approval and mark the official end of a session.

Commands run:

- `git checkout -b 016-session-git-workflow`
- `sed -n '1,260p' AGENT.md`
- `sed -n '1,220p' SESSION_PROMPT.md`
- `sed -n '1,220p' docs/IMPLEMENTATION_OPERATING_MODEL.md`
- `sed -n '1,220p' SESSION.md`
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...`

Assumptions made:

- Documentation changes are sufficient for this workflow requirement without additional enforcement code.

Risks introduced or reduced:

- Reduced: session publishing and closure expectations are now explicit and consistent across the repo’s workflow documents.
- Ongoing: the workflow is still convention-driven rather than automatically enforced.

Next recommended session:

- `016-fzf-picker`
