#!/bin/bash

set -e

no_test=0
no_lint=0

while (( $# > 0 )); do
   case "$1" in
   	--help)
			printf "run.sh [OPTION]... [DIR]\n"
			printf "options:\n"
			printf "\t--help			Show help\n"
			printf "\t--no-test		Skip tests\n"
			printf "\t--no-lint		Skip linting\n"
			exit 0
      	;;
      --no-test)
			no_test=1
			shift
      	;;
      --no-lint)
			no_lint=1
			shift
			;;
		*)
			break
	      ;;
   esac
done

package="${1%/}" # Strip trailing slash, if any
package="./${package:+$package/}..."

if [[ -n "$1" && -z "$CI" ]]; then
	# Delete project-specific `.mod` files.
	find "$1" -type f -name "*.mod" -maxdepth 1 -exec rm -f {} +
fi

if (( no_test == 0 )); then
	go test -v "$package"
fi

if (( no_lint == 0 )) ; then
	red=$(tput -Txterm-256color setaf 1)
	default_color=$(tput -Txterm-256color sgr0)
	if [[ -x "$(command -v golangci-lint)" ]]; then
		if [[ -z $CI ]]; then
    		golangci-lint fmt "$package"
    	fi
    	golangci-lint run "$package"
	else
		printf "%bgolangci-lint not found%b\n" "$red" "$default_color"
	fi
fi