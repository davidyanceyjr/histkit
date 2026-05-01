# 029-clean-apply

Status: completed

## Summary

Implemented the first `clean` CLI slice with default dry-run output, explicit apply mode, shell-aware rewrite behavior, per-run backups, atomic file replacement, and audit-log appends.

## Policy answer

`delete` rules remove matching history lines during `clean --apply`.

## Verification

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/sanitize ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

## Next step

Proceed to `030-restore-command`.
