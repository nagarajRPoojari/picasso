package compiler

import (
	"github.com/nagarajRPoojari/x-lang/ast"
)

type Compiler struct {
	llvm *LLVM
}

func NewCompiler() *Compiler {
	return &Compiler{
		llvm: NewLLVM(),
	}
}

func (t *Compiler) Compile(tree ast.BlockStatement) {
	t.llvm.ParseAST(&tree)
	t.llvm.Dump()
}
