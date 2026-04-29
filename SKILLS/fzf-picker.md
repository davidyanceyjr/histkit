# fzf Picker Skill

## Goal

Return a selected command from indexed history and snippets.

## Constraints

- Do not build an internal fuzzy finder.
- Use external `fzf`.
- Emit selected command to stdout only.
- Keep snippets separate from history.
- Merge history and snippets only at presentation time.

## Candidate labels

```text
[history]
[snippet]
```

## Candidate format

```text
[history]  find . -type f -name '*.tmp' -delete
[snippet]  find {{path}} -type f -name '{{pattern}}' -exec {{cmd}} {} \;
```

## Preview pane

Eventually show:

- full command
- source type
- tags
- last used time
- safety level
- snippet description

## Required tests

- candidate formatting
- selected line parsing
- no `fzf` installed behavior
- snippets are not written to real history
