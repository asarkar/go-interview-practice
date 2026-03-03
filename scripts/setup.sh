#!/usr/bin/env bash
set -eo pipefail
go mod tidy -v
go mod download
