#!/usr/bin/env bash
set -e

# build
echo "building picasso"
picasso clean "$(pwd)/benchmark/picasso"
picasso build "$(pwd)/benchmark/picasso"

echo "building go"
go build  -o "$(pwd)/benchmark/go/main" "$(pwd)/benchmark/go/main.go"

hyperfine \
  --warmup 5 \
  --runs 50 \
  -n "go 1.24.5" './benchmark/go/main' \
  -n "picasso 1.0.2" './benchmark/picasso/build/a.out' \
  -n "python 3.12.0" 'python3.12 ./benchmark/python/main.py' \
  --export-markdown benchmark_report.md \
  --export-json benchmark_report.json \
  --export-csv benchmark_report.csv