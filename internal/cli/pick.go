package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"histkit/internal/config"
	"histkit/internal/index"
	"histkit/internal/picker"
	"histkit/internal/snippets"
)

const defaultPickHistoryLimit = 200

func runPick(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("pick", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
	configPath := fs.String("config", "", "path to a histkit config file")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("pick: %w", err)
	}
	if *help {
		writePickUsage(stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("pick: unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	return executePick(context.Background(), stdout, pickOptions{ConfigPath: *configPath})
}

func writePickUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit pick [--config <path>]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Open an fzf picker over indexed history and snippets, then print the selected command to stdout.")
}

type pickOptions struct {
	ConfigPath string
}

func executePick(ctx context.Context, stdout io.Writer, opts pickOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("pick: detect home directory: %w", err)
	}

	cfg := config.Default()
	if opts.ConfigPath != "" {
		expandedConfigPath, err := config.ExpandUserPath(opts.ConfigPath, home)
		if err != nil {
			return fmt.Errorf("pick: %w", err)
		}
		cfg, err = config.Load(expandedConfigPath)
		if err != nil {
			return fmt.Errorf("pick: %w", err)
		}
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		return fmt.Errorf("pick: %w", err)
	}

	snippetPath := paths.SnippetsFile
	if cfg.Snippets.UserFile != "" {
		snippetPath, err = config.ExpandUserPath(cfg.Snippets.UserFile, home)
		if err != nil {
			return fmt.Errorf("pick: %w", err)
		}
	}

	db, err := index.Open(paths.HistoryDB)
	if err != nil {
		return fmt.Errorf("pick: %w", err)
	}
	defer db.Close()

	if err := index.InitSchema(db); err != nil {
		return fmt.Errorf("pick: %w", err)
	}

	store := snippets.Store{Path: snippetPath}
	candidates, err := picker.LoadCandidates(db, store, cfg.Snippets.Enabled, cfg.Snippets.Builtin, defaultPickHistoryLimit)
	if err != nil {
		return fmt.Errorf("pick: %w", err)
	}

	selected, ok, err := picker.Select(ctx, candidates)
	if err != nil {
		return fmt.Errorf("pick: %w", err)
	}
	if !ok {
		return nil
	}

	_, err = fmt.Fprintln(stdout, selected.Command)
	return err
}
