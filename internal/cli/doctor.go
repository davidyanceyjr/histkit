package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"histkit/internal/config"
	"histkit/internal/doctor"
)

func runDoctor(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	help := fs.Bool("help", false, "show help")
	fs.BoolVar(help, "h", false, "show help")
	configPath := fs.String("config", "", "path to a histkit config file")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("doctor: %w", err)
	}
	if *help {
		writeDoctorUsage(stdout)
		return nil
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("doctor: unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	return executeDoctor(stdout, doctorOptions{ConfigPath: *configPath})
}

func writeDoctorUsage(w io.Writer) {
	writeHelpBlocks(w,
		helpBlock{
			Title: "Usage",
			Lines: []string{
				"  histkit doctor [--config <path>]",
			},
		},
		helpBlock{
			Lines: []string{
				"Inspect the local histkit environment and report configuration and runtime checks.",
			},
		},
		helpBlock{
			Lines: []string{
				"doctor checks config loading, writable state paths, detected history sources, the history index,",
				"fzf availability, and optional systemd --user automation files.",
			},
		},
		helpBlock{
			Title: "Flags",
			Lines: []string{
				"  --config <path>   load or validate a specific histkit config file",
			},
		},
	)
}

type doctorOptions struct {
	ConfigPath string
}

func executeDoctor(stdout io.Writer, opts doctorOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("doctor: detect home directory: %w", err)
	}

	if opts.ConfigPath != "" {
		expandedConfigPath, err := config.ExpandUserPath(opts.ConfigPath, home)
		if err != nil {
			return fmt.Errorf("doctor: %w", err)
		}
		opts.ConfigPath = expandedConfigPath
	}

	report, err := doctor.Run(home, opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("doctor: %w", err)
	}

	fmt.Fprintf(stdout, "doctor overall status: %s\n", report.OverallStatus())
	for _, check := range report.Checks {
		fmt.Fprintf(stdout, "%s: %s - %s\n", check.Name, strings.ToUpper(check.Status), check.Detail)
	}

	return nil
}
