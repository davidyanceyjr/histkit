# Snippets Skill

## Goal

Store reusable command templates separately from real shell history.

## Constraints

- Snippets must not be written to real shell history.
- Snippets must remain distinct from indexed history in storage.
- Snippets may be merged with history only for picker presentation.
- Placeholder expansion can be deferred.

## Snippet fields

- `id`
- `title`
- `command`
- `description`
- `tags`
- `shells`
- `safety`
- optional `placeholders`

## Example

```toml
[[snippets]]
id = "find-delete-pyc"
title = "Delete Python cache files"
command = "find {{path}} -type f -name '*.pyc' -delete"
description = "Delete .pyc files under a path"
tags = ["find", "python", "cleanup"]
shells = ["bash", "zsh"]
safety = "medium"
```

## Required tests

- parse snippets TOML
- validate required fields
- reject duplicate IDs
- search/list snippets
- preserve command template exactly
