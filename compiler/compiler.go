package compiler

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/parser"
)

type Compiler struct {
	llvm *LLVM
}

func NewCompiler() *Compiler {
	return &Compiler{
		llvm: NewLLVM(),
	}
}

func (t *Compiler) Compile(path string) {
	m := LoadPackages(path)
	t.llvm.ParseAST(m)
}

func (t *Compiler) Dump(file string) {
	t.llvm.Dump(file)
}

func LoadPackages(packagePath string) map[string]ast.BlockStatement {
	cfg, err := ini.Load(packagePath)
	if err != nil {
		panic("failed to read package.ini")
	}

	m := make(map[string]ast.BlockStatement, 0)

	for _, section := range cfg.Sections()[1:] {
		sourceBytes, err := os.ReadFile(section.Key("path").String())
		if err != nil {
			panic(fmt.Sprintf("unable to read: %v", err))
		}
		source := string(sourceBytes)
		tree := parser.Parse(source)
		m[section.Name()] = tree
	}

	return m
}
