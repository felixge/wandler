#!/usr/bin/env bash
pkgs=$(
  find . \
    -type f \
    \! -path "./src*" \
    -name "*_test.go" \
  | xargs -n1 dirname | uniq \
)
echo "${pkgs}" | xargs go test
