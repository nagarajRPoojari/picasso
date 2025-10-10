package array

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	_types "github.com/nagarajRPoojari/x-lang/compiler/libs/type"
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

func (t *ArrayHandler) create(th *tf.TypeHandler, module *ir.Module, bh tf.BlockHolder, args []typedef.Var) (typedef.Var, tf.BlockHolder) {
	dims := make([]value.Value, 0)
	size, safe := _types.NewTypeHandler().Size(th, module, bh, []tf.Var{args[0]})
	bh = safe
	for _, i := range args[1:] {
		toInt, safeN := th.ImplicitIntCast(bh.N, i.Load(bh.N), types.I64)
		bh.N = safeN

		dims = append(dims, toInt)
	}
	return typedef.NewArray(bh, args[0].Type(), size.Load(bh.N), dims), bh
}
