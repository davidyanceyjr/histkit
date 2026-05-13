# 056-ci-gosec-scan

## Summary

- Added a `gosec` stage to the GitHub Actions workflow using a Go toolchain install.
- Scoped the scan to histkit's own packages and disabled test scanning so local cache trees and fixture-heavy test code do not drown out actionable findings.
- Tightened two user-local directory creation sites from `0o755` to `0o700` so the new scan passes without suppressing those permission findings.
- Pushed branch `056-ci-gosec-scan` and opened draft PR `#52`.

## Objective completed or not completed

- Completed.

## Files read

- `SESSION.md`: prior session state and carry-forward structure
- `ROADMAP.md`: roadmap boundary confirmation for the slice
- `SKILLS/testing.md`: verification expectations for CI-affecting changes
- `SKILLS/go-cli.md`: implementation/testing constraints for the CLI project
- `.github/workflows/ci.yml`: existing CI shape and insertion point for the scan stage
- `SESSIONS/049-basic-ci-workflow.md`: prior CI workflow decisions and conventions
- `DECISIONS.md`: durable project decisions before adding a CI policy decision
- `RISKS.md`: current risk register before adding the static-analysis exclusion risk
- `internal/index/schema.go`: existing directory permission for the SQLite parent path
- `internal/snippets/store.go`: existing directory permission for the snippets store parent path
- `internal/picker/fzf.go`: context for the current `G204` finding
- `internal/history/detect.go`: context for histkit's path-driven local filesystem behavior

## Files changed

- `.github/workflows/ci.yml`: installed `gosec` in CI and added a scoped scan step with narrow rule exclusions
- `internal/index/schema.go`: reduced SQLite parent directory creation permissions from `0o755` to `0o700`
- `internal/snippets/store.go`: reduced snippet store parent directory creation permissions from `0o755` to `0o700`
- `DECISIONS.md`: recorded the durable CI scanning policy decision
- `RISKS.md`: recorded the residual risk from the temporary `gosec` exclusions
- `SESSION.md`: replaced the previous session carry-forward with this session's working state
- `SESSIONS/056-ci-gosec-scan.md`: recorded the completed session

## Tests added

- None.

## Tests run

- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`
- `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath /tmp/histkit-gobin/gosec -tests=false -exclude=G204,G304 ./cmd/... ./internal/...`
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build ./cmd/histkit`

## Known failures

- None after switching local Go cache paths to writable temp directories.

## Commands run

- `git status --short --branch`: confirmed the starting worktree state on `main`
- `ls -1 SESSIONS | sort | tail -n 12`: identified the next session number
- `sed -n '1,240p' .github/workflows/ci.yml`: reviewed the existing CI workflow
- `sed -n '1,220p' SKILLS/go-cli.md`: reviewed relevant implementation constraints
- `sed -n '1,240p' SESSION.md`: reviewed the stale carry-forward session state
- `sed -n '1,220p' ROADMAP.md`: confirmed roadmap boundaries
- `git checkout -b 056-ci-gosec-scan`: created the implementation branch
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath GOBIN=/tmp/histkit-gobin go install github.com/securego/gosec/v2/cmd/gosec@latest`: installed `gosec` locally with writable temp Go paths
- `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath gosec ./...`: showed that scanning `./...` was too noisy because repo-local `.cache` trees were treated as source
- `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath gosec -tests=false ./cmd/... ./internal/...`: reduced the scan to histkit packages and surfaced 12 project-local findings
- `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath gosec -tests=false -exclude=G204,G304 ./cmd/... ./internal/...`: confirmed the remaining findings were the two `G301` permission issues addressed in this slice
- `sed -n '1,200p' internal/index/schema.go`: inspected the SQLite parent directory creation path
- `sed -n '1,200p' internal/snippets/store.go`: inspected the snippet store directory creation path
- `sed -n '1,220p' internal/picker/fzf.go`: reviewed the current `fzf` subprocess launch flow
- `sed -n '1,220p' internal/history/detect.go`: reviewed the project's path-driven local history detection behavior
- `gofmt -w internal/index/schema.go internal/snippets/store.go`: formatted the changed Go files
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go test ./...`: passed
- `env PATH=/tmp/histkit-gobin:$PATH GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache GOPATH=/tmp/histkit-gopath /tmp/histkit-gobin/gosec -tests=false -exclude=G204,G304 ./cmd/... ./internal/...`: passed
- `env GOCACHE=/tmp/histkit-gocache GOMODCACHE=/tmp/histkit-gomodcache go build ./cmd/histkit`: passed
- `rm -f histkit`: removed the local build artifact after verification
- `git add .github/workflows/ci.yml DECISIONS.md RISKS.md SESSION.md SESSIONS/056-ci-gosec-scan.md internal/index/schema.go internal/snippets/store.go && git commit -m "Add gosec scan to CI"`: committed the slice
- `git push -u origin 056-ci-gosec-scan`: pushed the branch to GitHub
- GitHub PR create via connector: opened draft PR `#52`

## Decisions made

- Keep the CI change in the existing single Ubuntu workflow rather than adding a separate security job.
- Scan only `./cmd/...` and `./internal/...` so repo-local cache directories do not become CI noise.
- Disable test scanning for this stage because the deliverable is production-focused SAST signal, not fixture-heavy test analysis.
- Exclude `G204` and `G304` for now and document them as a deliberate CI policy rather than burying them in unexplained failures.
- Fix the two `G301` permission findings in code instead of suppressing them.

## Assumptions made

- `NON-BLOCKING`: excluding `G204` and `G304` is safe for this slice because the findings arise from known local `fzf` execution and user-local path access patterns that already sit inside histkit's explicit CLI contract; the reversal cost is low because the CI flags can be tightened independently later.

## Unresolved questions

- None.

## Risks introduced or reduced

- Reduced: CI now runs a repeatable Go SAST stage against histkit's own packages.
- Reduced: user-local SQLite and snippet parent directories are now created with more restrictive permissions.
- Remaining: CI still suppresses `G204` and `G304`, so future hardening work should revisit process execution and path handling explicitly.

## Next slice recommendation

- Review draft PR `#52`, then merge and clean up the branch after human approval.
- After that, add `govulncheck` to the existing workflow so dependency vulnerability reachability is covered alongside the new `gosec` stage.
