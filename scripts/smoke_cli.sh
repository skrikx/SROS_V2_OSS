#!/usr/bin/env sh
set -eu

go run ./cmd/sros --help
go run ./cmd/sros config --help
go run ./cmd/sros run --help
go run ./cmd/sros tools --help
go run ./cmd/sros doctor
go run ./cmd/sros status
go run ./cmd/sros inspect
