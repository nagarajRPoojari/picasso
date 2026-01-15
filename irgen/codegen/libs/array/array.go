package array

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	function "github.com/nagarajRPoojari/picasso/irgen/codegen/libs/func"
	_types "github.com/nagarajRPoojari/picasso/irgen/codegen/libs/type"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	typedef "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

type ArrayHandler struct {
}

func NewArrayHandler() *ArrayHandler {
	return &ArrayHandler{}
}

func (t *ArrayHandler) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs["create"] = t.create
	funcs["len"] = t.len
	funcs["shape"] = t.shape
	return funcs
}

func (t *ArrayHandler) create(_ *ir.Func, th *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	dims := make([]value.Value, 0)
	size := _types.NewTypeHandler().Size(nil, th, module, bh, []tf.Var{args[0]})
	for _, i := range args[1:] {
		toInt := th.ImplicitIntCast(bh, i.Load(bh), types.I64)

		dims = append(dims, toInt)
	}
	return typedef.NewArray(bh, args[0].Type(), size.Load(bh), dims, args[0].NativeTypeString())
}

func (t *ArrayHandler) len(_ *ir.Func, th *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	arr := args[0].(*tf.Array)
	length := arr.Len(bh)
	return length
}

func (t *ArrayHandler) shape(_ *ir.Func, th *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	arr := args[0].(*tf.Array)
	shape := arr.LoadShapeArray(bh)
	return shape
}
