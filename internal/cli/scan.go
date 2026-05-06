package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/davidyanceyjr/histkit/internal/config"
	"github.com/davidyanceyjr/histkit/internal/history"
	"github.com/davidyanceyjr/histkit/internal/index"
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
	writeHelpBlocks(w,
		helpBlock{
			Title: "Usage",
			Lines: []string{
				"  histkit scan [--shell <shell>] [--config <path>]",
			},
		},
		helpBlock{
			Lines: []string{
				"Parse detected shell history sources, refresh the local SQLite history index, and report what was indexed.",
			},
		},
		helpBlock{
			Lines: []string{
				"scan reads supported history files and writes normalized entries into the local index.",
				"It does not rewrite shell history.",
			},
		},
		helpBlock{
			Title: "Flags",
			Lines: []string{
				"  --shell <shell>   scan only one supported shell source (bash or zsh)",
				"  --config <path>   load a specific histkit config file before scanning",
			},
		},
	)
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

const scanWriteBatchSize = 1000

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
		parser, err := history.StreamParserForShell(source.Shell)
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		file, err := os.Open(source.Path)
		if err != nil {
			return fmt.Errorf("scan: open history source %q: %w", source.Path, err)
		}

		batch := make([]history.HistoryEntry, 0, scanWriteBatchSize)
		flushBatch := func() error {
			if len(batch) == 0 {
				return nil
			}

			result, err := index.WriteHistoryEntries(db, batch)
			if err != nil {
				return err
			}

			summary.Attempted += result.Attempted
			summary.Inserted += result.Inserted
			summary.Skipped += result.Skipped
			batch = batch[:0]
			return nil
		}

		parseErr := parser(
			source.Path,
			file,
			func(entry history.HistoryEntry) error {
				batch = append(batch, entry)
				if len(batch) < scanWriteBatchSize {
					return nil
				}
				return flushBatch()
			},
			func(warning history.ParseWarning) error {
				summary.Warnings++
				return nil
			},
		)
		closeErr := file.Close()
		if parseErr != nil {
			return fmt.Errorf("scan: %w", parseErr)
		}
		if closeErr != nil {
			return fmt.Errorf("scan: close history source %q: %w", source.Path, closeErr)
		}
		if err := flushBatch(); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
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
