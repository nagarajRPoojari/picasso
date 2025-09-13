package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nagarajRPoojari/x-lang/compiler"
	"github.com/nagarajRPoojari/x-lang/parser"
	"github.com/sanity-io/litter"
)

// "github.com/nagarajRPoojari/x-lang/compiler"

func main() {
	sourceBytes, _ := os.ReadFile("test.rs")
	source := string(sourceBytes)
	start := time.Now()
	ast := parser.Parse(source)
	duration := time.Since(start)

	litter.Dump(ast)
	fmt.Printf("Duration: %v\n", duration)

	compiler.NewCompiler().Compile(ast)
}

// func main() {
// 	fmt.Println(base64.StdEncoding.EncodeToString([]byte("abcd")))
// }
