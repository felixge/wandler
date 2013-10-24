#!/usr/bin/env bash
set -eu
go install github.com/felixge/wandler/cmd/wandler-server
go install github.com/felixge/wandler/cmd/wandler-worker
