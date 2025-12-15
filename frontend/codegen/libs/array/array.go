package array

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	function "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/func"
	_types "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/type"
	tf "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	typedef "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
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

func (t *ArrayHandler) create(th *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	dims := make([]value.Value, 0)
	size := _types.NewTypeHandler().Size(th, module, bh, []tf.Var{args[0]})
	for _, i := range args[1:] {
		toInt := th.ImplicitIntCast(bh, i.Load(bh), types.I64)

		dims = append(dims, toInt)
	}
	return typedef.NewArray(bh, args[0].Type(), size.Load(bh), dims)
}
