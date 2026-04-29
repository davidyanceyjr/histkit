# Backup and Restore Skill

## Goal

Make destructive cleanup recoverable.

## Constraints

- No `clean --apply` without backup.
- Backups need checksums.
- Restores must be atomic.
- Audit every apply and restore.
- Refuse unsafe overwrite unless forced.
- Never use the user's real history files in tests.

## Backup requirements

- per-run backup directory
- timestamped identifiers
- checksum for integrity
- metadata stored in SQLite and/or audit log
- explicit restore command

## Restore requirements

- enumerate available backups
- validate backup readability
- replace target file atomically
- log restore event
- refuse unsafe overwrite unless forced

## Required tests

- backup creation
- checksum validation
- atomic replace
- restore by backup ID
- forced overwrite behavior
