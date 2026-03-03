#!/usr/bin/env bash
set -eo pipefail

dir="${1%/}"
package="./${dir:+$dir/}..."

if [ -z "${CI:-}" ]; then
    golangci-lint fmt "$package"
fi
golangci-lint run "$package"
