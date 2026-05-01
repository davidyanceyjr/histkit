package audit

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"histkit/internal/sanitize"
)

type Record struct {
	RunID              string
	StartedAt          time.Time
	CompletedAt        time.Time
	Shell              string
	RuleNames          []string
	CountsByAction     map[sanitize.ActionType]int
	CountsByConfidence map[sanitize.Confidence]int
	BackupID           string
	Apply              bool
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.RunID) == "" {
		return fmt.Errorf("audit record run id is required")
	}
	if r.StartedAt.IsZero() {
		return fmt.Errorf("audit record started time is required")
	}
	if r.CompletedAt.IsZero() {
		return fmt.Errorf("audit record completed time is required")
	}
	if r.CompletedAt.Before(r.StartedAt) {
		return fmt.Errorf("audit record completed time must not be before started time")
	}
	if strings.TrimSpace(r.Shell) == "" {
		return fmt.Errorf("audit record shell is required")
	}
	if len(r.RuleNames) == 0 {
		return fmt.Errorf("audit record rule names are required")
	}

	seenRules := make(map[string]struct{}, len(r.RuleNames))
	for _, name := range r.RuleNames {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			return fmt.Errorf("audit record rule names must not contain empty values")
		}
		if _, ok := seenRules[trimmed]; ok {
			return fmt.Errorf("audit record rule names must be unique")
		}
		seenRules[trimmed] = struct{}{}
	}

	if err := validateActionCounts(r.CountsByAction); err != nil {
		return err
	}
	if err := validateConfidenceCounts(r.CountsByConfidence); err != nil {
		return err
	}

	return nil
}

func RenderLine(record Record) string {
	rules := slices.Clone(record.RuleNames)
	slices.Sort(rules)

	backupID := strings.TrimSpace(record.BackupID)
	if backupID == "" {
		backupID = "-"
	}

	return fmt.Sprintf(
		"run_id=%s started_at=%s completed_at=%s shell=%s apply=%t backup_id=%s rules=%s action_counts=%s confidence_counts=%s",
		record.RunID,
		record.StartedAt.UTC().Format(time.RFC3339),
		record.CompletedAt.UTC().Format(time.RFC3339),
		record.Shell,
		record.Apply,
		backupID,
		strings.Join(rules, ","),
		renderActionCounts(record.CountsByAction),
		renderConfidenceCounts(record.CountsByConfidence),
	)
}

func validateActionCounts(counts map[sanitize.ActionType]int) error {
	for action, count := range counts {
		if !isValidAction(action) {
			return fmt.Errorf("audit record action %q is invalid", action)
		}
		if count < 0 {
			return fmt.Errorf("audit record action count for %q must not be negative", action)
		}
	}
	return nil
}

func validateConfidenceCounts(counts map[sanitize.Confidence]int) error {
	for confidence, count := range counts {
		if !isValidConfidence(confidence) {
			return fmt.Errorf("audit record confidence %q is invalid", confidence)
		}
		if count < 0 {
			return fmt.Errorf("audit record confidence count for %q must not be negative", confidence)
		}
	}
	return nil
}

func renderActionCounts(counts map[sanitize.ActionType]int) string {
	parts := make([]string, 0, len(orderedActions()))
	for _, action := range orderedActions() {
		parts = append(parts, fmt.Sprintf("%s:%d", action, counts[action]))
	}
	return strings.Join(parts, ",")
}

func renderConfidenceCounts(counts map[sanitize.Confidence]int) string {
	parts := make([]string, 0, len(orderedConfidence()))
	for _, confidence := range orderedConfidence() {
		parts = append(parts, fmt.Sprintf("%s:%d", confidence, counts[confidence]))
	}
	return strings.Join(parts, ",")
}

func orderedActions() []sanitize.ActionType {
	return []sanitize.ActionType{
		sanitize.ActionKeep,
		sanitize.ActionDelete,
		sanitize.ActionRedact,
		sanitize.ActionQuarantine,
	}
}

func orderedConfidence() []sanitize.Confidence {
	return []sanitize.Confidence{
		sanitize.ConfidenceLow,
		sanitize.ConfidenceMedium,
		sanitize.ConfidenceHigh,
	}
}

func isValidAction(action sanitize.ActionType) bool {
	switch action {
	case sanitize.ActionKeep, sanitize.ActionDelete, sanitize.ActionRedact, sanitize.ActionQuarantine:
		return true
	default:
		return false
	}
}

func isValidConfidence(confidence sanitize.Confidence) bool {
	switch confidence {
	case sanitize.ConfidenceLow, sanitize.ConfidenceMedium, sanitize.ConfidenceHigh:
		return true
	default:
		return false
	}
}
