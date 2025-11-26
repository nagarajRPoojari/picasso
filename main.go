package main

import (
	"github.com/nagarajRPoojari/x-lang/generator"
)

func main() {
	c := generator.NewGenerator()
	c.Build("./project/package.ini")
	c.Compile()
	c.Dump("bin/output.ll")
}
