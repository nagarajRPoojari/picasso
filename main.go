package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nagarajRPoojari/x-lang/compiler"
	"github.com/nagarajRPoojari/x-lang/parser"
	"github.com/sanity-io/litter"
)

func main() {
	sourceBytes, _ := os.ReadFile("test.rs")
	source := string(sourceBytes)
	start := time.Now()
	ast := parser.Parse(source)

	litter.Dump(ast)

	duration := time.Since(start)

	fmt.Printf("Duration: %v\n", duration)

	c := compiler.NewCompiler()
	c.Compile(ast)
	c.Dump("bin/output.ll")
}
