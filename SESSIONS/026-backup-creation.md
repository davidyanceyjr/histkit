# 026-backup-creation

## Goal

Create real on-disk backup files from source history files, using the existing backup record model and deterministic backup identifiers.

## Scope

- compute source-file checksums
- create backup directories as needed
- copy source files into the derived backup path
- verify copied bytes by checksum
- add focused temp-file tests

## Out of scope

- restore behavior
- atomic rewrite support
- audit logging
- CLI wiring

## Notes

- Backup paths continue to use the existing `<backup-dir>/<backup-id>/<basename(source-file)>` shape.
- This slice should fail rather than overwrite an existing backup artifact.
