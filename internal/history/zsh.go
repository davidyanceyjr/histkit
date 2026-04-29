package history

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseZsh(sourceFile string, r io.Reader) ([]HistoryEntry, []ParseWarning, error) {
	if strings.TrimSpace(sourceFile) == "" {
		return nil, nil, fmt.Errorf("zsh parser source file is required")
	}
	if r == nil {
		return nil, nil, fmt.Errorf("zsh parser reader is required")
	}

	scanner := bufio.NewScanner(r)
	var entries []HistoryEntry
	var warnings []ParseWarning
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		rawLine := scanner.Text()

		switch {
		case rawLine == "":
			continue
		case strings.TrimSpace(rawLine) == "":
			warnings = append(warnings, ParseWarning{
				Shell:      ShellZsh,
				SourceFile: sourceFile,
				LineNumber: lineNumber,
				RawLine:    rawLine,
				Message:    "whitespace-only Zsh history line",
			})
		case strings.HasPrefix(rawLine, ": "):
			entry, warning := parseZshExtendedLine(sourceFile, lineNumber, rawLine)
			if warning != nil {
				warnings = append(warnings, *warning)
				continue
			}
			entries = append(entries, entry)
		default:
			entries = append(entries, HistoryEntry{
				Shell:      ShellZsh,
				SourceFile: sourceFile,
				RawLine:    rawLine,
				Command:    rawLine,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("read Zsh history from %q: %w", sourceFile, err)
	}

	return entries, warnings, nil
}

func parseZshExtendedLine(sourceFile string, lineNumber int, rawLine string) (HistoryEntry, *ParseWarning) {
	metadataAndCommand := strings.TrimPrefix(rawLine, ": ")
	parts := strings.SplitN(metadataAndCommand, ";", 2)
	if len(parts) != 2 {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: missing command separator",
		}
	}

	command := parts[1]
	if strings.TrimSpace(command) == "" {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: empty command",
		}
	}

	metadata := strings.SplitN(parts[0], ":", 2)
	if len(metadata) != 2 {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: invalid metadata fields",
		}
	}

	unixSeconds, err := strconv.ParseInt(metadata[0], 10, 64)
	if err != nil {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: invalid timestamp",
		}
	}

	if _, err := strconv.Atoi(metadata[1]); err != nil {
		return HistoryEntry{}, &ParseWarning{
			Shell:      ShellZsh,
			SourceFile: sourceFile,
			LineNumber: lineNumber,
			RawLine:    rawLine,
			Message:    "malformed Zsh extended history line: invalid duration",
		}
	}

	timestamp := time.Unix(unixSeconds, 0).UTC()

	return HistoryEntry{
		Shell:      ShellZsh,
		SourceFile: sourceFile,
		RawLine:    rawLine,
		Command:    command,
		Timestamp:  &timestamp,
	}, nil
}
