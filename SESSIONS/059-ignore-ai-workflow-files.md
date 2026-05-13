# 059-ignore-ai-workflow-files

## Summary

- Added `.gitignore` rules for common AI assistant workflow metadata files and directories.
- Removed `AGENTS.md` from the git index so it is no longer tracked, while preserving the local file.
- Kept `.github/` tracked because the repository uses it for CI rather than AI workflow metadata.

## Objective completed or not completed

- Completed.

## Files read

- `SESSION.md`: previous closed session state and carry-forward structure
- `ROADMAP.md`: roadmap boundary confirmation for the slice
- `AGENTS.md`: required implementation workflow and session record expectations
- `.gitignore`: existing ignore rules and insertion point for AI metadata patterns
- `/home/opsman/.codex/plugins/cache/openai-curated/github/1141b764/skills/yeet/SKILL.md`: publish workflow expectations for commit, push, and draft PR creation

## Files changed

- `.gitignore`: added ignore rules for AI workflow metadata
- `AGENTS.md`: removed from git tracking while leaving the local file intact
- `SESSION.md`: replaced the previous closed session carry-forward with this session's working state
- `SESSIONS/059-ignore-ai-workflow-files.md`: recorded the session

## Tests added

- None.

## Tests run

- `git check-ignore -v AGENTS.md .claude/settings.json .codex/config.toml CLAUDE.md`

## Known failures

- None.

## Commands run

- `git status --short --branch`: inspected repository state on `main`
- `sed -n '1,220p' SESSION.md`: reviewed the previous closed session state
- `sed -n '1,240p' ROADMAP.md`: confirmed roadmap boundaries
- `sed -n '1,220p' .gitignore`: reviewed current ignore rules
- `sed -n '1,240p' AGENTS.md`: reviewed required session workflow
- `ls -1 SESSIONS | sort | tail -n 8`: identified the next session number
- `git ls-files | rg '(^|/)(\\.github|\\.claude|\\.cursor|\\.copilot|\\.aider|\\.codex|agents|AGENTS\\.md|CLAUDE\\.md|GEMINI\\.md|copilot-instructions\\.md|promptfoo|\\.windsurf)($|/)'`: checked currently tracked workflow-style files
- `find .. -maxdepth 3 \\( -name '.github' -o -name '.claude' -o -name '.cursor' -o -name '.copilot' -o -name '.aider*' -o -name '.codex' -o -name 'CLAUDE.md' -o -name 'GEMINI.md' -o -name 'copilot-instructions.md' -o -name 'AGENTS.md' \\) | sed -n '1,120p'`: checked nearby filesystem matches
- `git checkout -b 059-ignore-ai-workflow-files`: created the session branch
- `sed -n '1,240p' /home/opsman/.codex/plugins/cache/openai-curated/github/1141b764/skills/yeet/SKILL.md`: reviewed publish workflow guidance
- `git remote -v`: confirmed GitHub remote details
- `git rm --cached -- AGENTS.md`: removed `AGENTS.md` from git tracking while keeping the local file
- `git check-ignore -v AGENTS.md .claude/settings.json .codex/config.toml CLAUDE.md`: verified the new ignore patterns
- `git status --short`: confirmed the intended staged deletion and ignore-rule modification

## Decisions made

- Treat AI workflow files as assistant-specific metadata rather than all GitHub workflow content.
- Stop tracking `AGENTS.md` by removing it from the index instead of deleting the local file.

## Assumptions made

- None.

## Unresolved questions

- None.

## Risks introduced or reduced

- Reduced: AI assistant metadata is less likely to be committed accidentally.
- Remaining: future repository-specific agent instructions in `AGENTS.md` will no longer be shared through git unless managed through another channel.

## Next slice recommendation

- Review the draft PR and confirm the ignore scope before merge.
