package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"histkit/internal/config"
	"histkit/internal/index"
	"histkit/internal/picker"
	"histkit/internal/snippets"
)

const defaultPickHistoryLimit = 200

func runPick(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("pick", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
	configPath := fs.String("config", "", "path to a histkit config file")
	debug := fs.Bool("debug", false, "print pick progress diagnostics to stderr")

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

	return executePick(context.Background(), stdout, pickOptions{
		ConfigPath: *configPath,
		DebugLog:   newPickDebugLog(stderr, *debug),
	})
}

func writePickUsage(w io.Writer) {
	writeHelpBlocks(w,
		helpBlock{
			Title: "Usage",
			Lines: []string{
				"  histkit pick [--config <path>]",
			},
		},
		helpBlock{
			Lines: []string{
				"Open an fzf picker over indexed history and snippets, then print the selected command to stdout.",
			},
		},
		helpBlock{
			Lines: []string{
				"pick reads from the local history index and the snippet store.",
				"It does not write shell history or expand snippet placeholders on its own.",
			},
		},
		helpBlock{
			Title: "Flags",
			Lines: []string{
				"  --config <path>   load a specific histkit config file before opening the picker",
				"  --debug           print pick progress diagnostics to stderr",
			},
		},
	)
}

type pickOptions struct {
	ConfigPath string
	DebugLog   func(string, ...any)
}

func newPickDebugLog(w io.Writer, enabled bool) func(string, ...any) {
	if !enabled || w == nil {
		return nil
	}

	return func(format string, args ...any) {
		fmt.Fprintf(w, "pick debug: "+format+"\n", args...)
	}
}

func (opts pickOptions) logf(format string, args ...any) {
	if opts.DebugLog != nil {
		opts.DebugLog(format, args...)
	}
}

func (opts pickOptions) timedStep(name string) func() {
	started := time.Now()
	opts.logf("%s start", name)
	return func() {
		opts.logf("%s done in %s", name, time.Since(started).Round(time.Millisecond))
	}
}

func executePick(ctx context.Context, stdout io.Writer, opts pickOptions) error {
	done := opts.timedStep("detect home directory")
	home, err := os.UserHomeDir()
	done()
	if err != nil {
		return fmt.Errorf("pick: detect home directory: %w", err)
	}

	cfg := config.Default()
	if opts.ConfigPath != "" {
		done = opts.timedStep("load config")
		expandedConfigPath, err := config.ExpandUserPath(opts.ConfigPath, home)
		if err != nil {
			return fmt.Errorf("pick: %w", err)
		}
		cfg, err = config.Load(expandedConfigPath)
		done()
		if err != nil {
			return fmt.Errorf("pick: %w", err)
		}
	}

	done = opts.timedStep("resolve default paths")
	paths, err := config.DefaultPaths(home)
	done()
	if err != nil {
		return fmt.Errorf("pick: %w", err)
	}

	snippetPath := paths.SnippetsFile
	if cfg.Snippets.UserFile != "" {
		done = opts.timedStep("resolve snippet path")
		snippetPath, err = config.ExpandUserPath(cfg.Snippets.UserFile, home)
		done()
		if err != nil {
			return fmt.Errorf("pick: %w", err)
		}
	}

	done = opts.timedStep("open history database")
	db, err := index.Open(paths.HistoryDB)
	done()
	if err != nil {
		return fmt.Errorf("pick: %w", err)
	}
	defer db.Close()

	done = opts.timedStep("initialize history schema")
	if err := index.InitSchema(db); err != nil {
		done()
		return fmt.Errorf("pick: %w", err)
	}
	done()

	store := snippets.Store{Path: snippetPath}
	done = opts.timedStep("load picker candidates")
	candidates, err := picker.LoadCandidates(db, store, cfg.Snippets.Enabled, cfg.Snippets.Builtin, defaultPickHistoryLimit)
	done()
	if err != nil {
		return fmt.Errorf("pick: %w", err)
	}
	opts.logf("candidate summary: history_limit=%d total=%d snippets_enabled=%t builtin_snippets=%t snippet_path=%s history_db=%s",
		defaultPickHistoryLimit,
		len(candidates),
		cfg.Snippets.Enabled,
		cfg.Snippets.Builtin,
		snippetPath,
		paths.HistoryDB,
	)

	done = opts.timedStep("launch fzf")
	selected, ok, err := picker.Select(ctx, candidates)
	done()
	if err != nil {
		return fmt.Errorf("pick: %w", err)
	}
	if !ok {
		opts.logf("no selection returned from fzf")
		return nil
	}

	opts.logf("selected %s candidate", selected.Label)
	_, err = fmt.Fprintln(stdout, selected.Command)
	return err
}
