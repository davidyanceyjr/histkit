# 030-restore-command

Status: completed

## Summary

Added persisted backup metadata plus an initial `histkit restore` command that lists backups, restores by backup ID, validates checksum integrity, rewrites history atomically, and appends an audit entry.

## Verification

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

## Next step

Proceed to `031-failure-recovery-tests`.
