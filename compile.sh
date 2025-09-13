#!/bin/bash

go run main.go
llvm-as output.ll -o output.bc
llc -filetype=obj output.bc -o output.o
clang output.o -o output
./output

echo "ran succesfully.."
echo $?