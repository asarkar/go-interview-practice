#!/usr/bin/env bash
set -eo pipefail

dir="${1%/}"
package="./${dir:+$dir/}..."

if [[ -n "$dir" && -z "${CI:-}" ]]; then
    find "$dir" -type f -name "*.mod" -maxdepth 1 -exec rm -f {} +
    num_files=$(ls -1F "$dir" | wc -l)
    min_files=5
    if (( num_files < min_files )); then
        red=$(tput -Txterm-256color setaf 1)
        default_color=$(tput -Txterm-256color sgr0)
        printf "%bProject '%s' has fewer than %d files%b\n" "$red" "$dir" $min_files "$default_color"
        exit 1
    fi
fi

go test -v "$package"
