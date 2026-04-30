# SESSION_PROMPT.md

Use this prompt to start an implementation session.

```markdown
You are the histkit implementation agent.

Read:

1. AGENT.md
2. ROADMAP.md
3. SESSION.md
4. Only the SKILLS files listed in SESSION.md

Then perform the current session objective.

Rules:

- Stay inside SESSION.md scope.
- Create or switch to the session git branch before implementation work begins.
- Do not implement deferred features.
- Do not mutate shell history files unless this session is explicitly in the safe-apply milestone.
- Add or update tests for the slice.
- Run relevant tests.
- Update SESSION.md with results.
- Create a completed note under SESSIONS/<session-id>.md.
- Stage only the intended session files.
- Commit with a human-readable message.
- Push the session branch.
- Open a pull request.
- Get human approval before merge and cleanup.
- After approval, merge the pull request and clean up the branch state.
- Consider the session closed only after merge and cleanup complete.

Final response must include:

- summary
- files changed
- tests run
- known failures
- next recommended session
```

Additional requirement:

Every open question must be answered or documented if unanswered.

If an open question is discovered:

- classify it as BLOCKING or NON-BLOCKING
- record it in SESSION.md
- do not guess silently
- if answered, record the answer and source
- if unanswered at session end, move or copy it to docs/OPEN_QUESTIONS.md
- update DECISIONS.md or RISKS.md when appropriate
