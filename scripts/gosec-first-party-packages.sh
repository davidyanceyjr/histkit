#!/usr/bin/env bash

set -euo pipefail

module_path="$(go list -m)"
repo_root="$(pwd -P)"

go list -f '{{with .Module}}{{.Path}} {{end}}{{.Dir}}' ./... | awk -v module_path="$module_path" -v repo_root="$repo_root/" '
$1 == module_path && index($2, repo_root) == 1 {
  print "./" substr($2, length(repo_root) + 1)
}
' | LC_ALL=C sort -u
