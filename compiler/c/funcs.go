package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// mod.NewFunc(ALLOC, types.I8Ptr, ir.NewParam("", types.I64)),
// mod.NewFunc(RUNTIME_INIT, types.Void),
// mod.NewFunc(ARRAY_ALLOC, types.NewPointer(ARRAYSTRUCT), ir.NewParam("", types.I64), ir.NewParam("", types.I64)),

func (t *Interface) registerFuncs(mod *ir.Module) {

	t.Funcs[ALLOC] = mod.NewFunc(ALLOC, types.I8Ptr, ir.NewParam("", types.I64))
	t.Funcs[RUNTIME_INIT] = mod.NewFunc(RUNTIME_INIT, types.Void)
	t.Funcs[ARRAY_ALLOC] = mod.NewFunc(ARRAY_ALLOC, types.NewPointer(Array), ir.NewParam("", types.I64), ir.NewParam("", types.I64))

}
