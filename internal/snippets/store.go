package snippets

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Store struct {
	Path string
}

type storeDocument struct {
	Snippets []Snippet `toml:"snippets"`
}

func (s Store) List() ([]Snippet, error) {
	if s.Path == "" {
		return nil, fmt.Errorf("snippet store path is required")
	}

	if _, err := os.Stat(s.Path); err != nil {
		if os.IsNotExist(err) {
			return []Snippet{}, nil
		}
		return nil, fmt.Errorf("stat snippet store %q: %w", s.Path, err)
	}

	var doc storeDocument
	if _, err := toml.DecodeFile(s.Path, &doc); err != nil {
		return nil, fmt.Errorf("load snippet store %q: %w", s.Path, err)
	}
	if err := ValidateCollection(doc.Snippets); err != nil {
		return nil, fmt.Errorf("load snippet store %q: %w", s.Path, err)
	}

	return doc.Snippets, nil
}

func (s Store) Save(snippets []Snippet) error {
	if s.Path == "" {
		return fmt.Errorf("snippet store path is required")
	}
	if err := ValidateCollection(snippets); err != nil {
		return fmt.Errorf("save snippet store %q: %w", s.Path, err)
	}
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return fmt.Errorf("save snippet store %q: %w", s.Path, err)
	}

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(storeDocument{Snippets: snippets}); err != nil {
		return fmt.Errorf("save snippet store %q: %w", s.Path, err)
	}
	if err := os.WriteFile(s.Path, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("save snippet store %q: %w", s.Path, err)
	}

	return nil
}

func (s Store) Add(snippet Snippet) error {
	snippets, err := s.List()
	if err != nil {
		return err
	}
	snippets = append(snippets, snippet)
	return s.Save(snippets)
}

func (s Store) Remove(id string) error {
	if id == "" {
		return fmt.Errorf("snippet id is required")
	}

	snippets, err := s.List()
	if err != nil {
		return err
	}

	filtered := snippets[:0]
	removed := false
	for _, snippet := range snippets {
		if snippet.ID == id {
			removed = true
			continue
		}
		filtered = append(filtered, snippet)
	}
	if !removed {
		return fmt.Errorf("snippet %q not found", id)
	}

	return s.Save(filtered)
}
