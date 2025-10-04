package array

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	typedef "github.com/nagarajRPoojari/x-lang/compiler/type"
)

type ArrayHandler struct {
}

func NewArrayHandler() *ArrayHandler {
	return &ArrayHandler{}
}

func (t *ArrayHandler) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs["create"] = t.create
	return funcs
}

func (t *ArrayHandler) create(ArrayHandler *tf.TypeHandler, module *ir.Module, block *ir.Block, args []typedef.Var) (typedef.Var, *ir.Block) {
	dims := make([]value.Value, 0)
	for _, i := range args {
		dims = append(dims, i.Load(block))
	}
	return typedef.NewArray(block, types.I64, dims), block
}
