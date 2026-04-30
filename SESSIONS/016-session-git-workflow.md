# Session 016: Session Git Workflow

## Objective

Add an explicit git/GitHub workflow requirement to every session.

## Completed

- Updated `AGENT.md` so sessions now require branch creation at start and add/commit/push/PR/human approval/merge/cleanup at end.
- Updated `SESSION_PROMPT.md` to include the same GitHub workflow expectations in the session-start prompt.
- Updated `docs/IMPLEMENTATION_OPERATING_MODEL.md` so the recommended workflow now includes branch creation, PR creation, and human approval before merge.
- Recorded the workflow change in `SESSION.md` and this completed session note.

## Files changed

- AGENT.md
- SESSION_PROMPT.md
- docs/IMPLEMENTATION_OPERATING_MODEL.md
- SESSION.md
- SESSIONS/016-session-git-workflow.md

## Files read

- AGENT.md
- SESSION_PROMPT.md
- docs/IMPLEMENTATION_OPERATING_MODEL.md
- SESSION.md
- ROADMAP.md

## Tests added

- None.

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git checkout -b 016-session-git-workflow
sed -n '1,260p' AGENT.md
sed -n '1,220p' SESSION_PROMPT.md
sed -n '1,220p' docs/IMPLEMENTATION_OPERATING_MODEL.md
sed -n '1,220p' SESSION.md
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Every session must start by creating or switching to the session branch before implementation work begins.
- Every session must end by staging intended files, committing with a human-readable message, pushing, and opening a pull request.
- Merge and cleanup require explicit human approval and mark the official end of a session.

## Assumptions

- Documentation changes are sufficient for this workflow requirement without additional enforcement code.

## Known issues

- The workflow remains convention-driven rather than automatically enforced.

## Risks reduced

- Session publishing and closure expectations are now explicit and consistent across the repo’s workflow documents.

## Next recommended session

`016-fzf-picker`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
