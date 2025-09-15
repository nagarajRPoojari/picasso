package function

import (
	"github.com/llir/llvm/ir"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	typedef "github.com/nagarajRPoojari/x-lang/compiler/type"
)

type Func func(typeHandler *tf.TypeHandler, module *ir.Module, block *ir.Block, args []typedef.Var) typedef.Var
