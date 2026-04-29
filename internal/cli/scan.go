package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"histkit/internal/config"
	"histkit/internal/history"
	"histkit/internal/index"
)

func runScan(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
	shell := fs.String("shell", "", "scan only one shell history source")
	configPath := fs.String("config", "", "path to a histkit config file")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	if *help {
		writeScanUsage(stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("scan: unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	return executeScan(stdout, scanOptions{
		Shell:      *shell,
		ConfigPath: *configPath,
	})
}

func writeScanUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit scan [--shell <shell>] [--config <path>]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Parse detected shell history sources, update the local SQLite index, and report a summary.")
}

type scanOptions struct {
	Shell      string
	ConfigPath string
}

type scanSummary struct {
	Sources   int
	Warnings  int
	Attempted int
	Inserted  int
	Skipped   int
}

func executeScan(stdout io.Writer, opts scanOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("scan: detect home directory: %w", err)
	}

	if opts.ConfigPath != "" {
		expandedConfigPath, err := config.ExpandUserPath(opts.ConfigPath, home)
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}
		if _, err := config.Load(expandedConfigPath); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	sources, err := history.DetectSources(home, opts.Shell)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	db, err := index.Open(paths.HistoryDB)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	defer db.Close()

	if err := index.InitSchema(db); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	summary := scanSummary{Sources: len(sources)}
	for _, source := range sources {
		parser, err := history.ParserForShell(source.Shell)
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		file, err := os.Open(source.Path)
		if err != nil {
			return fmt.Errorf("scan: open history source %q: %w", source.Path, err)
		}

		entries, warnings, parseErr := parser(source.Path, file)
		closeErr := file.Close()
		if parseErr != nil {
			return fmt.Errorf("scan: %w", parseErr)
		}
		if closeErr != nil {
			return fmt.Errorf("scan: close history source %q: %w", source.Path, closeErr)
		}

		result, err := index.WriteHistoryEntries(db, entries)
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		summary.Warnings += len(warnings)
		summary.Attempted += result.Attempted
		summary.Inserted += result.Inserted
		summary.Skipped += result.Skipped
	}

	fmt.Fprintf(
		stdout,
		"scan complete: %d source(s), %d entries parsed, %d inserted, %d skipped, %d warning(s). history index: %s\n",
		summary.Sources,
		summary.Attempted,
		summary.Inserted,
		summary.Skipped,
		summary.Warnings,
		paths.HistoryDB,
	)

	return nil
}
