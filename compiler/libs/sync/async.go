package sync

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/compiler/c"

	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	"github.com/nagarajRPoojari/x-lang/compiler/libs/libutils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

type Sync struct {
}

func NewIO() *Sync {
	return &Sync{}

}

func (t *Sync) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.PRINTF] = t.printf
	return funcs
}

func (t *Sync) printf(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	printfFn := c.NewInterface(module).Funcs[c.PRINTF]
	return libutils.CallCFunc(typeHandler, printfFn, bh, args)
}
