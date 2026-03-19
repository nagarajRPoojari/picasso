#!/usr/bin/env bash
set -e

bazel run //cli:picasso -- clean "$(pwd)/benchmark/picasso"
bazel run //cli:picasso -- build "$(pwd)/benchmark/picasso"

go build -o "$(pwd)/benchmark/go/main" "$(pwd)/benchmark/go/main.go"

javac -d "$(pwd)/benchmark/java" "$(pwd)/benchmark/java/mandelbrot.java"

cd benchmark/rust
cargo build --release
cd ../..

gcc -O3 -march=native "$(pwd)/benchmark/c/main.c" -o "$(pwd)/benchmark/c/main_bin" -lm

g++ -O3 -std=c++11 -march=native "$(pwd)/benchmark/cpp/main.cpp" -o "$(pwd)/benchmark/cpp/main_bin" -lpthread

hyperfine \
  --warmup 5 \
  --runs 50 \
  -n "Picasso 1.0.2 (no-opt)" './benchmark/picasso/build/a.out 1000' \
  -n "Go 1.24.5 (gc)" './benchmark/go/main 1000' \
  -n "Java 17 (HotSpot JIT)" "java -cp $(pwd)/benchmark/java mandelbrot 1000" \
  -n "Rust 1.94 (Release/LTO)" "./benchmark/rust/target/release/mandelbrot_rust 1000" \
  -n "C (GCC -O3)" './benchmark/c/main_bin 1000' \
  -n "C++ (G++ -O3)" './benchmark/cpp/main_bin 1000' \
  -n "Python 3.12 (CPython)" 'python3.12 ./benchmark/python/main.py 1000' \
  -n "Node.js 24.1 (V8 JIT)" 'node benchmark/nodejs/main.js 1000' \
  --export-markdown benchmark_report.md \
  --export-csv benchmark_results.csv

python3 scripts/generate_chart.py