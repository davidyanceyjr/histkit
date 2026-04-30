# Session 014: Builtin Snippets

## Objective

Implement the initial builtin snippet catalog.

## Completed

- Added a curated builtin snippet catalog in `internal/snippets`.
- Added an import helper that seeds missing builtin snippets into the user snippet store.
- Preserved user-owned snippets by refusing to overwrite existing snippet IDs during builtin import.
- Kept builtin import idempotent so repeated imports do not duplicate or reorder builtins unexpectedly.
- Added deterministic tests for builtin validation, initial import, non-overwriting behavior, and idempotent import.

## Files changed

- SESSION.md
- internal/snippets/builtin.go
- internal/snippets/builtin_test.go
- SESSIONS/014-builtin-snippets.md

## Files read

- AGENT.md
- ROADMAP.md
- SESSION.md
- SKILLS/snippets.md
- README.md
- docs/histkit-implementation-plan.md
- internal/snippets/model.go
- internal/snippets/store.go
- internal/snippets/store_test.go

## Tests added

- TestBuiltinsValidate
- TestImportBuiltinsSeedsMissingSnippets
- TestImportBuiltinsDoesNotOverwriteExistingSnippet
- TestImportBuiltinsIdempotent

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git status --short --branch
git checkout -b 014-builtin-snippets
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
rg -n "builtin snippets|builtin = true|Builtin snippets|import builtins|builtin catalog|snippets]\" README.md docs/histkit-implementation-plan.md SKILLS internal
sed -n '330,345p' README.md
sed -n '805,825p' docs/histkit-implementation-plan.md
sed -n '1,240p' internal/snippets/store.go
sed -n '1,240p' internal/snippets/store_test.go
gofmt -w internal/snippets/builtin.go internal/snippets/builtin_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Keep builtin snippets as ordinary `Snippet` values validated through the same model rules as user snippets.
- Import builtins only when their IDs are missing from the user store; existing user snippets always win.
- Keep the builtin import helper idempotent so repeated imports do not reorder or duplicate entries.

## Assumptions

- A small curated builtin catalog is sufficient for the first builtin slice before any CLI import surface exists.

## Known issues

- No CLI surface exists yet for listing or importing builtins directly.
- Picker integration is still out of scope.

## Risks reduced

- Builtin snippets now have a deterministic catalog and a safe import path that preserves user-owned entries.

## Next recommended session

`015-pick-candidates`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
