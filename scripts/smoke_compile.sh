#!/usr/bin/env sh
set -eu

go run ./cmd/sros compile --help
go run ./cmd/sros compile --input examples/intents/minimal.txt
go run ./cmd/sros compile --input examples/intents/file_patch.txt --format json
go run ./cmd/sros compile --input examples/intents/research_task.txt --emit-srxml
