package main

import "github.com/nagarajRPoojari/x-lang/compiler"

func main() {
	c := compiler.NewCompiler()
	c.Compile("./project/package.ini")
	c.Dump("bin/output.ll")
}
