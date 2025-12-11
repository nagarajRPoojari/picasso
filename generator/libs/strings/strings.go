package strings

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/generator/c"

	function "github.com/nagarajRPoojari/x-lang/generator/libs/func"
	"github.com/nagarajRPoojari/x-lang/generator/libs/libutils"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
)

type StringHandler struct {
}

func NewStringHandler() *StringHandler {
	return &StringHandler{}

}

func (t *StringHandler) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.ALIAS_FORMAT] = t.format
	funcs[c.ALIAS_LEN] = t.len
	funcs[c.ALIAS_COMPARE] = t.compare
	return funcs
}

func (t *StringHandler) format(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	scanFn := c.NewInterface(module).Funcs[c.FUNC_FORMAT]
	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}
func (t *StringHandler) len(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	scanFn := c.NewInterface(module).Funcs[c.FUNC_LEN]
	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}
func (t *StringHandler) compare(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	scanFn := c.NewInterface(module).Funcs[c.FUNC_COMPARE]
	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}
