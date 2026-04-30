package snippets

import (
	"fmt"
	"strings"
)

const (
	SafetyLow    = "low"
	SafetyMedium = "medium"
	SafetyHigh   = "high"
)

type Snippet struct {
	ID           string
	Title        string
	Command      string
	Description  string
	Tags         []string
	Placeholders map[string]string
	Shells       []string
	Safety       string
}

func (s Snippet) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("snippet id is required")
	}
	if strings.TrimSpace(s.Title) == "" {
		return fmt.Errorf("snippet title is required")
	}
	if strings.TrimSpace(s.Command) == "" {
		return fmt.Errorf("snippet command is required")
	}
	if strings.TrimSpace(s.Description) == "" {
		return fmt.Errorf("snippet description is required")
	}
	if strings.TrimSpace(s.Safety) == "" {
		return fmt.Errorf("snippet safety is required")
	}
	if !isValidSafety(s.Safety) {
		return fmt.Errorf("snippet safety %q is invalid", s.Safety)
	}

	for _, tag := range s.Tags {
		if strings.TrimSpace(tag) == "" {
			return fmt.Errorf("snippet tags must not contain empty values")
		}
	}
	for _, shell := range s.Shells {
		if strings.TrimSpace(shell) == "" {
			return fmt.Errorf("snippet shells must not contain empty values")
		}
	}
	for key := range s.Placeholders {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("snippet placeholders must not contain empty keys")
		}
	}

	return nil
}

func ValidateCollection(snippets []Snippet) error {
	seen := make(map[string]struct{}, len(snippets))
	for _, snippet := range snippets {
		if err := snippet.Validate(); err != nil {
			return err
		}
		if _, ok := seen[snippet.ID]; ok {
			return fmt.Errorf("duplicate snippet id %q", snippet.ID)
		}
		seen[snippet.ID] = struct{}{}
	}

	return nil
}

func isValidSafety(safety string) bool {
	switch safety {
	case SafetyLow, SafetyMedium, SafetyHigh:
		return true
	default:
		return false
	}
}
