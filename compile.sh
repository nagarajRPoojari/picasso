#!/bin/bash
set -e

go run main.go

llvm-as bin/output.ll -o bin/output.bc

llc -filetype=obj bin/output.bc -o bin/output.o

brew_prefix=$(brew --prefix bdw-gc)
clang -c c/runtime.c -I"$brew_prefix/include" -o bin/runtime.o

clang bin/output.o bin/runtime.o -L"$brew_prefix/lib" -lgc -o bin/output

start=$(gdate +%s%N)
./bin/output
end=$(gdate +%s%N)
echo "Time taken: $((end - start)) ns"
