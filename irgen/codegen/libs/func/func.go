package function

import (
	"github.com/llir/llvm/ir"
	typedef "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

type Func func(f *ir.Func, typeHandler *typedef.TypeHandler, module *ir.Module, block *bc.BlockHolder, args []typedef.Var) typedef.Var
