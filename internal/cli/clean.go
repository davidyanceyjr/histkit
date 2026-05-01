package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"histkit/internal/audit"
	"histkit/internal/backup"
	"histkit/internal/config"
	"histkit/internal/history"
	"histkit/internal/sanitize"
)

func runClean(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("clean", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
	shell := fs.String("shell", "", "clean only one shell history source")
	configPath := fs.String("config", "", "path to a histkit config file")
	apply := fs.Bool("apply", false, "apply cleanup changes to history files")
	dryRun := fs.Bool("dry-run", false, "render the cleanup preview without changing files")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("clean: %w", err)
	}
	if *help {
		writeCleanUsage(stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("clean: unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	if *apply && *dryRun {
		return fmt.Errorf("clean: --apply and --dry-run are mutually exclusive")
	}

	return executeClean(stdout, cleanOptions{
		Shell:      *shell,
		ConfigPath: *configPath,
		Apply:      *apply,
	})
}

func writeCleanUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit clean [--apply] [--dry-run] [--shell <shell>] [--config <path>]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Preview or apply built-in shell history cleanup rules with backup, atomic rewrite, and audit logging.")
}

type cleanOptions struct {
	Shell      string
	ConfigPath string
	Apply      bool
}

func executeClean(stdout io.Writer, opts cleanOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("clean: detect home directory: %w", err)
	}

	cfg := config.Default()
	if opts.ConfigPath != "" {
		expandedConfigPath, err := config.ExpandUserPath(opts.ConfigPath, home)
		if err != nil {
			return fmt.Errorf("clean: %w", err)
		}
		cfg, err = config.Load(expandedConfigPath)
		if err != nil {
			return fmt.Errorf("clean: %w", err)
		}
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		return fmt.Errorf("clean: %w", err)
	}

	sources, err := history.DetectSources(home, opts.Shell)
	if err != nil {
		return fmt.Errorf("clean: %w", err)
	}
	if len(sources) == 0 {
		fmt.Fprintln(stdout, "clean: no history sources detected")
		return nil
	}

	if opts.Apply {
		if !cfg.General.BackupHistory {
			return fmt.Errorf("clean: apply requires backup_history=true")
		}
		return executeCleanApply(stdout, sources, paths)
	}

	return executeCleanDryRun(stdout, sources)
}

func executeCleanDryRun(stdout io.Writer, sources []history.Source) error {
	for i, source := range sources {
		if i > 0 {
			fmt.Fprintln(stdout)
		}

		entries, warnings, err := readHistoryEntries(source)
		if err != nil {
			return fmt.Errorf("clean: %w", err)
		}

		report, err := sanitize.PreviewEntries(entries)
		if err != nil {
			return fmt.Errorf("clean: %w", err)
		}

		fmt.Fprintf(stdout, "source: shell=%s path=%s\n", source.Shell, source.Path)
		if len(warnings) > 0 {
			fmt.Fprintf(stdout, "parse warnings: %d\n", len(warnings))
		}
		fmt.Fprint(stdout, sanitize.RenderPreviewText(report))
	}

	return nil
}

func executeCleanApply(stdout io.Writer, sources []history.Source, paths config.Paths) error {
	startedAt := time.Now().UTC()
	runID := cleanRunID(startedAt)
	backupDir := filepath.Join(paths.StateDir, "backups")

	for i, source := range sources {
		content, err := os.ReadFile(source.Path)
		if err != nil {
			return fmt.Errorf("clean: read history source %q: %w", source.Path, err)
		}

		report, err := sanitize.ApplyToSource(source, content)
		if err != nil {
			return fmt.Errorf("clean: %w", err)
		}
		if report.MatchedEntries == 0 {
			fmt.Fprintf(stdout, "clean apply: shell=%s source=%s no matching entries\n", source.Shell, source.Path)
			continue
		}

		backupRecord, err := backup.Create(source.Path, backupDir, startedAt, i+1)
		if err != nil {
			return fmt.Errorf("clean: %w", err)
		}
		if err := backup.RewriteAtomic(source.Path, report.RewrittenContent); err != nil {
			return fmt.Errorf("clean: %w", err)
		}

		completedAt := time.Now().UTC()
		record := audit.Record{
			RunID:              runID,
			StartedAt:          startedAt,
			CompletedAt:        completedAt,
			Shell:              source.Shell,
			RuleNames:          report.RuleNames,
			CountsByAction:     report.CountsByAction,
			CountsByConfidence: report.CountsByConfidence,
			BackupID:           backupRecord.ID,
			Apply:              true,
		}
		if err := audit.Append(paths.AuditLog, record); err != nil {
			return fmt.Errorf("clean: %w", err)
		}

		fmt.Fprintf(
			stdout,
			"clean apply: shell=%s source=%s matched=%d rewritten=%d deleted=%d backup=%s audit=%s\n",
			source.Shell,
			source.Path,
			report.MatchedEntries,
			report.RewrittenEntries,
			report.DeletedEntries,
			backupRecord.BackupPath,
			paths.AuditLog,
		)
	}

	return nil
}

func readHistoryEntries(source history.Source) ([]history.HistoryEntry, []history.ParseWarning, error) {
	parser, err := history.ParserForShell(source.Shell)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Open(source.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("open history source %q: %w", source.Path, err)
	}
	defer file.Close()

	entries, warnings, err := parser(source.Path, file)
	if err != nil {
		return nil, nil, err
	}
	return entries, warnings, nil
}

func cleanRunID(startedAt time.Time) string {
	return "run_" + startedAt.UTC().Format("20060102T150405.000000000Z")
}
