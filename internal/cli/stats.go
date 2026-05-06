package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"histkit/internal/config"
	"histkit/internal/index"
)

func runStats(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("stats", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
	configPath := fs.String("config", "", "path to a histkit config file")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("stats: %w", err)
	}
	if *help {
		writeStatsUsage(stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("stats: unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	return executeStats(stdout, statsOptions{ConfigPath: *configPath})
}

func writeStatsUsage(w io.Writer) {
	writeHelpBlocks(w,
		helpBlock{
			Title: "Usage",
			Lines: []string{
				"  histkit stats [--config <path>]",
			},
		},
		helpBlock{
			Lines: []string{
				"Read the local SQLite history index and print indexed entry counts by shell and source.",
			},
		},
		helpBlock{
			Lines: []string{
				"stats reports what has already been indexed.",
				"It does not rescan history files or modify them.",
			},
		},
		helpBlock{
			Title: "Flags",
			Lines: []string{
				"  --config <path>   load a specific histkit config file before reading stats",
			},
		},
	)
}

type statsOptions struct {
	ConfigPath string
}

func executeStats(stdout io.Writer, opts statsOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("stats: detect home directory: %w", err)
	}

	if opts.ConfigPath != "" {
		expandedConfigPath, err := config.ExpandUserPath(opts.ConfigPath, home)
		if err != nil {
			return fmt.Errorf("stats: %w", err)
		}
		if _, err := config.Load(expandedConfigPath); err != nil {
			return fmt.Errorf("stats: %w", err)
		}
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		return fmt.Errorf("stats: %w", err)
	}

	db, err := index.Open(paths.HistoryDB)
	if err != nil {
		return fmt.Errorf("stats: %w", err)
	}
	defer db.Close()

	if err := index.InitSchema(db); err != nil {
		return fmt.Errorf("stats: %w", err)
	}

	stats, err := index.QueryHistoryStats(db)
	if err != nil {
		return fmt.Errorf("stats: %w", err)
	}

	fmt.Fprintf(stdout, "Indexed history entries: %d\n", stats.TotalEntries)
	fmt.Fprintln(stdout, "Counts by shell:")
	writeGroupCounts(stdout, stats.ByShell)
	fmt.Fprintln(stdout, "Counts by source:")
	writeGroupCounts(stdout, stats.BySource)

	return nil
}

func writeGroupCounts(w io.Writer, counts []index.GroupCount) {
	if len(counts) == 0 {
		fmt.Fprintln(w, "  (none)")
		return
	}

	for _, item := range counts {
		fmt.Fprintf(w, "  %s: %d\n", item.Name, item.Count)
	}
}
