# 058-ci-binary-smoke-checks

## Summary

- Added binary smoke checks to the GitHub Actions workflow after the build step.
- The workflow now builds a named `./histkit` artifact and runs `./histkit --help` plus `./histkit doctor --help`.
- Pushed branch `058-ci-binary-smoke-checks` and opened draft PR `#54`.

## Objective completed or not completed

- Completed.

## Files read

- `SESSION.md`: previous closed session state and carry-forward structure
- `ROADMAP.md`: roadmap boundary confirmation for the slice
- `SKILLS/testing.md`: verification expectations
- `SKILLS/go-cli.md`: implementation constraints relevant to the CLI repository
- `.github/workflows/ci.yml`: existing CI job shape and insertion point for the smoke stage

## Files changed

- `.github/workflows/ci.yml`: updated the build output and added the binary help smoke checks
- `SESSION.md`: replaced the previous closed session carry-forward with this session's working state
- `SESSIONS/058-ci-binary-smoke-checks.md`: recorded the session

## Tests added

- None.

## Tests run

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
- `mkdir -p /tmp/histkit-smoke && env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build -o /tmp/histkit-smoke/histkit ./cmd/histkit && /tmp/histkit-smoke/histkit --help && /tmp/histkit-smoke/histkit doctor --help`
- `git add .github/workflows/ci.yml SESSION.md SESSIONS/058-ci-binary-smoke-checks.md && git commit -m "Add binary smoke checks to CI"`: committed the slice
- `git push -u origin 058-ci-binary-smoke-checks`: pushed the branch to GitHub
- GitHub PR create via connector: opened draft PR `#54`

## Known failures

- None.

## Commands run

- `git status --short --branch`: inspected repository state on `main`
- `sed -n '1,240p' SESSION.md`: reviewed the previous closed session state
- `sed -n '1,220p' ROADMAP.md`: confirmed roadmap boundaries
- `sed -n '1,220p' SKILLS/testing.md`: reviewed verification expectations
- `sed -n '1,220p' .github/workflows/ci.yml`: reviewed the current CI workflow
- `ls -1 SESSIONS | sort | tail -n 4`: identified the next session number
- `git checkout -b 058-ci-binary-smoke-checks`: created the session branch
- `sed -n '1,220p' SKILLS/go-cli.md`: reviewed relevant implementation constraints
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`: passed
- `mkdir -p /tmp/histkit-smoke && env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build -o /tmp/histkit-smoke/histkit ./cmd/histkit && /tmp/histkit-smoke/histkit --help && /tmp/histkit-smoke/histkit doctor --help`: passed

## Decisions made

- Keep the smoke-check change in the existing single Ubuntu workflow.
- Build a named binary and smoke-test the same artifact in CI.

## Assumptions made

- None.

## Unresolved questions

- None.

## Risks introduced or reduced

- Reduced: CI now validates that the built CLI actually starts and prints help output.
- Remaining: smoke checks only cover help paths, not command execution with real state.

## Next slice recommendation

- Review draft PR `#54`, then merge and clean up the branch after human approval.
- Add a small CI step or job if the project later wants to split smoke checks from the main build.
