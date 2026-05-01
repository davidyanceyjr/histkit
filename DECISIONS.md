# DECISIONS.md

## 001: Go implementation

histkit is implemented in Go.

## 002: Non-destructive first

Initial releases parse, index, preview, and pick. They do not rewrite shell history by default.

## 003: Separate snippets from history

Snippets are never inserted into real shell history by histkit itself.

## 004: SQLite as index

SQLite is the local index and metadata store, not the raw source of truth.

## 005: fzf as picker

histkit uses external `fzf` rather than implementing a fuzzy finder.

## 006: systemd --user for automation

Automation is user-level and optional. Default automation runs `scan`, not destructive apply.

## 007: SESSION.md controls current work

The implementation agent must treat `SESSION.md` as the active work contract.

## 008: Skills are loaded narrowly

The implementation agent should read only the skill files listed in `SESSION.md`.

## 009: Delete actions remove matching history lines during apply

For `histkit clean --apply`, rules whose action is `delete` remove matching history lines from the rewritten history file. This behavior is guarded by per-run backups, atomic rewrite, and audit logging.
