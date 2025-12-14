package function

import (
	"github.com/llir/llvm/ir"
	typedef "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
)

type Func func(typeHandler *typedef.TypeHandler, module *ir.Module, block *bc.BlockHolder, args []typedef.Var) typedef.Var
