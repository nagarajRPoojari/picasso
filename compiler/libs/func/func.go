package function

import (
	"github.com/llir/llvm/ir"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	typedef "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

type Func func(typeHandler *tf.TypeHandler, module *ir.Module, block *bc.BlockHolder, args []typedef.Var) typedef.Var
