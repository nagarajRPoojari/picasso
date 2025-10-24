package main

import (
	"github.com/nagarajRPoojari/x-lang/compiler"
)

func main() {
	c := compiler.NewCompiler()
	c.Build("./project/package.ini")
	c.Compile()
	c.Dump("bin/output.ll")
}
