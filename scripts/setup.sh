#!/usr/bin/env bash
set -eo pipefail
go mod download
go mod verify
