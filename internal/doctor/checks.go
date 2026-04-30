package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"histkit/internal/config"
	"histkit/internal/history"
)

const (
	StatusOK   = "ok"
	StatusWarn = "warn"
	StatusFail = "fail"
	writeOK    = 2
)

type CheckResult struct {
	Name   string
	Status string
	Detail string
}

type Report struct {
	Checks []CheckResult
}

func (r Report) OverallStatus() string {
	overall := StatusOK
	for _, check := range r.Checks {
		switch check.Status {
		case StatusFail:
			return StatusFail
		case StatusWarn:
			overall = StatusWarn
		}
	}

	return overall
}

func Run(home, configPath string) (Report, error) {
	if home == "" {
		return Report{}, fmt.Errorf("doctor run: home directory is required")
	}

	paths, err := config.DefaultPaths(home)
	if err != nil {
		return Report{}, fmt.Errorf("doctor run: %w", err)
	}

	report := Report{}
	report.Checks = append(report.Checks, checkConfig(home, paths.ConfigFile, configPath))
	report.Checks = append(report.Checks, checkStateDir(paths.StateDir))
	report.Checks = append(report.Checks, checkHistorySources(home))
	report.Checks = append(report.Checks, checkHistoryDB(paths.HistoryDB))
	report.Checks = append(report.Checks, checkFZF())

	return report, nil
}

func checkConfig(home, defaultPath, requestedPath string) CheckResult {
	path := defaultPath
	label := "default config"
	if requestedPath != "" {
		expanded, err := config.ExpandUserPath(requestedPath, home)
		if err != nil {
			return CheckResult{Name: "config", Status: StatusFail, Detail: err.Error()}
		}
		path = expanded
		label = "requested config"
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if requestedPath != "" {
				return CheckResult{Name: "config", Status: StatusFail, Detail: fmt.Sprintf("%s file not found: %s", label, path)}
			}
			return CheckResult{Name: "config", Status: StatusOK, Detail: fmt.Sprintf("default config not present; using built-in defaults (%s)", path)}
		}
		return CheckResult{Name: "config", Status: StatusFail, Detail: fmt.Sprintf("stat config %s: %v", path, err)}
	}
	if info.IsDir() {
		return CheckResult{Name: "config", Status: StatusFail, Detail: fmt.Sprintf("config path is a directory: %s", path)}
	}
	if _, err := config.Load(path); err != nil {
		return CheckResult{Name: "config", Status: StatusFail, Detail: err.Error()}
	}

	return CheckResult{Name: "config", Status: StatusOK, Detail: fmt.Sprintf("%s loaded: %s", label, path)}
}

func checkStateDir(path string) CheckResult {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			parent, accessErr := writableAncestor(path)
			if accessErr != nil {
				return CheckResult{Name: "state_dir", Status: StatusFail, Detail: fmt.Sprintf("state directory missing and no writable ancestor found for: %s", path)}
			}
			return CheckResult{Name: "state_dir", Status: StatusWarn, Detail: fmt.Sprintf("state directory does not exist yet; writable parent detected: %s", parent)}
		}
		return CheckResult{Name: "state_dir", Status: StatusFail, Detail: fmt.Sprintf("stat state directory %s: %v", path, err)}
	}
	if !info.IsDir() {
		return CheckResult{Name: "state_dir", Status: StatusFail, Detail: fmt.Sprintf("state path is not a directory: %s", path)}
	}
	if err := syscall.Access(path, writeOK); err != nil {
		return CheckResult{Name: "state_dir", Status: StatusFail, Detail: fmt.Sprintf("state directory is not writable: %s", path)}
	}

	return CheckResult{Name: "state_dir", Status: StatusOK, Detail: fmt.Sprintf("state directory is writable: %s", path)}
}

func checkHistorySources(home string) CheckResult {
	sources, err := history.CandidateSources(home)
	if err != nil {
		return CheckResult{Name: "history_sources", Status: StatusFail, Detail: err.Error()}
	}

	var found []string
	for _, source := range sources {
		info, err := os.Stat(source.Path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return CheckResult{Name: "history_sources", Status: StatusFail, Detail: fmt.Sprintf("stat history source %s: %v", source.Path, err)}
		}
		if info.IsDir() {
			continue
		}
		file, err := os.Open(source.Path)
		if err != nil {
			return CheckResult{Name: "history_sources", Status: StatusFail, Detail: fmt.Sprintf("open history source %s: %v", source.Path, err)}
		}
		if err := file.Close(); err != nil {
			return CheckResult{Name: "history_sources", Status: StatusFail, Detail: fmt.Sprintf("close history source %s: %v", source.Path, err)}
		}
		found = append(found, fmt.Sprintf("%s (%s)", source.Shell, source.Path))
	}

	if len(found) == 0 {
		return CheckResult{Name: "history_sources", Status: StatusWarn, Detail: "no supported history files detected"}
	}

	return CheckResult{Name: "history_sources", Status: StatusOK, Detail: fmt.Sprintf("readable history sources: %s", join(found))}
}

func checkHistoryDB(path string) CheckResult {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			parent, accessErr := writableAncestor(path)
			if accessErr != nil {
				return CheckResult{Name: "history_db", Status: StatusFail, Detail: fmt.Sprintf("history database missing and no writable ancestor found for: %s", path)}
			}
			return CheckResult{Name: "history_db", Status: StatusWarn, Detail: fmt.Sprintf("history database does not exist yet; writable parent detected: %s", parent)}
		}
		return CheckResult{Name: "history_db", Status: StatusFail, Detail: fmt.Sprintf("stat history database %s: %v", path, err)}
	}
	if info.IsDir() {
		return CheckResult{Name: "history_db", Status: StatusFail, Detail: fmt.Sprintf("history database path is a directory: %s", path)}
	}
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return CheckResult{Name: "history_db", Status: StatusFail, Detail: fmt.Sprintf("open history database %s: %v", path, err)}
	}
	if err := file.Close(); err != nil {
		return CheckResult{Name: "history_db", Status: StatusFail, Detail: fmt.Sprintf("close history database %s: %v", path, err)}
	}

	return CheckResult{Name: "history_db", Status: StatusOK, Detail: fmt.Sprintf("history database is accessible: %s", path)}
}

func writableAncestor(path string) (string, error) {
	current := filepath.Dir(path)
	for {
		info, err := os.Stat(current)
		if err == nil {
			if !info.IsDir() {
				return "", fmt.Errorf("ancestor is not a directory: %s", current)
			}
			if accessErr := syscall.Access(current, writeOK); accessErr != nil {
				return "", accessErr
			}
			return current, nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}

		next := filepath.Dir(current)
		if next == current {
			return "", fmt.Errorf("reached filesystem root without finding writable ancestor")
		}
		current = next
	}
}

func checkFZF() CheckResult {
	path, err := exec.LookPath("fzf")
	if err != nil {
		return CheckResult{Name: "fzf", Status: StatusWarn, Detail: "fzf not found in PATH"}
	}

	return CheckResult{Name: "fzf", Status: StatusOK, Detail: fmt.Sprintf("fzf available at %s", path)}
}

func join(values []string) string {
	switch len(values) {
	case 0:
		return ""
	case 1:
		return values[0]
	default:
		result := values[0]
		for _, value := range values[1:] {
			result += ", " + value
		}
		return result
	}
}
