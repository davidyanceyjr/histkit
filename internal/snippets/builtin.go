package snippets

import "fmt"

func Builtins() []Snippet {
	return []Snippet{
		{
			ID:          "find-delete-pyc",
			Title:       "Delete Python cache files",
			Command:     "find {{path}} -type f -name '*.pyc' -delete",
			Description: "Delete .pyc files under a target path",
			Tags:        []string{"find", "python", "cleanup"},
			Shells:      []string{"bash", "zsh"},
			Safety:      SafetyMedium,
		},
		{
			ID:          "git-clean-merged-branches",
			Title:       "Delete merged Git branches",
			Command:     "git branch --merged | grep -v '\\*\\|main\\|master' | xargs -r git branch -d",
			Description: "Delete local Git branches that are already merged",
			Tags:        []string{"git", "cleanup"},
			Shells:      []string{"bash", "zsh"},
			Safety:      SafetyHigh,
		},
		{
			ID:          "find-large-files",
			Title:       "Find large files",
			Command:     "find {{path}} -type f -size +{{size}} -print",
			Description: "List files above a chosen size threshold",
			Tags:        []string{"find", "disk", "inspection"},
			Shells:      []string{"bash", "zsh"},
			Safety:      SafetyLow,
		},
	}
}

func ImportBuiltins(store Store) (int, error) {
	if store.Path == "" {
		return 0, fmt.Errorf("snippet store path is required")
	}

	existing, err := store.List()
	if err != nil {
		return 0, err
	}

	seen := make(map[string]struct{}, len(existing))
	for _, snippet := range existing {
		seen[snippet.ID] = struct{}{}
	}

	merged := append([]Snippet{}, existing...)
	imported := 0
	for _, snippet := range Builtins() {
		if _, ok := seen[snippet.ID]; ok {
			continue
		}
		merged = append(merged, snippet)
		seen[snippet.ID] = struct{}{}
		imported++
	}

	if imported == 0 {
		return 0, nil
	}
	if err := store.Save(merged); err != nil {
		return 0, err
	}

	return imported, nil
}
