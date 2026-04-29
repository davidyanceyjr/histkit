package history

import (
	"fmt"
	"strings"
	"time"
)

const (
	ShellBash = "bash"
	ShellZsh  = "zsh"
)

type HistoryEntry struct {
	ID         string
	Shell      string
	SourceFile string
	RawLine    string
	Command    string
	Timestamp  *time.Time
	ExitCode   *int
	SessionID  string
	Hash       string
}

type ParseWarning struct {
	Shell      string
	SourceFile string
	LineNumber int
	RawLine    string
	Message    string
}

func (e HistoryEntry) HasTimestamp() bool {
	return e.Timestamp != nil
}

func (e HistoryEntry) HasExitCode() bool {
	return e.ExitCode != nil
}

func (e HistoryEntry) Validate() error {
	if strings.TrimSpace(e.Shell) == "" {
		return fmt.Errorf("history entry shell is required")
	}
	if strings.TrimSpace(e.SourceFile) == "" {
		return fmt.Errorf("history entry source file is required")
	}
	if e.RawLine == "" {
		return fmt.Errorf("history entry raw line is required")
	}
	if strings.TrimSpace(e.Command) == "" {
		return fmt.Errorf("history entry command is required")
	}

	return nil
}

func (w ParseWarning) Validate() error {
	if strings.TrimSpace(w.Shell) == "" {
		return fmt.Errorf("parse warning shell is required")
	}
	if strings.TrimSpace(w.SourceFile) == "" {
		return fmt.Errorf("parse warning source file is required")
	}
	if w.LineNumber <= 0 {
		return fmt.Errorf("parse warning line number must be positive")
	}
	if w.RawLine == "" {
		return fmt.Errorf("parse warning raw line is required")
	}
	if strings.TrimSpace(w.Message) == "" {
		return fmt.Errorf("parse warning message is required")
	}

	return nil
}
