package strings

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/c"
	function "github.com/nagarajRPoojari/picasso/irgen/codegen/libs/func"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	typedef "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

type StringsHandler struct {
}

func NewStringsHandler() *StringsHandler {
	return &StringsHandler{}
}

func (t *StringsHandler) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs["substring"] = t.substring
	return funcs
}

func (t *StringsHandler) substring(_ *ir.Func, th *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_STRING_SUBSTRING]
	s := bh.N.NewCall(fn, args[0].Load(bh), args[1].Load(bh), args[2].Load(bh))
	return &tf.String{
		NativeType: types.NewPointer(tf.STRINGSTRUCT),
		Value:      s,
	}
}
