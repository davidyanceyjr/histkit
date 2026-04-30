# Session 013: Snippet Store

## Objective

Implement the initial snippet store over a TOML user file.

## Completed

- Added a TOML-backed snippet store with load/save/list/add/remove operations.
- Treated a missing user snippet file as an empty store for first-run usability.
- Added snippet config defaults for the user snippet file path and default path tracking in `config.Paths`.
- Preserved exact command templates through round-trip storage.
- Added deterministic tests for round-trip storage, duplicate-ID rejection, and snippet add/remove behavior.

## Files changed

- SESSION.md
- internal/config/config.go
- internal/config/config_test.go
- internal/snippets/store.go
- internal/snippets/store_test.go
- SESSIONS/013-snippet-store.md

## Files read

- AGENT.md
- ROADMAP.md
- SESSION.md
- SKILLS/testing.md
- SKILLS/snippets.md
- SKILLS/config.md
- README.md
- docs/histkit-implementation-plan.md
- internal/config/config.go
- internal/config/config_test.go
- internal/snippets/model.go
- internal/snippets/model_test.go

## Tests added

- TestStoreListMissingFileReturnsEmpty
- TestStoreSaveAndListRoundTrip
- TestStoreListRejectsDuplicateIDs
- TestStoreAddAndRemove
- TestStoreRemoveMissingIDFails

## Tests run

```bash
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Results

All tests passed.

## Commands run

```bash
git status --short --branch
git checkout -b 013-snippet-store
sed -n '1,220p' SESSION.md
sed -n '1,220p' ROADMAP.md
sed -n '1,220p' SKILLS/snippets.md
sed -n '635,715p' docs/histkit-implementation-plan.md
sed -n '390,430p' README.md
sed -n '1,220p' internal/config/config.go
sed -n '1,220p' internal/snippets/model.go
sed -n '1,220p' internal/snippets/model_test.go
rg -n "snippet store|snippets.toml|snippets list|snippets add|snippets remove|user_file|parse snippets TOML|search/list snippets" docs/histkit-implementation-plan.md README.md SKILLS
sed -n '618,680p' docs/histkit-implementation-plan.md
sed -n '190,225p' README.md
sed -n '1,220p' internal/config/config_test.go
sed -n '1,220p' SKILLS/config.md
gofmt -w internal/config/config.go internal/config/config_test.go internal/snippets/store.go internal/snippets/store_test.go
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/go-mod GOPATH=$(pwd)/.cache/go-path go test ./...
```

## Decisions

- Treat a missing user snippet file as an empty store rather than an error.
- Persist snippets as `[[snippets]]` TOML records using the existing model validation rules.
- Keep add/remove operations file-backed and in-memory ordered; do not introduce search or sorting semantics yet.

## Assumptions

- A single user TOML file is sufficient for the first snippet-store slice before builtin catalogs are introduced.

## Known issues

- Builtin snippet catalogs are still out of scope.
- Snippet CLI surfaces are still out of scope.
- Picker integration is still out of scope.

## Risks reduced

- Snippet persistence now exists in a store separate from shell history and SQLite index data.

## Next recommended session

`014-builtin-snippets`

## Open questions discovered

### BLOCKING

None.

### NON-BLOCKING

None.

## Questions answered

None.

## Questions moved to docs/OPEN_QUESTIONS.md

None.
