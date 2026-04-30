package snippets

import "testing"

func TestSnippetValidate(t *testing.T) {
	snippet := Snippet{
		ID:          "find-delete-pyc",
		Title:       "Delete Python cache files",
		Command:     "find {{path}} -type f -name '*.pyc' -delete",
		Description: "Delete .pyc files under a path",
		Tags:        []string{"find", "python", "cleanup"},
		Placeholders: map[string]string{
			"path": ".",
		},
		Shells: []string{"bash", "zsh"},
		Safety: SafetyMedium,
	}

	if err := snippet.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
	if snippet.Command != "find {{path}} -type f -name '*.pyc' -delete" {
		t.Fatalf("Command = %q, want exact template preservation", snippet.Command)
	}
}

func TestSnippetValidateRequiresFields(t *testing.T) {
	tests := []struct {
		name    string
		snippet Snippet
	}{
		{
			name: "missing id",
			snippet: Snippet{
				Title:       "Title",
				Command:     "echo hi",
				Description: "desc",
				Safety:      SafetyLow,
			},
		},
		{
			name: "missing title",
			snippet: Snippet{
				ID:          "snippet-001",
				Command:     "echo hi",
				Description: "desc",
				Safety:      SafetyLow,
			},
		},
		{
			name: "missing command",
			snippet: Snippet{
				ID:          "snippet-001",
				Title:       "Title",
				Description: "desc",
				Safety:      SafetyLow,
			},
		},
		{
			name: "missing description",
			snippet: Snippet{
				ID:      "snippet-001",
				Title:   "Title",
				Command: "echo hi",
				Safety:  SafetyLow,
			},
		},
		{
			name: "missing safety",
			snippet: Snippet{
				ID:          "snippet-001",
				Title:       "Title",
				Command:     "echo hi",
				Description: "desc",
			},
		},
		{
			name: "invalid safety",
			snippet: Snippet{
				ID:          "snippet-001",
				Title:       "Title",
				Command:     "echo hi",
				Description: "desc",
				Safety:      "critical",
			},
		},
		{
			name: "empty tag",
			snippet: Snippet{
				ID:          "snippet-001",
				Title:       "Title",
				Command:     "echo hi",
				Description: "desc",
				Tags:        []string{"shell", " "},
				Safety:      SafetyLow,
			},
		},
		{
			name: "empty shell",
			snippet: Snippet{
				ID:          "snippet-001",
				Title:       "Title",
				Command:     "echo hi",
				Description: "desc",
				Shells:      []string{"bash", ""},
				Safety:      SafetyLow,
			},
		},
		{
			name: "empty placeholder key",
			snippet: Snippet{
				ID:          "snippet-001",
				Title:       "Title",
				Command:     "echo hi",
				Description: "desc",
				Placeholders: map[string]string{
					"": "value",
				},
				Safety: SafetyLow,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.snippet.Validate(); err == nil {
				t.Fatal("Validate returned nil error")
			}
		})
	}
}

func TestValidateCollectionRejectsDuplicateIDs(t *testing.T) {
	snippets := []Snippet{
		{
			ID:          "dup-id",
			Title:       "One",
			Command:     "echo one",
			Description: "first",
			Safety:      SafetyLow,
		},
		{
			ID:          "dup-id",
			Title:       "Two",
			Command:     "echo two",
			Description: "second",
			Safety:      SafetyMedium,
		},
	}

	if err := ValidateCollection(snippets); err == nil {
		t.Fatal("ValidateCollection returned nil error for duplicate ids")
	}
}

func TestValidateCollectionAcceptsDistinctSnippets(t *testing.T) {
	snippets := []Snippet{
		{
			ID:          "find-delete-pyc",
			Title:       "Delete Python cache files",
			Command:     "find {{path}} -type f -name '*.pyc' -delete",
			Description: "Delete .pyc files under a path",
			Safety:      SafetyMedium,
		},
		{
			ID:          "git-clean-branches",
			Title:       "Delete merged branches",
			Command:     "git branch --merged | grep -v '\\*\\|main\\|master' | xargs -r git branch -d",
			Description: "Delete local branches already merged",
			Safety:      SafetyHigh,
		},
	}

	if err := ValidateCollection(snippets); err != nil {
		t.Fatalf("ValidateCollection returned error: %v", err)
	}
}
