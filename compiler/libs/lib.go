package libs

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler"
	typedef "github.com/nagarajRPoojari/x-lang/compiler/type"
)

type Func func(llvm *compiler.LLVM, block *ir.Block, args []*ast.Expression) typedef.Var

type Module struct {
	Name  string
	Funcs map[string]Func
}
