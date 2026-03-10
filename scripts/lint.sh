#!/usr/bin/env bash
set -eo pipefail

dir="${1%/}"
package="./${dir:+$dir/}..."

go mod tidy

if [ -z "${CI:-}" ]; then
  go tool golangci-lint fmt "$package"
fi
go tool golangci-lint run "$package"
