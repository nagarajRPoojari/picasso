package array

import (
	"github.com/llir/llvm/ir"
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

func (t *ArrayHandler) create(th *tf.TypeHandler, module *ir.Module, block *ir.Block, args []typedef.Var) (typedef.Var, *ir.Block) {
	dims := make([]value.Value, 0)
	tp := args[0].NativeTypeString()
	size, safe := _types.NewTypeHandler().Size(th, module, block, []tf.Var{args[0]})
	block = safe
	for _, i := range args[1:] {
		dims = append(dims, i.Load(block))
	}
	return typedef.NewArray(block, th.GetLLVMType(tf.Type(tp)), size.Load(block), dims), block
}
