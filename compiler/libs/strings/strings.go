package strings

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/compiler/c"

	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	"github.com/nagarajRPoojari/x-lang/compiler/libs/libutils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

type StringHandler struct {
}

func NewStringHandler() *StringHandler {
	return &StringHandler{}

}

func (t *StringHandler) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.FORMAT] = t.format
	funcs[c.LEN] = t.len
	funcs[c.COMPARE] = t.compare
	return funcs
}

func (t *StringHandler) format(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	scanFn := c.NewInterface(module).Funcs[c.FORMAT]
	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}
func (t *StringHandler) len(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	scanFn := c.NewInterface(module).Funcs[c.LEN]
	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}
func (t *StringHandler) compare(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	scanFn := c.NewInterface(module).Funcs[c.COMPARE]
	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}
