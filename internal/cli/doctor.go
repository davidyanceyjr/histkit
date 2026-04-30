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
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  histkit doctor [--config <path>]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Inspect the local histkit environment and report basic configuration and runtime checks.")
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
