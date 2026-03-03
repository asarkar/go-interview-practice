#!/usr/bin/env bash
set -eo pipefail

dir="${1%/}"
package="./${dir:+$dir/}..."

if [[ -n "$dir" && -z "${CI:-}" ]]; then
    find "$dir" -type f -name "*.mod" -maxdepth 1 -exec rm -f {} +
    red=$(tput -Txterm-256color setaf 1)
    default_color=$(tput -Txterm-256color sgr0)
    if ! find "$dir" -name "*_test.go" -type f | grep -q .; then
        printf "%bProject '%s' has no test files%b\n" "$red" "$dir" "$default_color"
        exit 1
    fi
    for f in hints.md learning.md README.md; do
        if [[ ! -f "$dir/$f" ]]; then
            printf "%bProject '%s' is missing file '%s'%b\n" "$red" "$dir" "$f" "$default_color"
            exit 1
        fi
    done
fi

go test -v "$package"
