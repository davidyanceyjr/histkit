# SQLite Skill

## Goal

Store normalized history entries and run metadata.

## Constraints

- Use SQLite as local index.
- Do not treat SQLite as the raw source of truth.
- Schema migrations must be explicit.
- Use restrictive file permissions where applicable.
- Tests must use temporary databases.

## Initial tables

- `history_entries`
- `runs`

## Deferred tables

- `history_actions`
- `snippets`
- `backups`

## Suggested `history_entries` columns

- `id`
- `shell`
- `source_file`
- `raw_line`
- `command`
- `timestamp`
- `exit_code`
- `session_id`
- `hash`
- `ingested_at`

## Suggested `runs` columns

- `id`
- `command`
- `started_at`
- `finished_at`
- `status`
- `notes`

## Required tests

- initialize schema
- insert entries
- dedupe by hash/source
- query stats
- reopen database
