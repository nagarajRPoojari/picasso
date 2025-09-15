#!/bin/bash

go run main.go
llvm-as bin/output.ll -o bin/output.bc
llc -filetype=obj bin/output.bc -o bin/output.o
clang bin/output.o -o bin/output
./bin/output


echo $?