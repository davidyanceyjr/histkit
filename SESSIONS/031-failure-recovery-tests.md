# 031-failure-recovery-tests

Status: completed

## Summary

Added focused recovery tests covering backup metadata-write failure cleanup, restore checksum-mismatch preservation, and audit-log append failures after successful apply or restore operations.

## Verification

- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./internal/backup ./internal/cli`
- `env GOCACHE=/home/opsman/project_git/histkit/.cache/go-build GOMODCACHE=/home/opsman/project_git/histkit/.cache/go-mod GOPATH=/home/opsman/project_git/histkit/.cache/go-path go test ./...`

## Next step

Proceed to `032-systemd-user-service`.
