#!/usr/bin/env bash
set -eu
./build.bash
pkgs=$(
  find . \
    -type f \
    \! -path "./src*" \
    -name "*_test.go" \
  | xargs -n1 dirname | uniq \
)
echo "${pkgs}" | xargs -n1 go test -v
