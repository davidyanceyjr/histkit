package sanitize

import (
	"fmt"
	"regexp"
	"strings"
)

type RuleType string

const (
	RuleExact        RuleType = "exact"
	RuleContains     RuleType = "contains"
	RuleRegex        RuleType = "regex"
	RuleKeywordGroup RuleType = "keyword_group"
	RuleHeuristic    RuleType = "heuristic"
)

type ActionType string

const (
	ActionKeep       ActionType = "keep"
	ActionDelete     ActionType = "delete"
	ActionRedact     ActionType = "redact"
	ActionQuarantine ActionType = "quarantine"
)

type Confidence string

const (
	ConfidenceLow    Confidence = "low"
	ConfidenceMedium Confidence = "medium"
	ConfidenceHigh   Confidence = "high"
)

type Rule struct {
	Name       string
	Type       RuleType
	Pattern    string
	Keywords   []string
	Detector   string
	Action     ActionType
	Confidence Confidence
	Reason     string
}

type RuleMatch struct {
	RuleName   string
	Reason     string
	Confidence Confidence
	Action     ActionType
	Before     string
	After      string
}

func (r Rule) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("rule name is required")
	}
	if !isValidRuleType(r.Type) {
		return fmt.Errorf("rule type %q is invalid", r.Type)
	}
	if !isValidAction(r.Action) {
		return fmt.Errorf("rule action %q is invalid", r.Action)
	}
	if !isValidConfidence(r.Confidence) {
		return fmt.Errorf("rule confidence %q is invalid", r.Confidence)
	}
	if strings.TrimSpace(r.Reason) == "" {
		return fmt.Errorf("rule reason is required")
	}

	switch r.Type {
	case RuleExact, RuleContains:
		if strings.TrimSpace(r.Pattern) == "" {
			return fmt.Errorf("rule pattern is required for type %q", r.Type)
		}
	case RuleRegex:
		if strings.TrimSpace(r.Pattern) == "" {
			return fmt.Errorf("rule pattern is required for type %q", r.Type)
		}
		if _, err := regexp.Compile(r.Pattern); err != nil {
			return fmt.Errorf("rule regex pattern %q is invalid: %w", r.Pattern, err)
		}
	case RuleKeywordGroup:
		if len(r.Keywords) == 0 {
			return fmt.Errorf("rule keywords are required for type %q", r.Type)
		}
		for _, keyword := range r.Keywords {
			if strings.TrimSpace(keyword) == "" {
				return fmt.Errorf("rule keywords must not contain empty values")
			}
		}
	case RuleHeuristic:
		if strings.TrimSpace(r.Detector) == "" {
			return fmt.Errorf("rule detector is required for type %q", r.Type)
		}
	}

	return nil
}

func (m RuleMatch) Validate() error {
	if strings.TrimSpace(m.RuleName) == "" {
		return fmt.Errorf("rule match name is required")
	}
	if strings.TrimSpace(m.Reason) == "" {
		return fmt.Errorf("rule match reason is required")
	}
	if !isValidConfidence(m.Confidence) {
		return fmt.Errorf("rule match confidence %q is invalid", m.Confidence)
	}
	if !isValidAction(m.Action) {
		return fmt.Errorf("rule match action %q is invalid", m.Action)
	}
	if strings.TrimSpace(m.Before) == "" {
		return fmt.Errorf("rule match before value is required")
	}
	if m.Action == ActionRedact && strings.TrimSpace(m.After) == "" {
		return fmt.Errorf("rule match after value is required for action %q", m.Action)
	}

	return nil
}

func ValidateRules(rules []Rule) error {
	seen := make(map[string]struct{}, len(rules))
	for _, rule := range rules {
		if err := rule.Validate(); err != nil {
			return err
		}
		if _, ok := seen[rule.Name]; ok {
			return fmt.Errorf("duplicate rule name %q", rule.Name)
		}
		seen[rule.Name] = struct{}{}
	}

	return nil
}

func isValidRuleType(ruleType RuleType) bool {
	switch ruleType {
	case RuleExact, RuleContains, RuleRegex, RuleKeywordGroup, RuleHeuristic:
		return true
	default:
		return false
	}
}

func isValidAction(action ActionType) bool {
	switch action {
	case ActionKeep, ActionDelete, ActionRedact, ActionQuarantine:
		return true
	default:
		return false
	}
}

func isValidConfidence(confidence Confidence) bool {
	switch confidence {
	case ConfidenceLow, ConfidenceMedium, ConfidenceHigh:
		return true
	default:
		return false
	}
}
