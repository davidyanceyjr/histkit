# Sanitizer Skill

## Goal

Classify risky or noisy history entries.

## Constraints

- Do not delete by default.
- Prefer quarantine or redact for sensitive entries.
- Avoid broad rules that destroy trust.
- Dry-run output must explain matched rule and action.
- Redaction should avoid re-exposing secrets in logs.

## Rule classes

- exact match
- contains
- regex
- keyword group
- heuristic detector

## Supported actions

- keep
- redact
- delete
- quarantine

## Initial detections

- private key block markers
- bearer tokens
- inline password flags
- URLs with embedded credentials
- obvious cloud key patterns
- suspicious high-entropy tokens

## Bad initial detections

Do not start with broad rules like:

- every command containing `ssh`
- every command containing `sudo`
- every command containing `openssl`
- every command containing `kubectl`

## Required tests

- true positives
- false positive guard cases
- redaction output
- action classification
- no mutation in dry-run
