# Session 046: Pick fzf TTY Wiring Plan

## Objective

Define a focused implementation slice to make `histkit pick` drive `fzf` through the controlling terminal instead of buffered process pipes.

## Completed

- Identified `internal/picker/fzf.go` as the primary runtime change surface.
- Reviewed the existing picker tests and shell wrapper contracts.
- Recorded a concrete next slice with file scope, risks, acceptance criteria, and one explicit non-blocking design question.

## Files changed

- `SESSION.md`
- `SESSIONS/046-pick-fzf-tty-wiring-plan.md`

## Tests added

- None.

## Tests run

```bash
# None. Planning-only session.
```

## Results

The next implementation slice is defined as `047-pick-fzf-tty-wiring`. No runtime behavior changed in this planning session.

## Decisions

- Keep the follow-up change narrowly centered on `internal/picker/fzf.go`.
- Preserve the existing stdout contract for the selected command.
- Treat shell wrapper compatibility as a validation target, not as primary implementation scope.

## Known issues

- The `fzf` interactive issue remains unresolved until the implementation slice lands.
- The implementation still needs to choose between inheriting stdio and opening `/dev/tty` explicitly.

## Next recommended session

`047-pick-fzf-tty-wiring`


## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

- Should the implementation prefer inheriting the parent stdio streams directly or opening `/dev/tty` explicitly for `fzf`?

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
